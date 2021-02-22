package main

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func (g *Generator) openapiGenerate(component string, plugin *protogen.Plugin) error {
	for _, file := range plugin.Files {

		if !file.Generate {
			continue
		}

		// 1. Initialise a buffer to hold the generated code
		var buf bytes.Buffer

		// 2. Write the package name
		pkg := fmt.Sprintf("package %s", file.GoPackageName)
		buf.Write([]byte(pkg))

		// 3. For each message add our Foo() method
		for _, msg := range file.Proto.MessageType {
			buf.Write([]byte(fmt.Sprintf(`func (x %s) Foo() string {
return "bar"}`, *msg.Name)))
		}

		// 4. Specify the output filename, in this case test.foo.go
		filename := file.GeneratedFilenamePrefix + ".foo.go"
		file := plugin.NewGeneratedFile(filename, ".")

		// 5. Pass the data from our buffer to the plugin file struct
		file.Write(buf.Bytes())
	}

	return nil
}
