package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/unistack-org/protoc-gen-micro/v3/assets"
	pgghelpers "github.com/unistack-org/protoc-gen-micro/v3/helpers"
)

type GenericTemplateBasedEncoder struct {
	templateDir    string
	service        *descriptor.ServiceDescriptorProto
	file           *descriptor.FileDescriptorProto
	enum           []*descriptor.EnumDescriptorProto
	debug          bool
	destinationDir string
}

type Ast struct {
	BuildDate      time.Time                          `json:"build-date"`
	BuildHostname  string                             `json:"build-hostname"`
	BuildUser      string                             `json:"build-user"`
	GoPWD          string                             `json:"go-pwd,omitempty"`
	PWD            string                             `json:"pwd"`
	Debug          bool                               `json:"debug"`
	DestinationDir string                             `json:"destination-dir"`
	File           *descriptor.FileDescriptorProto    `json:"file"`
	RawFilename    string                             `json:"raw-filename"`
	Filename       string                             `json:"filename"`
	TemplateDir    string                             `json:"template-dir"`
	Service        *descriptor.ServiceDescriptorProto `json:"service"`
	Enum           []*descriptor.EnumDescriptorProto  `json:"enum"`
}

func NewGenericServiceTemplateBasedEncoder(templateDir string, service *descriptor.ServiceDescriptorProto, file *descriptor.FileDescriptorProto, debug bool, destinationDir string) (e *GenericTemplateBasedEncoder) {
	e = &GenericTemplateBasedEncoder{
		service:        service,
		file:           file,
		templateDir:    templateDir,
		debug:          debug,
		destinationDir: destinationDir,
		enum:           file.GetEnumType(),
	}
	if debug {
		log.Printf("new encoder: file=%q service=%q template-dir=%q", file.GetName(), service.GetName(), templateDir)
	}
	pgghelpers.InitPathMap(file)

	return
}

func NewGenericTemplateBasedEncoder(templateDir string, file *descriptor.FileDescriptorProto, debug bool, destinationDir string) (e *GenericTemplateBasedEncoder) {
	e = &GenericTemplateBasedEncoder{
		service:        nil,
		file:           file,
		templateDir:    templateDir,
		enum:           file.GetEnumType(),
		debug:          debug,
		destinationDir: destinationDir,
	}
	if debug {
		log.Printf("new encoder: file=%q template-dir=%q", file.GetName(), templateDir)
	}
	pgghelpers.InitPathMap(file)

	return
}

func (e *GenericTemplateBasedEncoder) templates() ([]string, error) {
	filenames := []string{}

	if e.templateDir == "" {
		dir, err := assets.Assets.Open("/")
		if err != nil {
			return nil, fmt.Errorf("failed to open assets dir")
		}

		fi, err := dir.Readdir(-1)
		if err != nil {
			return nil, fmt.Errorf("failed to get assets files")
		}

		if debug {
			log.Printf("components to generate: %v", components)
		}

		for _, f := range fi {
			name := f.Name()
			skip := true
			if dname, err := base64.StdEncoding.DecodeString(name); err == nil {
				name = string(dname)
			}
			for _, component := range components {
				if component == "all" || strings.Contains(name, "_"+component+".pb.go.tmpl") {
					skip = false
				}
			}
			if skip {
				if debug {
					log.Printf("skip template %s", name)
				}
				continue
			}

			if f.IsDir() {
				continue
			}
			if filepath.Ext(name) != ".tmpl" {
				continue
			}
			if e.debug {
				log.Printf("new template: %q", name)
			}

			filenames = append(filenames, name)
		}

		return filenames, nil
	}

	err := filepath.Walk(e.templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}
		rel, err := filepath.Rel(e.templateDir, path)
		if err != nil {
			return err
		}
		if e.debug {
			log.Printf("new template: %q", rel)
		}

		filenames = append(filenames, rel)
		return nil
	})
	return filenames, err
}

func (e *GenericTemplateBasedEncoder) genAst(templateFilename string) (*Ast, error) {
	// prepare the ast passed to the template engine
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	goPwd := ""
	if os.Getenv("GOPATH") != "" {
		goPwd, err = filepath.Rel(os.Getenv("GOPATH")+"/src", pwd)
		if err != nil {
			return nil, err
		}
		if strings.Contains(goPwd, "../") {
			goPwd = ""
		}
	}
	ast := Ast{
		BuildDate:      time.Now(),
		BuildHostname:  hostname,
		BuildUser:      os.Getenv("USER"),
		PWD:            pwd,
		GoPWD:          goPwd,
		File:           e.file,
		TemplateDir:    e.templateDir,
		DestinationDir: e.destinationDir,
		RawFilename:    templateFilename,
		Filename:       "",
		Service:        e.service,
		Enum:           e.enum,
	}
	buffer := new(bytes.Buffer)

	unescaped, err := url.QueryUnescape(templateFilename)
	if err != nil {
		log.Printf("failed to unescape filepath %q: %v", templateFilename, err)
	} else {
		templateFilename = unescaped
	}

	tmpl, err := template.New("").Funcs(pgghelpers.ProtoHelpersFuncMap).Parse(templateFilename)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buffer, ast); err != nil {
		return nil, err
	}
	ast.Filename = buffer.String()
	return &ast, nil
}

func (e *GenericTemplateBasedEncoder) buildContent(templateFilename string) (string, string, error) {
	var tmpl *template.Template
	var err error

	if e.templateDir == "" {
		fs, err := assets.Assets.Open("/" + string(base64.StdEncoding.EncodeToString([]byte(templateFilename))))
		if err != nil {
			fs, err = assets.Assets.Open("/" + templateFilename)
		}
		if err != nil {
			return "", "", err
		}
		buf, err := ioutil.ReadAll(fs)
		if err != nil {
			return "", "", err
		}
		if err = fs.Close(); err == nil {
			tmpl, err = template.New("/" + templateFilename).Funcs(pgghelpers.ProtoHelpersFuncMap).Parse(string(buf))
		}
	} else {
		// initialize template engine
		fullPath := filepath.Join(e.templateDir, templateFilename)
		templateName := filepath.Base(fullPath)
		tmpl, err = template.New(templateName).Funcs(pgghelpers.ProtoHelpersFuncMap).ParseFiles(fullPath)
	}
	if err != nil {
		return "", "", err
	} else if tmpl == nil {
		return "", "", fmt.Errorf("template for %s is nil", templateFilename)
	}

	ast, err := e.genAst(templateFilename)
	if err != nil {
		return "", "", err
	}

	// generate the content
	buffer := new(bytes.Buffer)
	if err := tmpl.Execute(buffer, ast); err != nil {
		return "", "", err
	}

	return buffer.String(), ast.Filename, nil
}

func (e *GenericTemplateBasedEncoder) Files() []*plugin_go.CodeGeneratorResponse_File {
	templates, err := e.templates()
	if err != nil {
		log.Fatalf("cannot get templates from %q: %v", e.templateDir, err)
	}

	length := len(templates)
	files := make([]*plugin_go.CodeGeneratorResponse_File, 0, length)
	errChan := make(chan error, length)
	resultChan := make(chan *plugin_go.CodeGeneratorResponse_File, length)
	for _, templateFilename := range templates {
		go func(tmpl string) {
			var translatedFilename, content string
			content, translatedFilename, err = e.buildContent(tmpl)
			if err != nil {
				errChan <- err
				return
			}
			filename := translatedFilename[:len(translatedFilename)-len(".tmpl")]

			resultChan <- &plugin_go.CodeGeneratorResponse_File{
				Content: &content,
				Name:    &filename,
			}
		}(templateFilename)
	}
	for i := 0; i < length; i++ {
		select {
		case f := <-resultChan:
			files = append(files, f)
		case err = <-errChan:
			panic(err)
		}
	}
	return files
}
