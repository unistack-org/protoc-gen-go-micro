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
	microMetadataPackage   = protogen.GoImportPath("go.unistack.org/micro/v3/metadata")
	microClientPackage     = protogen.GoImportPath("go.unistack.org/micro/v3/client")
	microServerPackage     = protogen.GoImportPath("go.unistack.org/micro/v3/server")
	microClientHttpPackage = protogen.GoImportPath("go.unistack.org/micro-client-http/v3")
	microServerHttpPackage = protogen.GoImportPath("go.unistack.org/micro-server-http/v3")
	microCodecPackage      = protogen.GoImportPath("go.unistack.org/micro/v3/codec")
	microErrorsPackage     = protogen.GoImportPath("go.unistack.org/micro/v3/errors")
	grpcPackage            = protogen.GoImportPath("google.golang.org/grpc")
	timePackage            = protogen.GoImportPath("time")
	deprecationComment     = "// Deprecated: Do not use."
	versionComment         = "v3.10.3"
)
