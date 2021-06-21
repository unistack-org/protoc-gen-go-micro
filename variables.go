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
	microApiPackage        = protogen.GoImportPath("github.com/unistack-org/micro/v3/api")
	microClientPackage     = protogen.GoImportPath("github.com/unistack-org/micro/v3/client")
	microServerPackage     = protogen.GoImportPath("github.com/unistack-org/micro/v3/server")
	microClientHttpPackage = protogen.GoImportPath("github.com/unistack-org/micro-client-http/v3")
	timePackage            = protogen.GoImportPath("time")
	deprecationComment     = "// Deprecated: Do not use."
	versionComment         = "v3.4.0"
)
