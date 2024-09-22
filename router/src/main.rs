use anyhow::Result;

mod jwt;

// `cargo run -- -s ../../graphql/supergraph.graphql -c ./router.yaml`
fn main() -> Result<()> {
    apollo_router::main()
}