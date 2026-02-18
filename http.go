package main

import "google.golang.org/protobuf/compiler/protogen"

func (g *Generator) httpGenerate(component string, plugin *protogen.Plugin, genClient bool, genServer bool) error {
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		if len(file.Services) == 0 {
			continue
		}

		gfile := g.newGeneratedFile(plugin, file, component, genClient, genServer)
		if genClient {
			gfile.Import(microClientHttpPackage)
		}

		for _, service := range file.Services {
			g.generateServiceEndpoints(gfile, service, component)
			if genClient {
				g.generateServiceClient(gfile, file, service)
				g.generateServiceClientMethods(gfile, file, service, component)
			}
			if genServer {
				g.generateServiceServer(gfile, file, service)
				g.generateServiceServerMethods(gfile, service)
				g.generateServiceRegister(gfile, file, service, component)
			}
		}
	}

	return nil
}
