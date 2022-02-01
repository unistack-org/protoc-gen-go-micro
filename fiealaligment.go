package main

import (
	"log"
	"os"

	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/singlechecker"
	"google.golang.org/protobuf/compiler/protogen"
)

func (g *Generator) fieldAlign(plugin *protogen.Plugin) error {
	if !g.fieldaligment {
		return nil
	}

	log.Printf("%v\n", []string{"fieldalignment", "-fix", g.tagPath})
	origArgs := os.Args
	os.Args = []string{"fieldalignment", "-fix", g.tagPath}
	singlechecker.Main(fieldalignment.Analyzer)
	os.Args = origArgs

	return nil
}
