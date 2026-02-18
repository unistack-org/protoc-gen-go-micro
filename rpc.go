package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func (g *Generator) rpcGenerate(component string, plugin *protogen.Plugin, genClient bool, genServer bool) error {
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if len(file.Services) == 0 {
			continue
		}

		gFile := g.newGeneratedFile(plugin, file, component, genClient, genServer)

		for _, service := range file.Services {
			if genClient {
				g.generateServiceClient(gFile, file, service)
				g.generateServiceClientMethods(gFile, file, service, component)
			}
			if genServer {
				g.generateServiceServer(gFile, file, service)
				g.generateServiceServerMethods(gFile, service)
				g.generateServiceRegister(gFile, file, service, component)
			}
			if component == "grpc" && g.reflection {
				g.generateServiceDesc(gFile, file, service)
			}
		}
	}

	return nil
}
