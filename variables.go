package main

import (
	"runtime/debug"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	ioPackage              = protogen.GoImportPath("io")
	graphqlPackage         = protogen.GoImportPath("github.com/99designs/gqlgen/graphql")
	reflectPackage         = protogen.GoImportPath("reflect")
	stringsPackage         = protogen.GoImportPath("strings")
	fmtPackage             = protogen.GoImportPath("fmt")
	contextPackage         = protogen.GoImportPath("context")
	httpPackage            = protogen.GoImportPath("net/http")
	gorillaMuxPackage      = protogen.GoImportPath("github.com/gorilla/mux")
	chiPackage             = protogen.GoImportPath("github.com/go-chi/chi/v5")
	chiMiddlewarePackage   = protogen.GoImportPath("github.com/go-chi/chi/v5/middleware")
	microMetadataPackage   = protogen.GoImportPath("go.unistack.org/micro/v5/metadata")
	microClientPackage     = protogen.GoImportPath("go.unistack.org/micro/v5/client")
	microServerPackage     = protogen.GoImportPath("go.unistack.org/micro/v5/server")
	microClientHttpPackage = protogen.GoImportPath("go.unistack.org/micro-client-http/v5")
	microServerHttpPackage = protogen.GoImportPath("go.unistack.org/micro-server-http/v5")
	microCodecPackage      = protogen.GoImportPath("go.unistack.org/micro-proto/v5/codec")
	microErrorsPackage     = protogen.GoImportPath("go.unistack.org/micro/v4/errors")
	grpcPackage            = protogen.GoImportPath("google.golang.org/grpc")
	protojsonPackage       = protogen.GoImportPath("google.golang.org/protobuf/encoding/protojson")
	timePackage            = protogen.GoImportPath("time")
	deprecationComment     = "// Deprecated: Do not use."
	versionComment         = func() string {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
		return "(devel)"
	}()
)
