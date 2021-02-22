package main

import (
	"flag"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet

	flagDebug := flags.Bool("debug", false, "")
	flagComponents := flags.String("components", "micro", "")
	flagPaths := flag.String("paths", "", "")
	flagModule := flag.String("module", "", "")

	opts := &protogen.Options{
		ParamFunc: flags.Set,
	}

	g := &Generator{
		debug:      *flagDebug,
		components: strings.Split(*flagComponents, "|"),
		paths:      *flagPaths,
		module:     *flagModule,
	}

	opts.Run(g.Generate)
}

type Generator struct {
	debug      bool
	components []string
	paths      string
	module     string
}

func (g *Generator) Generate(plugin *protogen.Plugin) error {
	var err error

	// Protoc passes a slice of File structs for us to process
	for _, component := range g.components {
		switch component {
		case "micro":
			err = g.microGenerate(component, plugin)
		case "http":
			err = g.httpGenerate(component, plugin)
		case "grpc", "rpc":
			err = g.rpcGenerate(component, plugin)
		case "openapi", "swagger":
			err = g.openapiGenerate(component, plugin)
		default:
			err = fmt.Errorf("unknown component: %s", component)
		}

		if err != nil {
			plugin.Error(err)
			return err
		}

	}

	return nil
}
