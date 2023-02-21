package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

var gorillaPackageFiles map[protogen.GoPackageName]struct{}

func (g *Generator) gorillaGenerate(component string, plugin *protogen.Plugin) error {
	gorillaPackageFiles = make(map[protogen.GoPackageName]struct{})
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if len(file.Services) == 0 {
			continue
		}
		if _, ok := gorillaPackageFiles[file.GoPackageName]; ok {
			continue
		}

		gorillaPackageFiles[file.GoPackageName] = struct{}{}
		gname := "micro" + "_" + component + ".pb.go"

		path := file.GoImportPath
		if g.standalone {
			path = "."
		}
		gfile := plugin.NewGeneratedFile(gname, path)

		gfile.P("// Code generated by protoc-gen-go-micro. DO NOT EDIT.")
		gfile.P("// protoc-gen-go-micro version: " + versionComment)
		gfile.P()
		gfile.P("package ", file.GoPackageName)
		gfile.P()

		gfile.Import(fmtPackage)
		gfile.Import(httpPackage)
		gfile.Import(reflectPackage)
		gfile.Import(stringsPackage)
		gfile.Import(gorillaMuxPackage)

		gfile.P("func RegisterHandlers(r *", gorillaMuxPackage.Ident("Router"), ", h interface{}, eps []", microServerHttpPackage.Ident("EndpointMetadata"), ") error {")
		gfile.P("v := ", reflectPackage.Ident("ValueOf"), "(h)")
		gfile.P("if v.NumMethod() < 1 {")
		gfile.P(`return `, fmtPackage.Ident("Errorf"), `("handler has no methods: %T", h)`)
		gfile.P("}")
		gfile.P("for _, ep := range eps {")
		gfile.P(`idx := `, stringsPackage.Ident("Index"), `(ep.Name, ".")`)
		gfile.P(`if idx < 1 || len(ep.Name) <= idx {`)
		gfile.P(`return `, fmtPackage.Ident("Errorf"), `("invalid endpoint name: %s", ep.Name)`)
		gfile.P("}")
		gfile.P(`name := ep.Name[idx+1:]`)
		gfile.P("m := v.MethodByName(name)")
		gfile.P("if !m.IsValid() || m.IsZero() {")
		gfile.P(`return `, fmtPackage.Ident("Errorf"), `("invalid handler, method %s not found", name)`)
		gfile.P("}")
		gfile.P("rh, ok := m.Interface().(func(", httpPackage.Ident("ResponseWriter"), ", *", httpPackage.Ident("Request"), "))")
		gfile.P("if !ok {")
		gfile.P(`return `, fmtPackage.Ident("Errorf"), `("invalid handler: %#+v", m.Interface())`)
		gfile.P("}")
		gfile.P(`r.HandleFunc(ep.Path, rh).Methods(ep.Method).Name(ep.Name)`)
		gfile.P("}")
		gfile.P("return nil")
		gfile.P("}")
	}

	return nil
}
