package main

//go:generate sh -xc "protoc -I./example -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v5) --go-micro_out=components=graphqls,graphql_file=./schema.graphql:./example example/example.proto"

//go:generate go build -o $GOPATH/bin/protoc-gen-go-micro .
//go:generate mkdir -p ./example/out/go/grpc ./example/out/align

// stream_test: micro  + grpc
//go:generate sh -xc "protoc -I./example -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v5) --go_out=paths=source_relative:./example/out/go --go-micro_out=components=micro,standalone=false,paths=source_relative:./example/out/go ./example/stream_test.proto"
//go:generate sh -xc "protoc -I./example -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v5) --go-micro_out=components=grpc,standalone=true,paths=source_relative:./example/out/go/grpc ./example/stream_test.proto"
