# `protoc-gen-go-micro`
protobuf plugin to generate helper code for micro framework

A generic **code**/script/data generator based on [Protobuf](https://developers.google.com/protocol-buffers/).

---

This project is a generator plugin for the Google Protocol Buffers compiler (`protoc`).

## Usage

```console
$> protoc --go_micro_out=debug=true,components="micro|http":. input.proto
```

| Option                | Default Value | Accepted Values           | Description
|-----------------------|---------------|---------------------------|-----------------------
| `tag_path`            | `.`           | `any local path`          | path contains generated protobuf code that needs to be tagged
| `debug`               | *false*       | `true` or `false`         | if *true*, `protoc` will generate a more verbose output
| `components`          | `micro`       | `micro rpc http chi gorilla client server` | some values can't coexists like gorilla/chi or rpc/http, values must be concatinated with pipe symbol

## Install

* Install the **go** compiler and tools from https://golang.org/doc/install
* Install **protoc-gen-go**: `go install google.golang.org/protobuf/cmd/protoc-gen-go`
* Install **protoc-gen-go-micro**: `go install go.unistack.org/protoc-gen-go-micro/v4`
