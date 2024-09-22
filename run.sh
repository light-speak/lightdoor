# /bin/zsh

rover supergraph compose --config ./supergraph.yml > supergraph-schema.graphql

cargo run -- -c ./router.yaml -s ./supergraph-schema.graphql --dev