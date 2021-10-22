package main

import (
	"flag"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	flagDebug       = flag.Bool("debug", false, "debug output")
	flagStandalone  = flag.Bool("standalone", false, "generate file to standalone dir")
	flagComponents  = flag.String("components", "micro|rpc|http|client|server|openapiv3", "specify components to generate")
	flagTagPath     = flag.String("tag_path", "", "tag rewriting dir")
	flagOpenapiFile = flag.String("openapi_file", "apidocs.swagger.json", "openapi file name")
	flagHelp        = flag.Bool("help", false, "display help message")
)

func main() {
	opts := &protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}

	flag.Parse()

	if *flagHelp {
		flag.PrintDefaults()
		return
	}

	g := &Generator{}

	opts.Run(g.Generate)
}

type Generator struct {
	components  string
	standalone  bool
	debug       bool
	tagPath     string
	openapiFile string
}

func (g *Generator) Generate(plugin *protogen.Plugin) error {
	var err error

	g.standalone = *flagStandalone
	g.debug = *flagDebug
	g.components = *flagComponents
	g.tagPath = *flagTagPath
	g.openapiFile = *flagOpenapiFile
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	var genClient bool
	var genServer bool

	if strings.Contains(g.components, "server") {
		genServer = true
	}
	if strings.Contains(g.components, "client") {
		genClient = true
	}
	if strings.Contains(g.components, "rpc") || strings.Contains(g.components, "http") {
		if !genServer && !genClient {
			genServer = true
			genClient = true
		}
	}

	// Protoc passes a slice of File structs for us to process
	for _, component := range strings.Split(g.components, "|") {
		switch component {
		case "server", "client":
			continue
		case "micro":
			err = g.microGenerate(component, plugin, genClient, genServer)
		case "http":
			err = g.httpGenerate(component, plugin, genClient, genServer)
		case "grpc", "rpc":
			err = g.rpcGenerate("rpc", plugin, genClient, genServer)
		case "gorilla":
			err = g.gorillaGenerate(component, plugin)
		case "chi":
			err = g.chiGenerate(component, plugin)
		case "openapiv3":
			err = g.openapiv3Generate(component, plugin)
		default:
			err = fmt.Errorf("unknown component: %s", component)
		}

		if err != nil {
			plugin.Error(err)
			return err
		}

	}

	if err = g.astGenerate(plugin); err != nil {
		plugin.Error(err)
		return err
	}

	return nil
}
