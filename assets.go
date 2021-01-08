// +build dev

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/unistack-org/protoc-gen-micro/assets"
)

func main() {
	err := vfsgen.Generate(assets.Assets, vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Assets",
		Filename:     "assets/vfsdata.go",
	})
	if err != nil {
		log.Fatal(err)
	}
}
