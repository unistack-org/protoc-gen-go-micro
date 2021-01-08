// +build dev

package assets

import (
	"go/build"
	"log"
	"net/http"
)

func importPathToDir(importPath string) string {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	return p.Dir
}

// Assets contains the project's assets.
var Assets http.FileSystem = http.Dir("./templates")
