package main

import (
	"flag"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	flags          flag.FlagSet
	flagDebug      *bool
	flagComponents *string
	flagPaths      *string
	flagModule     *string
)

func init() {
	flagDebug = flags.Bool("debug", false, "")
	flagComponents = flags.String("components", "micro", "")
	flagPaths = flag.String("paths", "", "")
	flagModule = flag.String("module", "", "")
}

func main() {
	opts := &protogen.Options{
		ParamFunc: flags.Set,
	}

	g := &Generator{}

	opts.Run(g.Generate)
}

type Generator struct {
}

func (g *Generator) Generate(plugin *protogen.Plugin) error {
	var err error

	// Protoc passes a slice of File structs for us to process
	for _, component := range strings.Split(*flagComponents, "|") {
		switch component {
		case "micro":
			err = g.microGenerate(component, plugin)
		case "http":
			err = g.httpGenerate(component, plugin)
		case "grpc", "rpc":
			err = g.rpcGenerate(component, plugin)
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
