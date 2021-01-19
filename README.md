# `protoc-gen-gotemplate`
:open_file_folder: protocol generator + golang text/template (protobuf)

A generic **code**/script/data generator based on [Protobuf](https://developers.google.com/protocol-buffers/).

---

This project is a generator plugin for the Google Protocol Buffers compiler (`protoc`).

The plugin parses **protobuf** files, generates an **ast**, and walks a local **templates directory** to generate files using the [Golang's `text/template` engine](https://golang.org/pkg/text/template/).

## Philosophy

* protobuf-first
* no built-in template, only user defined templates
* kiss, *keep it stupid simple*

## Under the hood

1. the *user* `protobuf` files are parsed by [`protoc`](https://github.com/google/protobuf/releases)
2. the `ast` is generated by [`protoc-gen-go` helpers](https://github.com/golang/protobuf/tree/master/protoc-gen-go)
3. the `ast` is given to [Golang's `text/template` engine](https://golang.org/pkg/text/template/) for each *user* template files
4. the *funcmap* enriching the template engine is based on [Masterminds/sprig](https://github.com/Masterminds/sprig), and contains type-manipulation, iteration and language-specific helpers

## Web editor

![Web editor screenshot](https://github.com/moul/protoc-gen-gotemplate/raw/master/assets/web-editor.jpg)

[Demo server](http://protoc-gen-gotemplate.m.42.am/)

## Usage

`protoc-gen-gotemplate` requires a **template_dir** directory *(by default `./templates`)*.

Every file ending with `.tmpl` will be processed and written to the destination folder, following the file hierarchy of the `template_dir`, and remove the `.tmpl` extension.

---

```console
$> ls -R
input.proto     templates/doc.txt.tmpl      templates/config.json.tmpl
$> protoc --gotemplate_out=. input.proto
$> ls -R
input.proto     templates/doc.txt.tmpl      templates/config.json.tmpl
doc.txt         config.json
```

### Options

You can specify custom options, as follow:

```console
$> protoc --gotemplate_out=debug=true,template_dir=/path/to/template/directory:. input.proto
```

| Option                | Default Value | Accepted Values           | Description
|-----------------------|---------------|---------------------------|-----------------------
| `template_repo`       | ``            | url in form schema://domain | path to repo with optional branch or revision after @ sign
| `template_dir`        | `./template`  | absolute or relative path | path to look for templates
| `destination_dir`     | `.`           | absolute or relative path | base path to write output
| `single-package-mode` | *false*       | `true` or `false`         | if *true*, `protoc` won't accept multiple packages to be compiled at once (*!= from `all`*), but will support `Message` lookup across the imported protobuf dependencies
| `debug`               | *false*       | `true` or `false`         | if *true*, `protoc` will generate a more verbose output
| `all`                 | *false*       | `true` or `false`         | if *true*, protobuf files without `Service` will also be parsed
| `components`          | `micro`       | `micro, grpc, http, chi, gorilla | some values cant coexists like gorilla/chi or grpc/http
##### Hints

Shipping the templates with your project is very smart and useful when contributing on git-based projects.

Another workflow consists in having a dedicated repository for generic templates which is then versioned and vendored with multiple projects (npm package, golang vendor package, ...)

See [examples](./examples).

## Funcmap

This project uses [Masterminds/sprig](https://github.com/Masterminds/sprig) library and additional functions to extend the builtin [text/template](https://golang.org/pkg/text/template) helpers.

Non-exhaustive list of new helpers:

* **all the functions from [sprig](https://github.com/Masterminds/sprig)**
* `add`
* `boolFieldExtension`
* `camelCase`
* `contains`
* `divide`
* `fieldMapKeyType`
* `fieldMapValueType`
* `first`
* `getEnumValue`
* `getMessageType`
* `getProtoFile`
* `goNormalize`
* `goTypeWithPackage`
* `goType`
* `goZeroValue`
* `haskellType`
* `httpBody`
* `httpPath`
* `httpPathsAdditionalBindings`
* `httpVerb`
* `index`
* `int64FieldExtension`
* `isFieldMap`
* `isFieldMessageTimeStamp`
* `isFieldMessage`
* `isFieldRepeated`
* `jsSuffixReserved`
* `jsType`
* `json`
* `kebabCase`
* `last`
* `leadingComment`
* `leadingDetachedComments`
* `lowerCamelCase`
* `lowerFirst`
* `lowerGoNormalize`
* `multiply`
* `namespacedFlowType`
* `prettyjson`
* `replaceDict`
* `shortType`
* `snakeCase`
* `splitArray`
* `stringFieldExtension`
* `stringMethodOptionsExtension`
* `string`
* `subtract`
* `trailingComment`
* `trimstr`
* `upperFirst`
* `urlHasVarsFromMessage`

See the project helpers for the complete list.

## Install

* Install the **Go** compiler and tools from https://golang.org/doc/install
* Install **protobuf**: `go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`
* Install **protoc-gen-gotemplate**: `go get -u moul.io/protoc-gen-gotemplate`

## Docker

* automated docker hub build: [https://hub.docker.com/r/moul/protoc-gen-gotemplate/](https://hub.docker.com/r/moul/protoc-gen-gotemplate/)
* Based on [http://github.com/znly/protoc](http://github.com/znly/protoc)

Usage:

```console
$> docker run --rm -v "$(pwd):$(pwd)" -w "$(pwd)" moul/protoc-gen-gotemplate -I. --gotemplate_out=./output/ ./*.proto
```

## Projects using `protoc-gen-gotemplate`

* [kafka-gateway](https://github.com/moul/kafka-gateway/): Kafka gateway/proxy (gRPC + http) using Go-Kit
* [translator](https://github.com/moul/translator): Translator Micro-service using Gettext and Go-Kit
* [acl](https://github.com/moul/acl): ACL micro-service (gRPC/protobuf + http/json)

## See also

* [pbhbs](https://github.com/gponsinet/pbhbs): protobuf gen based on handlebarjs template

## License

MIT
