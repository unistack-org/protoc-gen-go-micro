package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	flagSet           = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagDebug         = flagSet.Bool("debug", false, "debug output")
	flagStandalone    = flagSet.Bool("standalone", false, "generate file to standalone dir")
	flagFieldaligment = flagSet.Bool("fieldaligment", false, "align struct fields in generated code")
	flagComponents    = flagSet.String("components", "micro|rpc|http|client|server|openapiv3", "specify components to generate")
	flagTagPath       = flagSet.String("tag_path", "", "tag rewriting dir")
	flagOpenapiFile   = flagSet.String("openapi_file", "apidocs.swagger.json", "openapi file name")
	flagReflection    = flagSet.Bool("reflection", false, "enable server reflection support")
	flagHelp          = flagSet.Bool("help", false, "display help message")
)

func main() {
	opts := &protogen.Options{
		ParamFunc: flagSet.Set,
	}

	flagSet.Parse(os.Args[1:])

	if *flagHelp {
		flagSet.PrintDefaults()
		return
	}

	g := &Generator{}

	opts.Run(g.Generate)
}

type Generator struct {
	components    string
	standalone    bool
	debug         bool
	fieldaligment bool
	tagPath       string
	openapiFile   string
	reflection    bool
	plugin        *protogen.Plugin
}

func (g *Generator) Generate(plugin *protogen.Plugin) error {
	var err error

	g.plugin = plugin
	g.standalone = *flagStandalone
	g.debug = *flagDebug
	g.components = *flagComponents
	g.fieldaligment = *flagFieldaligment
	g.tagPath = *flagTagPath
	g.openapiFile = *flagOpenapiFile
	g.reflection = *flagReflection
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	var genClient bool
	var genServer bool
	var genNone bool

	if strings.Contains(g.components, "server") {
		genServer = true
	}
	if strings.Contains(g.components, "client") {
		genClient = true
	}
	if strings.Contains(g.components, "none") {
		genNone = true
	}
	if strings.Contains(g.components, "rpc") || strings.Contains(g.components, "http") {
		if !genServer && !genClient && !genNone {
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
		case "grpc", "drpc", "rpc":
			err = g.rpcGenerate(component, plugin, genClient, genServer)
		case "gorilla":
			err = g.gorillaGenerate(component, plugin)
		case "chi":
			err = g.chiGenerate(component, plugin)
		case "openapiv3":
			err = g.openapiv3Generate(component, plugin)
		case "none":
			break
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

	if err = g.fieldAlign(plugin); err != nil {
		plugin.Error(err)
		return err
	}

	return nil
}
