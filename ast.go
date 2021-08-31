package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
	tag_options "github.com/unistack-org/micro-proto/tag"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var astFields = make(map[string]map[string]map[string]*structtag.Tags) // map proto file with proto message ast struct

func (g *Generator) astFill(file *protogen.File, message *protogen.Message) error {
	for _, field := range message.Fields {
		if field.Desc.Options() == nil {
			continue
		}
		if !proto.HasExtension(field.Desc.Options(), tag_options.E_Tags) {
			continue
		}

		opts := proto.GetExtension(field.Desc.Options(), tag_options.E_Tags)
		if opts != nil {
			fpath := filepath.Join(g.tagPath, file.GeneratedFilenamePrefix+".pb.go")
			mp, ok := astFields[fpath]
			if !ok {
				mp = make(map[string]map[string]*structtag.Tags)
			}
			nmp, ok := mp[message.GoIdent.GoName]
			if !ok {
				nmp = make(map[string]*structtag.Tags)
			}
			tags, err := structtag.Parse(opts.(string))
			if err != nil {
				return err
			}
			nmp[field.GoName] = tags
			mp[message.GoIdent.GoName] = nmp
			astFields[fpath] = mp
		}
	}
	for _, nmessage := range message.Messages {
		if err := g.astFill(file, nmessage); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) astGenerate(plugin *protogen.Plugin) error {
	if g.tagPath == "" {
		return nil
	}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		for _, message := range file.Messages {
			if err := g.astFill(file, message); err != nil {
				return err
			}
		}
	}

	for file, mp := range astFields {
		fset := token.NewFileSet()
		pf, err := parser.ParseFile(fset, file, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			return err
		}

		r := retag{}
		f := func(n ast.Node) ast.Visitor {
			if r.err != nil {
				return nil
			}

			if v, ok := n.(*ast.TypeSpec); ok {
				r.fields = mp[v.Name.Name]
				return r
			}

			return nil
		}

		ast.Walk(structVisitor{f}, pf)

		if r.err != nil {
			return err
		}

		fp, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
		if err != nil {
			return err
		}
		if err = format.Node(fp, fset, pf); err != nil {
			fp.Close()
			return err
		}
		if err = fp.Close(); err != nil {
			return err
		}
	}

	return nil
}

type retag struct {
	err    error
	fields map[string]*structtag.Tags
}

func (v retag) Visit(n ast.Node) ast.Visitor {
	if v.err != nil {
		return nil
	}
	if f, ok := n.(*ast.Field); ok {
		if len(f.Names) == 0 {
			return nil
		}

		newTags := v.fields[f.Names[0].String()]
		if newTags == nil {
			return nil
		}
		if f.Tag == nil {
			f.Tag = &ast.BasicLit{
				Kind: token.STRING,
			}
		}

		oldTags, err := structtag.Parse(strings.Trim(f.Tag.Value, "`"))
		if err != nil {
			v.err = err
			return nil
		}
		for _, t := range newTags.Tags() {
			oldTags.Set(t)
		}

		f.Tag.Value = "`" + oldTags.String() + "`"

		return nil
	}

	return v
}

type structVisitor struct {
	visitor func(n ast.Node) ast.Visitor
}

func (v structVisitor) Visit(n ast.Node) ast.Visitor {
	if tp, ok := n.(*ast.TypeSpec); ok {
		if _, ok := tp.Type.(*ast.StructType); ok {
			ast.Walk(v.visitor(n), n)
			return nil // This will ensure this struct is no longer traversed
		}
	}
	return v
}
