package main

import "google.golang.org/protobuf/compiler/protogen"

var (
	reflectPackage         = protogen.GoImportPath("reflect")
	stringsPackage         = protogen.GoImportPath("strings")
	fmtPackage             = protogen.GoImportPath("fmt")
	contextPackage         = protogen.GoImportPath("context")
	httpPackage            = protogen.GoImportPath("net/http")
	gorillaMuxPackage      = protogen.GoImportPath("github.com/gorilla/mux")
	chiPackage             = protogen.GoImportPath("github.com/go-chi/chi/v5")
	chiMiddlewarePackage   = protogen.GoImportPath("github.com/go-chi/chi/v5/middleware")
	microMetadataPackage   = protogen.GoImportPath("go.unistack.org/micro/v4/metadata")
	microClientPackage     = protogen.GoImportPath("go.unistack.org/micro/v4/client")
	microServerPackage     = protogen.GoImportPath("go.unistack.org/micro/v4/server")
	microClientHttpPackage = protogen.GoImportPath("go.unistack.org/micro-client-http/v4")
	microServerHttpPackage = protogen.GoImportPath("go.unistack.org/micro-server-http/v4")
	microCodecPackage      = protogen.GoImportPath("go.unistack.org/micro-proto/v4/codec")
	microErrorsPackage     = protogen.GoImportPath("go.unistack.org/micro/v4/errors")
	microOptionsPackage    = protogen.GoImportPath("go.unistack.org/micro/v4/options")
	grpcPackage            = protogen.GoImportPath("google.golang.org/grpc")
	protojsonPackage       = protogen.GoImportPath("google.golang.org/protobuf/encoding/protojson")
	timePackage            = protogen.GoImportPath("time")
	deprecationComment     = "// Deprecated: Do not use."
	versionComment         = "v4.0.2"
)
