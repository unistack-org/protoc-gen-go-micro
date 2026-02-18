package main

//go:generate go build -o $GOPATH/bin/protoc-gen-go-micro .

// stream_test: micro  + grpc
//go:generate sh -xc "protoc -I./example -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components=micro,standalone=false,paths=source_relative:./example/out/go ./example/stream_test.proto"
//go:generate sh -xc "protoc -I./example -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components=grpc,standalone=true,paths=source_relative:./example/out/go/grpc ./example/stream_test.proto"

// align_test: fieldalignment
//go:generate sh -xc "protoc -I./example -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./example/out/align ./example/align_test.proto"
//go:generate sh -xc "protoc -I./example -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components=micro,standalone=false,fieldaligment=true,tag_path=./example/out/align,paths=source_relative:./example/out/align ./example/align_test.proto"
