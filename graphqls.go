package main

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/danielvladco/go-proto-gql/pkg/generator"
	"github.com/vektah/gqlparser/v2/formatter"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func (g *Generator) graphqlsGenerate(plugin *protogen.Plugin) error {
	descs, err := generator.CreateDescriptorsFromProto(plugin.Request)
	if err != nil {
		return err
	}

	gqlDesc, err := generator.NewSchemas(descs, true, true, plugin)
	if err != nil {
		return err
	}

	var outFiles []*pluginpb.CodeGeneratorResponse_File

	for _, schema := range gqlDesc {
		buf := &bytes.Buffer{}
		formatter.NewFormatter(buf).FormatSchema(schema.AsGraphql())
		protoFileName := schema.FileDescriptors[0].GetName()

		outFiles = append(outFiles, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(resolveGraphqlFilename(protoFileName)),
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

func resolveGraphqlFilename(protoFileName string) string {
	gqlFileName := "schema." + "graphqls"
	absProtoFileName, err := filepath.Abs(protoFileName)
	if err == nil {
		protoDirSlice := strings.Split(filepath.Dir(absProtoFileName), string(filepath.Separator))
		if len(protoDirSlice) > 0 {
			gqlFileName = protoDirSlice[len(protoDirSlice)-1] + "." + "graphqls"
		}
	}
	protoDir, _ := path.Split(protoFileName)
	return path.Join(protoDir, gqlFileName)
}
