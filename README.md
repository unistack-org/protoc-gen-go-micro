# `protoc-gen-micro`
protobuf plugin to generate helper code for micro framework

A generic **code**/script/data generator based on [Protobuf](https://developers.google.com/protocol-buffers/).

---

This project is a generator plugin for the Google Protocol Buffers compiler (`protoc`).

## Usage

```console
$> protoc --micro_out=debug=true,components="micro|http":. input.proto
```

| Option                | Default Value | Accepted Values           | Description
|-----------------------|---------------|---------------------------|-----------------------
| `debug`               | *false*       | `true` or `false`         | if *true*, `protoc` will generate a more verbose output
| `components`          | `micro`       | `micro rpc http chi gorilla` | some values cant coexists like gorilla/chi or rpc/http, values must be concatinated with pipe symbol

## Install

* Install the **go** compiler and tools from https://golang.org/doc/install
* Install **protoc-gen-go**: `go get google.golang.org/protobuf/cmd/protoc-gen-go`
* Install **protoc-gen-micro**: `go get github.com/unistack-org/protoc-gen-micro`
