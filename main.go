package main

import (
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/generator"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	ggdescriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/unistack-org/protoc-gen-micro/assets"
	pgghelpers "github.com/unistack-org/protoc-gen-micro/helpers"
	"google.golang.org/protobuf/proto"
)

var (
	registry *ggdescriptor.Registry // some helpers need access to registry
)

const (
	boolTrue  = "true"
	boolFalse = "false"
)

func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err = proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// Parse parameters
	var (
		templateDir       = ""
		templateRepo      = ""
		destinationDir    = "."
		debug             = false
		all               = false
		singlePackageMode = false
		fileMode          = false
		components        = []string{"micro", "grpc"}
	)
	if parameter := g.Request.GetParameter(); parameter != "" {
		for _, param := range strings.Split(parameter, ",") {
			parts := strings.Split(param, "=")
			switch parts[0] {
			case "template_dir":
				templateDir = parts[1]
			case "destination_dir":
				destinationDir = parts[1]
			case "single-package-mode":
				switch strings.ToLower(parts[1]) {
				case boolTrue, "t":
					singlePackageMode = true
				case boolFalse, "f":
				default:
					log.Printf("Err: invalid value for single-package-mode: %q", parts[1])
				}
			case "debug":
				switch strings.ToLower(parts[1]) {
				case boolTrue, "t":
					debug = true
				case boolFalse, "f":
				default:
					log.Printf("Err: invalid value for debug: %q", parts[1])
				}
			case "all":
				switch strings.ToLower(parts[1]) {
				case boolTrue, "t":
					all = true
				case boolFalse, "f":
				default:
					log.Printf("Err: invalid value for debug: %q", parts[1])
				}
			case "file-mode":
				switch strings.ToLower(parts[1]) {
				case boolTrue, "t":
					fileMode = true
				case boolFalse, "f":
				default:
					log.Printf("Err: invalid value for file-mode: %q", parts[1])
				}
			case "template_repo":
				_, err := url.Parse(parts[1])
				if err != nil {
					log.Printf("Err: invalid value for template_repo: %q", parts[1])
				}
				templateRepo = parts[1]
			case "components":
				_, err := url.Parse(parts[1])
				if err != nil {
					log.Printf("Err: invalid value for components: %q", parts[1])
				}
				components = strings.Split(parts[1], "|")
			case "paths":
				// TODO: handle paths=source_relative
			default:
				log.Printf("Err: unknown parameter: %q", param)
			}
		}
	}

	tmplMap := make(map[string]*plugin_go.CodeGeneratorResponse_File)
	concatOrAppend := func(file *plugin_go.CodeGeneratorResponse_File) {
		if val, ok := tmplMap[file.GetName()]; ok {
			*val.Content += file.GetContent()
		} else {
			tmplMap[file.GetName()] = file
			g.Response.File = append(g.Response.File, file)
		}
	}

	if singlePackageMode {
		registry = ggdescriptor.NewRegistry()
		pgghelpers.SetRegistry(registry)
		if err = registry.Load(g.Request); err != nil {
			g.Error(err, "registry: failed to load the request")
		}
	}

	if templateDir == "" || templateRepo != "" {
		if templateDir, err = ioutil.TempDir("", "gen-*"); err != nil {
			g.Error(err, "failed to create tmp dir")
		}
		defer func() {
			if err := os.RemoveAll(templateDir); err != nil {
				g.Error(err, "failed to remove tmp dir")
			}
		}()

		if templateRepo != "" {
			if err = clone(templateRepo, templateDir); err != nil {
				g.Error(err, "failed to clone repo")
			}
		} else {
			dir, err := assets.Assets.Open("/")
			if err != nil {
				g.Error(err, "failed to open assets dir")
			}
			fi, err := dir.Readdir(-1)
			if err != nil {
				g.Error(err, "failed to get assets files")
			}

			for _, f := range fi {
				skip := true
				for _, component := range components {
					if component == "all" || strings.Contains(f.Name(), "_"+component+".pb.go") {
						skip = false
					}
				}
				if skip {
					if debug {
						log.Printf("skip template %s", f.Name())
					}
					continue
				}
				if debug {
					log.Printf("copy template %s", f.Name())
				}
				fpath := filepath.Join(templateDir, f.Name())
				if err = os.MkdirAll(filepath.Dir(fpath), os.FileMode(0755)); err != nil {
					g.Error(err, "failed to create nested dir")
				}
				if f.IsDir() {
					continue
				}
				fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, f.Mode())
				if err != nil {
					g.Error(err, "failed to create template file")
				}
				fs, err := assets.Assets.Open(f.Name())
				if err != nil {
					g.Error(err, "failed to open template file")
				}
				if _, err = io.Copy(fd, fs); err != nil {
					fd.Close()
					fs.Close()
					g.Error(err, "failed to copy template file")
				}
				if err = fd.Close(); err != nil {
					g.Error(err, "failed to flush template file")
				}
				if err = fs.Close(); err != nil {
					g.Error(err, "failed to flush template file")
				}
			}
		}
	}

	// Generate the encoders
	for _, file := range g.Request.GetProtoFile() {
		if all {
			if singlePackageMode {
				if _, err = registry.LookupFile(file.GetName()); err != nil {
					g.Error(err, "registry: failed to lookup file %q", file.GetName())
				}
			}
			encoder := NewGenericTemplateBasedEncoder(templateDir, file, debug, destinationDir)
			for _, tmpl := range encoder.Files() {
				concatOrAppend(tmpl)
			}

			continue
		}

		if fileMode {
			if s := file.GetService(); s != nil && len(s) > 0 {
				encoder := NewGenericTemplateBasedEncoder(templateDir, file, debug, destinationDir)
				for _, tmpl := range encoder.Files() {
					concatOrAppend(tmpl)
				}
			}

			continue
		}

		for _, service := range file.GetService() {
			encoder := NewGenericServiceTemplateBasedEncoder(templateDir, service, file, debug, destinationDir)
			for _, tmpl := range encoder.Files() {
				concatOrAppend(tmpl)
			}
		}
	}

	// Generate the protobufs
	g.GenerateAllFiles()

	for _, f := range g.Response.File {
		fdata, err := format.Source([]byte(*f.Content))
		if err != nil {
			g.Error(err, fmt.Sprintf("failed to format output: %s", *f.Content))
		}
		*f.Content = string(fdata)
	}

	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
