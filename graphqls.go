package main

import (
	"bytes"
	"os"

	"github.com/vektah/gqlparser/v2/formatter"
	generator "go.unistack.org/protoc-gen-go-micro/v3/graphql"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func (g *Generator) graphqlsGenerate(plugin *protogen.Plugin) error {
	descs, err := generator.CreateDescriptorsFromProto(plugin.Request)
	if err != nil {
		return err
	}

	gqlDesc, err := generator.NewSchemas(descs, false, false, plugin)
	if err != nil {
		return err
	}

	var outFiles []*pluginpb.CodeGeneratorResponse_File

	for _, schema := range gqlDesc {
		buf := &bytes.Buffer{}
		formatter.NewFormatter(buf).FormatSchema(schema.AsGraphql())

		outFiles = append(outFiles, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(g.graphqlFile),
			Content: proto.String(buf.String()),
		})
	}

	res := &pluginpb.CodeGeneratorResponse{
		File:              outFiles,
		SupportedFeatures: proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)),
	}

	if err != nil {
		res.Error = proto.String(err.Error())
	}

	out, err := proto.Marshal(res)
	if err != nil {
		return err
	}

	if _, err := os.Stdout.Write(out); err != nil {
		return err
	}

	return nil
}
