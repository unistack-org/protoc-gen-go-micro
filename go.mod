module github.com/unistack-org/protoc-gen-micro/v3

go 1.16

require (
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/unistack-org/micro-proto v0.0.2-0.20210227213711-77c7563bd01e
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.25.0
)

//replace github.com/unistack-org/micro-proto => ../micro-proto
