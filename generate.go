package main

//go:generate sh -xc "protoc -I./example -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components=graphqls,graphql_file=./schema.graphql:./example example/example.proto"
