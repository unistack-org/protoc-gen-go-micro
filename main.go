package main

import (
	"flag"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	flagDebug      = flag.Bool("debug", false, "")
	flagStandalone = flag.Bool("standalone", false, "")
	flagComponents = flag.String("components", "micro|rpc", "")
)

func main() {
	opts := &protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}

	g := &Generator{}

	opts.Run(g.Generate)
}

type Generator struct {
	components string
	standalone bool
	debug      bool
}

func (g *Generator) Generate(plugin *protogen.Plugin) error {
	var err error

	g.standalone = *flagStandalone
	g.debug = *flagDebug
	g.components = *flagComponents
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	// Protoc passes a slice of File structs for us to process
	for _, component := range strings.Split(g.components, "|") {
		switch component {
		case "micro":
			err = g.microGenerate(component, plugin)
		case "http":
			err = g.httpGenerate(component, plugin)
		case "grpc", "rpc":
			err = g.rpcGenerate("rpc", plugin)
		case "gorilla":
			err = g.gorillaGenerate(component, plugin)
		case "chi":
			err = g.chiGenerate(component, plugin)
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
