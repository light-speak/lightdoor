# 使用 rust 构建镜像作为基础
# renovate-automation: rustc version
FROM rust:1.76 as build

# 设置构建的工作目录
WORKDIR /app

# 更新构建镜像并安装所需包
RUN apt-get update
RUN apt-get -y install \
    npm \
    protobuf-compiler \
    cmake

# 添加 rustfmt
RUN rustup component add rustfmt

# 将路由器源代码复制到构建环境
COPY . .

# 构建并安装自定义二进制文件
RUN cargo build --release

# 为配置和模式创建目录
RUN mkdir -p /dist/config && \
    mv target/release/router /dist && \
    mv router.yaml /dist/router.yaml

FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get -y install \
    ca-certificates 

# 从构建镜像复制所需文件
COPY --from=build --chown=root:root /dist /dist
COPY --from=build --chown=root:root /app/schema /dist/config

WORKDIR /dist

# 确保我们可以运行 router
RUN chmod 755 /dist/router

# ENTRYPOINT ["/bin/sh", "-c", "tail -f /dev/null"]
ENTRYPOINT ["/dist/router", "-c", "/dist/router.yaml", "-s", "/dist/config/supergraph-schema.graphql", "--hr"]

# 使用说明：
# 要通过 docker-compose 运行这个镜像并传入配置文件，可以在 docker-compose.yml 中设置：
# volumes:
#   - ./config:/dist/config
# environment:
#   - APOLLO_ROUTER_CONFIG_PATH=/dist/config/config.yaml
#   - APOLLO_ROUTER_SUPERGRAPH_PATH=/dist/config/supergraph-schema.graphql