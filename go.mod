module github.com/unistack-org/protoc-gen-go-micro/v3

go 1.16

require (
	github.com/fatih/structtag v1.2.0
	github.com/unistack-org/micro-codec-jsonpb/v3 v3.7.5
	go.unistack.org/micro-proto/v3 v3.1.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.27.1
)

//replace go.unistack.org/micro-proto => ../micro-proto
