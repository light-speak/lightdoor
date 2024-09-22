use std::ops::ControlFlow;
use std::sync::Arc;
use apollo_router::layers::ServiceBuilderExt;
use apollo_router::plugin::{Plugin, PluginInit};
use apollo_router::register_plugin;
use apollo_router::services::supergraph::{BoxService, Request, Response};
use tower::ServiceExt;
use schemars::JsonSchema;
use serde::Deserialize;
use tokio::sync::Mutex;
use tonic::transport::Channel;
use tower::{BoxError, ServiceBuilder};

use proto::security_service_client::SecurityServiceClient;
use proto::TokenRequest;
use crate::jwt::proto::TokenResponse;

pub mod proto {
    tonic::include_proto!("proto"); 
}

#[derive(Debug, Default, Deserialize, JsonSchema)]
struct Conf {
    enable: bool,
}

struct JwtPlugin {
    config: Conf,
    kitex_client: Arc<Mutex<SecurityServiceClient<Channel>>>,
}

#[async_trait::async_trait]
impl Plugin for JwtPlugin {
    type Config = Conf;

    // 插件初始化，创建 gRPC 客户端
    async fn new(init: PluginInit<Self::Config>) -> Result<Self, BoxError> {
        // 创建 gRPC 客户端连接到 Kitex 服务
        // 使用 TCP 连接替代 Unix 套接字连接
        let kitex_client = SecurityServiceClient::connect("http://security-service:4000").await?;

        Ok(JwtPlugin {
            config: init.config,
            kitex_client: Arc::new(Mutex::new(kitex_client)),
        })
    }

    // 定义插件的超图服务逻辑
    fn supergraph_service(
        &self,
        service: BoxService,
    ) -> BoxService {
        ServiceBuilder::new()
            .oneshot_checkpoint_async({
                let kitex_client = Arc::clone(&self.kitex_client);
                move |mut req: Request| {
                    let kitex_client = Arc::clone(&kitex_client);
                    async move {
                        if let Some(auth_header) = req.supergraph_request.headers().get("Authorization") {
                            if let Ok(token) = auth_header.to_str() {
                                if token.starts_with("Bearer ") {
                                    let token = &token[7..];
                                    // println!("token {}", token);
                                    let _guard = req.context.enter_active_request();

                                    let request = tonic::Request::new(TokenRequest {
                                        token: token.to_string(),
                                    });
                                    let user_id_result = {
                                        let mut kitex_client = kitex_client.lock().await;  // 加锁使用 gRPC 客户端
                                        kitex_client.get_security_user_id(request).await
                                    };

                                    drop(_guard);
                                    return match user_id_result {
                                        Ok(res) => {
                                            // req.context.insert("UserId", res.get_ref().user_id)?;
                                            req.supergraph_request.headers_mut().insert("X-User-Id", res.get_ref().user_id.to_string().parse().unwrap());
                                            // println!("{:?}", req.supergraph_request.headers().get("X-User-Id"));
                                            Ok(ControlFlow::Continue(req))
                                        }
                                        Err(err) => {
                                            // println!("Failed to get user ID: {:?}", err);
                                            Ok(ControlFlow::Continue(req))
                                        }
                                    };
                                }
                            }
                        }
                        Ok(ControlFlow::Continue(req))
                    }
                }
            })
            .service(service)
            .boxed()
    }
}


register_plugin!("cos", "jwt_plugin", JwtPlugin);
