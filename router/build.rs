fn main() -> Result<(), Box<dyn std::error::Error>> {
    let proto_root = "./idl";

    tonic_build::configure()
        .build_server(true)
        .compile(&["./idl/security.proto"], &[proto_root])
        .unwrap_or_else(|e| panic!("Failed to compile protos {:?}", e));

    Ok(())
}