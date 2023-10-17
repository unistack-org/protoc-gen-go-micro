package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	api_options "go.unistack.org/micro-proto/v3/api"
	v2 "go.unistack.org/micro-proto/v3/openapiv2"
	v3 "go.unistack.org/micro-proto/v3/openapiv3"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var httpMethodMap = map[string]string{
	http.MethodGet:     "MethodGet",
	http.MethodHead:    "MethodHead",
	http.MethodPost:    "MethodPost",
	http.MethodPut:     "MethodPut",
	http.MethodPatch:   "MethodPatch",
	http.MethodDelete:  "MethodDelete",
	http.MethodConnect: "MethodConnect",
	http.MethodOptions: "MethodOptions",
	http.MethodTrace:   "MethodTrace",
}

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func (g *Generator) generateServiceClient(gfile *protogen.GeneratedFile, file *protogen.File, service *protogen.Service) {
	serviceName := service.GoName
	// if rule, ok := getMicroApiService(service); ok {
	//		gfile.P("// client wrappers ", strings.Join(rule.ClientWrappers, ", "))
	//	}
	gfile.P("type ", unexport(serviceName), "Client struct {")
	gfile.P("c ", microClientPackage.Ident("Client"))
	gfile.P("name string")
	gfile.P("}")

	if g.standalone {
		gfile.P("func New", serviceName, "Client(name string, c ", microClientPackage.Ident("Client"), ") ", file.GoImportPath.Ident(serviceName), "Client {")
	} else {
		gfile.P("func New", serviceName, "Client(name string, c ", microClientPackage.Ident("Client"), ") ", serviceName, "Client {")
	}
	gfile.P("return &", unexport(serviceName), "Client{c: c, name: name}")
	gfile.P("}")
	gfile.P()
}

func (g *Generator) generateServiceClientMethods(gfile *protogen.GeneratedFile, service *protogen.Service, component string) {
	serviceName := service.GoName
	for _, method := range service.Methods {
		methodName := fmt.Sprintf("%s.%s", serviceName, method.GoName)
		if component == "drpc" {
			methodName = fmt.Sprintf("%s.%s", method.Parent.Desc.FullName(), method.Desc.Name())
		}
		g.generateClientFuncSignature(gfile, serviceName, method)

		if component == "http" && method.Desc.Options() != nil {
			if proto.HasExtension(method.Desc.Options(), v2.E_Openapiv2Operation) {
				opts := proto.GetExtension(method.Desc.Options(), v2.E_Openapiv2Operation)
				if opts != nil {
					r := opts.(*v2.Operation)
					if r.Responses == nil {
						goto labelMethod
					}
					gfile.P("errmap := make(map[string]interface{}, ", len(r.Responses.ResponseCode), ")")
					for _, rsp := range r.Responses.ResponseCode {
						if schema := rsp.Value.GetJsonReference(); schema != nil {
							xref := schema.XRef
							if strings.HasPrefix(xref, "."+string(service.Desc.ParentFile().Package())+".") {
								xref = strings.TrimPrefix(xref, "."+string(service.Desc.ParentFile().Package())+".")
							}
							if xref[0] == '.' {
								xref = xref[1:]
							}
							switch xref {
							case "micro.codec.Frame":
								gfile.P(`errmap["`, rsp.Name, `"] = &`, microCodecPackage.Ident("Frame"), "{}")
							case "micro.errors.Error":
								gfile.P(`errmap["`, rsp.Name, `"] = &`, microErrorsPackage.Ident("Error"), "{}")
							default:
								ident, err := g.getGoIdentByXref(strings.TrimPrefix(schema.XRef, "."))
								if err != nil {
									log.Printf("cant find message by ref %s\n", schema.XRef)
									continue
								}
								gfile.P(`errmap["`, rsp.Name, `"] = &`, gfile.QualifiedGoIdent(ident), "{}")
							}
						}
					}
				}
				gfile.P("opts = append(opts,")
				gfile.P(microClientHttpPackage.Ident("ErrorMap"), "(errmap),")
				gfile.P(")")
			}
			if proto.HasExtension(method.Desc.Options(), v3.E_Openapiv3Operation) {
				opts := proto.GetExtension(method.Desc.Options(), v3.E_Openapiv3Operation)
				if opts != nil {
					r := opts.(*v3.Operation)
					if r.Responses == nil {
						goto labelMethod
					}
					resps := r.Responses.ResponseOrReference
					if r.Responses.GetDefault() != nil {
						resps = append(resps, &v3.NamedResponseOrReference{Name: "default", Value: r.Responses.GetDefault()})
					}
					gfile.P("errmap := make(map[string]interface{}, ", len(resps), ")")
					for _, rsp := range resps {
						if schema := rsp.Value.GetReference(); schema != nil {
							xref := schema.XRef
							if strings.HasPrefix(xref, "."+string(service.Desc.ParentFile().Package())+".") {
								xref = strings.TrimPrefix(xref, "."+string(service.Desc.ParentFile().Package())+".")
							}
							if xref[0] == '.' {
								xref = xref[1:]
							}
							switch xref {
							case "micro.codec.Frame":
								gfile.P(`errmap["`, rsp.Name, `"] = &`, microCodecPackage.Ident("Frame"), "{}")
							case "micro.errors.Error":
								gfile.P(`errmap["`, rsp.Name, `"] = &`, microErrorsPackage.Ident("Error"), "{}")
							default:
								ident, err := g.getGoIdentByXref(strings.TrimPrefix(schema.XRef, "."))
								if err != nil {
									log.Printf("cant find message by ref %s\n", schema.XRef)
									continue
								}
								gfile.P(`errmap["`, rsp.Name, `"] = &`, gfile.QualifiedGoIdent(ident), "{}")
							}
						}
					}
				}
				gfile.P("opts = append(opts,")
				gfile.P(microClientHttpPackage.Ident("ErrorMap"), "(errmap),")
				gfile.P(")")
			}

		labelMethod:
			if proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
				gfile.P("opts = append(opts,")
				endpoints, _ := generateEndpoints(method)
				path, method, body := getEndpoint(endpoints[0])
				if vmethod, ok := httpMethodMap[method]; ok {
					gfile.P(microClientHttpPackage.Ident("Method"), `(`, httpPackage.Ident(vmethod), `),`)
				} else {
					gfile.P(microClientHttpPackage.Ident("Method"), `("`, method, `"),`)
				}
				gfile.P(microClientHttpPackage.Ident("Path"), `("`, path, `"),`)
				if body != "" {
					gfile.P(microClientHttpPackage.Ident("Body"), `("`, body, `"),`)
				}
				gfile.P(")")
			}

			parameters := make(map[string]map[string]string)
			// Build a list of header parameters.
			e2opt := proto.GetExtension(method.Desc.Options(), v2.E_Openapiv2Operation)
			if e2opt != nil && e2opt != v2.E_Openapiv2Operation.InterfaceOf(v2.E_Openapiv2Operation.Zero()) {
				opt := e2opt.(*v2.Operation)
				for _, paramOrRef := range opt.Parameters {
					parameter := paramOrRef.GetParameter()
					// NonBodyParameter()
					if parameter == nil {
						continue
					}
					nonBodyParameter := parameter.GetNonBodyParameter()
					if nonBodyParameter == nil {
						continue
					}
					headerParameter := nonBodyParameter.GetHeaderParameterSubSchema()
					if headerParameter.In != "header" && headerParameter.In != "cookie" {
						continue
					}
					in, ok := parameters[headerParameter.In]
					if !ok {
						in = make(map[string]string)
						parameters[headerParameter.In] = in
					}
					in[headerParameter.Name] = fmt.Sprintf("%v", headerParameter.Required)
				}
			}
			e3opt := proto.GetExtension(method.Desc.Options(), v3.E_Openapiv3Operation)
			if e3opt != nil && e3opt != v3.E_Openapiv3Operation.InterfaceOf(v3.E_Openapiv3Operation.Zero()) {
				opt := e3opt.(*v3.Operation)
				for _, paramOrRef := range opt.Parameters {

					parameter := paramOrRef.GetParameter()
					if parameter == nil {
						continue
					}
					if parameter.In != "header" && parameter.In != "cookie" {
						continue
					}
					in, ok := parameters[parameter.In]
					if !ok {
						in = make(map[string]string)
						parameters[parameter.In] = in
					}
					in[parameter.Name] = fmt.Sprintf("%v", parameter.Required)
				}
			}

			if len(parameters) > 0 {
				gfile.P("opts = append(opts,")
				for pk, pv := range parameters {
					params := make([]string, 0, len(pv)/2)
					for k, v := range pv {
						params = append(params, k, v)
					}
					gfile.P(microClientHttpPackage.Ident(strings.Title(pk)), `("`, strings.Join(params, `" ,"`), `"),`)
				}
				gfile.P(")")
			}
		}

		if rule, ok := getMicroApiMethod(method); ok {
			if rule.Timeout != "" {
				td, err := time.ParseDuration(rule.Timeout)
				if err != nil {
					log.Printf("parse duration error %s\n", err.Error())
				} else {
					gfile.P("td := ", timePackage.Ident("Duration"), "(", td.Nanoseconds(), ")")
					gfile.P("opts = append(opts, ", microClientPackage.Ident("WithRequestTimeout"), "(td))")
				}
			}
		}

		if !method.Desc.IsStreamingServer() && !method.Desc.IsStreamingClient() {
			gfile.P("rsp := &", gfile.QualifiedGoIdent(method.Output.GoIdent), "{}")
			gfile.P(`err := c.c.Call(ctx, c.c.NewRequest(c.name, "`, methodName, `", req), rsp, opts...)`)
			gfile.P("if err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
			gfile.P("return rsp, nil")
			gfile.P("}")
			gfile.P()
			continue
		}

		gfile.P(`stream, err := c.c.Stream(ctx, c.c.NewRequest(c.name, "`, methodName, `", &`, gfile.QualifiedGoIdent(method.Input.GoIdent), `{}), opts...)`)
		gfile.P("if err != nil {")
		gfile.P("return nil, err")
		gfile.P("}")

		if !method.Desc.IsStreamingClient() {
			gfile.P("if err := stream.Send(req); err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
		}
		gfile.P("return &", unexport(serviceName), "Client", method.GoName, "{stream}, nil")
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			gfile.P("type ", unexport(serviceName), "Client", method.GoName, " struct {")
			gfile.P("stream ", microClientPackage.Ident("Stream"))
			gfile.P("}")
		}

		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") CloseAndRecv() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error) {")
			gfile.P("msg := &", gfile.QualifiedGoIdent(method.Output.GoIdent), "{}")
			gfile.P("err := s.RecvMsg(msg)")
			gfile.P("if err == nil {")
			gfile.P("err = s.Close()")
			gfile.P("}")
			gfile.P("if err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
			gfile.P("return msg, nil")
			gfile.P("}")
		}

		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") Close() error {")
		gfile.P("return s.stream.Close()")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") CloseSend() error {")
		gfile.P("return s.stream.CloseSend()")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") Context() ", contextPackage.Ident("Context"), " {")
		gfile.P("return s.stream.Context()")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") SendMsg(msg interface{}) error {")
		gfile.P("return s.stream.Send(msg)")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") RecvMsg(msg interface{}) error {")
		gfile.P("return s.stream.Recv(msg)")
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingClient() {
			gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") Header() ", microMetadataPackage.Ident("Metadata"), "{")
			gfile.P("return s.stream.Response().Header()")
			gfile.P("}")
			gfile.P()

			gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") Send(msg *", gfile.QualifiedGoIdent(method.Input.GoIdent), ") error {")
			gfile.P("return s.stream.Send(msg)")
			gfile.P("}")
			gfile.P()
		}

		if method.Desc.IsStreamingServer() {
			gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") Recv() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error) {")
			gfile.P("msg := &", gfile.QualifiedGoIdent(method.Output.GoIdent), "{}")
			gfile.P("if err := s.stream.Recv(msg); err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
			gfile.P("return msg, nil")
			gfile.P("}")
			gfile.P()
		}
	}
}

func (g *Generator) generateServiceServer(gfile *protogen.GeneratedFile, file *protogen.File, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("type ", unexport(serviceName), "Server struct {")
	if g.standalone {
		gfile.P(file.GoImportPath.Ident(serviceName), "Server")
	} else {
		gfile.P(serviceName, "Server")
	}
	gfile.P("}")
	gfile.P()
}

func (g *Generator) generateServiceServerMethods(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	for _, method := range service.Methods {
		generateServerFuncSignature(gfile, serviceName, method, true)
		if rule, ok := getMicroApiMethod(method); ok {
			if rule.Timeout != "" {
				td, err := time.ParseDuration(rule.Timeout)
				if err != nil {
					log.Printf("parse duration error %s\n", err.Error())
				} else {
					gfile.P("var cancel ", contextPackage.Ident("CancelFunc"))
					gfile.P("td := ", timePackage.Ident("Duration"), "(", td.Nanoseconds(), ")")
					gfile.P("ctx, cancel = ", contextPackage.Ident("WithTimeout"), "(ctx, ", "td", ")")
					gfile.P("defer cancel()")
				}
			}
		}
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			if !method.Desc.IsStreamingClient() {
				gfile.P("msg := &", gfile.QualifiedGoIdent(method.Input.GoIdent), "{}")
				gfile.P("if err := stream.Recv(msg); err != nil {")
				gfile.P("return err")
				gfile.P("}")
				gfile.P("return h.", serviceName, "Server.", method.GoName, "(ctx, msg, &", unexport(serviceName), method.GoName, "Stream{stream})")
			} else {
				gfile.P("return h.", serviceName, "Server.", method.GoName, "(ctx, &", unexport(serviceName), method.GoName, "Stream{stream})")
			}
		} else {
			parameters := make(map[string]map[string]string)
			// Build a list of header parameters.
			e2opt := proto.GetExtension(method.Desc.Options(), v2.E_Openapiv2Operation)
			if e2opt != nil && e2opt != v2.E_Openapiv2Operation.InterfaceOf(v2.E_Openapiv2Operation.Zero()) {
				opt := e2opt.(*v2.Operation)
				for _, paramOrRef := range opt.Parameters {
					parameter := paramOrRef.GetParameter()
					// NonBodyParameter()
					if parameter == nil {
						continue
					}
					nonBodyParameter := parameter.GetNonBodyParameter()
					if nonBodyParameter == nil {
						continue
					}
					headerParameter := nonBodyParameter.GetHeaderParameterSubSchema()
					if headerParameter.In != "header" && headerParameter.In != "cookie" {
						continue
					}
					in, ok := parameters[headerParameter.In]
					if !ok {
						in = make(map[string]string)
						parameters[headerParameter.In] = in
					}
					in[headerParameter.Name] = fmt.Sprintf("%v", headerParameter.Required)
				}
			}
			e3opt := proto.GetExtension(method.Desc.Options(), v3.E_Openapiv3Operation)
			if e3opt != nil && e3opt != v3.E_Openapiv3Operation.InterfaceOf(v3.E_Openapiv3Operation.Zero()) {
				opt := e3opt.(*v3.Operation)
				for _, paramOrRef := range opt.Parameters {
					parameter := paramOrRef.GetParameter()
					if parameter == nil {
						continue
					}
					if parameter.In != "header" && parameter.In != "cookie" {
						continue
					}
					in, ok := parameters[parameter.In]
					if !ok {
						in = make(map[string]string)
						parameters[parameter.In] = in
					}
					in[parameter.Name] = fmt.Sprintf("%v", parameter.Required)
				}
			}

			if len(parameters) > 0 {
				gfile.P(microServerHttpPackage.Ident("FillRequest"), `(ctx, req, `)
				for pk, pv := range parameters {
					params := make([]string, 0, len(pv)/2)
					for k, v := range pv {
						params = append(params, k, v)
					}
					gfile.P(microServerHttpPackage.Ident(strings.Title(pk)), `("`, strings.Join(params, `" ,"`), `"),`)
				}
				gfile.P(")")
			}
			gfile.P("return h.", serviceName, "Server.", method.GoName, "(ctx, req, rsp)")
		}
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			gfile.P("type ", unexport(serviceName), method.GoName, "Stream struct {")
			gfile.P("stream ", microServerPackage.Ident("Stream"))
			gfile.P("}")
		}

		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("func (s *", unexport(serviceName), method.GoName, "Stream) SendAndClose(msg *", gfile.QualifiedGoIdent(method.Output.GoIdent), ") error {")
			gfile.P("err := s.SendMsg(msg)")
			gfile.P("if err == nil {")
			gfile.P("err = s.stream.Close()")
			gfile.P("}")
			gfile.P("return err")
			gfile.P("}")
		}

		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			gfile.P("func (s *", unexport(serviceName), method.GoName, "Stream) Close() error {")
			gfile.P("return s.stream.Close()")
			gfile.P("}")
			gfile.P()

			gfile.P("func (s *", unexport(serviceName), method.GoName, "Stream) Context() ", contextPackage.Ident("Context"), " {")
			gfile.P("return s.stream.Context()")
			gfile.P("}")
			gfile.P()

			gfile.P("func (s *", unexport(serviceName), method.GoName, "Stream) SendMsg(msg interface{}) error {")
			gfile.P("return s.stream.Send(msg)")
			gfile.P("}")
			gfile.P()

			gfile.P("func (s *", unexport(serviceName), method.GoName, "Stream) RecvMsg(msg interface{}) error {")
			gfile.P("return s.stream.Recv(msg)")
			gfile.P("}")
			gfile.P()

			if method.Desc.IsStreamingServer() {
				gfile.P("func (s *", unexport(serviceName), method.GoName, "Stream) Send(msg *", gfile.QualifiedGoIdent(method.Output.GoIdent), ") error {")
				gfile.P("return s.stream.Send(msg)")
				gfile.P("}")
				gfile.P()
			}

			if method.Desc.IsStreamingClient() {
				gfile.P("func (s *", unexport(serviceName), method.GoName, "Stream) Recv() (*", gfile.QualifiedGoIdent(method.Input.GoIdent), ", error) {")
				gfile.P("msg := &", gfile.QualifiedGoIdent(method.Input.GoIdent), "{}")
				gfile.P("if err := s.stream.Recv(msg); err != nil {")
				gfile.P("return nil, err")
				gfile.P("}")
				gfile.P("return msg, nil")
				gfile.P("}")
				gfile.P()
			}
		}

	}
}

func (g *Generator) generateServiceRegister(gfile *protogen.GeneratedFile, file *protogen.File, service *protogen.Service, component string) {
	serviceName := service.GoName
	if g.standalone {
		gfile.P("func Register", serviceName, "Server(s ", microServerPackage.Ident("Server"), ", sh ", file.GoImportPath.Ident(serviceName), "Server, opts ...", microServerPackage.Ident("HandlerOption"), ") error {")
	} else {
		gfile.P("func Register", serviceName, "Server(s ", microServerPackage.Ident("Server"), ", sh ", serviceName, "Server, opts ...", microServerPackage.Ident("HandlerOption"), ") error {")
	}
	gfile.P("type ", unexport(serviceName), " interface {")
	for _, method := range service.Methods {
		generateServerSignature(gfile, serviceName, method, true)
	}
	gfile.P("}")
	gfile.P("type ", serviceName, " struct {")
	gfile.P(unexport(serviceName))
	gfile.P("}")
	gfile.P("h := &", unexport(serviceName), "Server{sh}")
	gfile.P("var nopts []", microServerPackage.Ident("HandlerOption"))
	if component == "http" {
		//	if g.standalone {
		//		gfile.P("nopts = append(nopts, ", microServerHttpPackage.Ident("HandlerEndpoints"), "(", file.GoImportPath.Ident(serviceName), "ServerEndpoints))")
		//	} else {
		gfile.P("nopts = append(nopts, ", microServerHttpPackage.Ident("HandlerEndpoints"), "(", serviceName, "ServerEndpoints))")
		//	}
	}
	gfile.P("return s.Handle(s.NewHandler(&", serviceName, "{h}, append(nopts, opts...)...))")
	gfile.P("}")
}

func generateServerFuncSignature(gfile *protogen.GeneratedFile, serviceName string, method *protogen.Method, private bool) {
	args := append([]interface{}{},
		"func (h *", unexport(serviceName), "Server) ", method.GoName,
		"(ctx ", contextPackage.Ident("Context"),
	)
	if private && (method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer()) {
		args = append(args, ", stream ", microServerPackage.Ident("Stream"))
	} else {
		if !method.Desc.IsStreamingClient() {
			args = append(args, ", req *", gfile.QualifiedGoIdent(method.Input.GoIdent))
		}
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			args = append(args, ", stream ", serviceName, "_", method.GoName, "Stream")
		}
		if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			args = append(args, ", rsp *", gfile.QualifiedGoIdent(method.Output.GoIdent))
		}
	}
	args = append(args, ") error {")
	gfile.P(args...)
}

func generateServerSignature(gfile *protogen.GeneratedFile, serviceName string, method *protogen.Method, private bool) {
	args := append([]interface{}{},
		method.GoName,
		"(ctx ", contextPackage.Ident("Context"),
	)
	if private && (method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer()) {
		args = append(args, ", stream ", microServerPackage.Ident("Stream"))
	} else {
		if !method.Desc.IsStreamingClient() {
			args = append(args, ", req *", gfile.QualifiedGoIdent(method.Input.GoIdent))
		}
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			args = append(args, ", stream ", serviceName, "_", method.GoName, "Stream")
		}
		if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			args = append(args, ", rsp *", gfile.QualifiedGoIdent(method.Output.GoIdent))
		}
	}
	args = append(args, ") error")
	gfile.P(args...)
}

func (g *Generator) generateClientFuncSignature(gfile *protogen.GeneratedFile, serviceName string, method *protogen.Method) {
	args := append([]interface{}{},
		"func (c *",
		unexport(serviceName),
		"Client) ",
		method.GoName,
		"(ctx ", contextPackage.Ident("Context"), ", ",
	)
	if !method.Desc.IsStreamingClient() {
		args = append(args, "req *", gfile.QualifiedGoIdent(method.Input.GoIdent), ", ")
	}
	args = append(args, "opts ...", microClientPackage.Ident("CallOption"), ") (")
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		args = append(args, "*", gfile.QualifiedGoIdent(method.Output.GoIdent))
	} else {
		args = append(args, serviceName, "_", method.GoName, "Client")
	}
	args = append(args, ", error) {")
	gfile.P(args...)
}

func generateClientSignature(gfile *protogen.GeneratedFile, serviceName string, method *protogen.Method) {
	args := append([]interface{}{},
		method.GoName,
		"(ctx ", contextPackage.Ident("Context"), ", ",
	)
	if !method.Desc.IsStreamingClient() {
		args = append(args, "req *", gfile.QualifiedGoIdent(method.Input.GoIdent), ", ")
	}
	args = append(args, "opts ...", microClientPackage.Ident("CallOption"), ") (")
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		args = append(args, "*", gfile.QualifiedGoIdent(method.Output.GoIdent))
	} else {
		args = append(args, serviceName, "_", method.GoName, "Client")
	}
	args = append(args, ", error)")
	gfile.P(args...)
}

func (g *Generator) generateServiceClientInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("type ", serviceName, "Client interface {")
	for _, method := range service.Methods {
		generateClientSignature(gfile, serviceName, method)
	}
	gfile.P("}")
	gfile.P()
}

func (g *Generator) generateServiceServerInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("type ", serviceName, "Server interface {")
	for _, method := range service.Methods {
		generateServerSignature(gfile, serviceName, method, false)
	}
	gfile.P("}")
	gfile.P()
}

func (g *Generator) generateServiceClientStreamInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	for _, method := range service.Methods {
		if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			continue
		}
		methodName := method.GoName
		gfile.P("type ", serviceName, "_", methodName, "Client interface {")
		gfile.P("Context() ", contextPackage.Ident("Context"))
		gfile.P("SendMsg(msg interface{}) error")
		gfile.P("RecvMsg(msg interface{}) error")
		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("CloseAndRecv() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error)")
			gfile.P("CloseSend() error")
		}
		gfile.P("Close() error")
		if method.Desc.IsStreamingClient() {
			gfile.P("Header() ", microMetadataPackage.Ident("Metadata"))
			gfile.P("Send(msg *", gfile.QualifiedGoIdent(method.Input.GoIdent), ") error")
		}
		if method.Desc.IsStreamingServer() {
			gfile.P("Recv() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error)")
		}
		gfile.P("}")
		gfile.P()
	}
}

func (g *Generator) generateServiceServerStreamInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	for _, method := range service.Methods {
		if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			continue
		}
		methodName := method.GoName
		gfile.P("type ", serviceName, "_", methodName, "Stream interface {")
		gfile.P("Context() ", contextPackage.Ident("Context"))
		gfile.P("SendMsg(msg interface{}) error")
		gfile.P("RecvMsg(msg interface{}) error")
		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("SendAndClose(msg *", gfile.QualifiedGoIdent(method.Output.GoIdent), ") error")
			gfile.P("CloseSend() error")
		}
		gfile.P("Close() error")
		if method.Desc.IsStreamingClient() {
			gfile.P("Recv() (*", gfile.QualifiedGoIdent(method.Input.GoIdent), ", error)")
		}
		if method.Desc.IsStreamingServer() {
			gfile.P("Send(msg *", gfile.QualifiedGoIdent(method.Output.GoIdent), ") error")
		}
		gfile.P("}")
		gfile.P()
	}
}

func generateEndpoints(method *protogen.Method) ([]*api_options.HttpRule, bool) {
	if method.Desc.Options() == nil {
		return nil, false
	}

	if !proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
		return nil, false
	}

	r := proto.GetExtension(method.Desc.Options(), api_options.E_Http)
	if r == nil {
		return nil, false
	}

	rule := r.(*api_options.HttpRule)
	rules := []*api_options.HttpRule{rule}
	rules = append(rules, rule.GetAdditionalBindings()...)

	return rules, method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient()
}

func getMicroApiMethod(method *protogen.Method) (*api_options.MicroMethod, bool) {
	if method.Desc.Options() == nil {
		return nil, false
	}

	if !proto.HasExtension(method.Desc.Options(), api_options.E_MicroMethod) {
		return nil, false
	}

	r := proto.GetExtension(method.Desc.Options(), api_options.E_MicroMethod)
	if r == nil {
		return nil, false
	}

	rule := r.(*api_options.MicroMethod)
	return rule, true
}

func getMicroApiService(service *protogen.Service) (*api_options.MicroService, bool) {
	if service.Desc.Options() == nil {
		return nil, false
	}

	if !proto.HasExtension(service.Desc.Options(), api_options.E_MicroService) {
		return nil, false
	}

	r := proto.GetExtension(service.Desc.Options(), api_options.E_MicroService)
	if r == nil {
		return nil, false
	}

	rule := r.(*api_options.MicroService)
	return rule, true
}

func getEndpoint(rule *api_options.HttpRule) (string, string, string) {
	var meth string
	var path string
	var body string

	switch {
	case len(rule.GetDelete()) > 0:
		meth = "DELETE"
		path = rule.GetDelete()
	case len(rule.GetGet()) > 0:
		meth = "GET"
		path = rule.GetGet()
	case len(rule.GetPatch()) > 0:
		meth = "PATCH"
		path = rule.GetPatch()
	case len(rule.GetPost()) > 0:
		meth = "POST"
		path = rule.GetPost()
	case len(rule.GetPut()) > 0:
		meth = "PUT"
		path = rule.GetPut()
	case rule.GetCustom() != nil:
		crule := rule.GetCustom()
		meth = crule.Kind
		path = crule.Path
	}

	body = rule.GetBody()
	return path, meth, body
}

func generateEndpoint(gfile *protogen.GeneratedFile, serviceName string, methodName string, rule *api_options.HttpRule, streaming bool) {
	path, meth, body := getEndpoint(rule)
	gfile.P("Name:", fmt.Sprintf(`"%s.%s",`, serviceName, methodName))
	gfile.P("Path:", fmt.Sprintf(`[]string{"%s"},`, path))
	//if vmethod, ok := httpMethodMap[meth]; ok {
	//	gfile.P("Method:", `[]string{`, httpPackage.Ident(vmethod), `},`)
	//} else {
	gfile.P("Method:", fmt.Sprintf(`[]string{"%s"},`, meth))
	//	}
	if len(rule.GetGet()) == 0 && body != "" {
		gfile.P("Body:", fmt.Sprintf(`"%s",`, body))
	}
	if streaming {
		gfile.P("Stream: true,")
	}
	gfile.P(`Handler: "rpc",`)
}

func (g *Generator) getGoIdentByXref(xref string) (protogen.GoIdent, error) {
	idx := strings.LastIndex(xref, ".")
	pkg := xref[:idx]
	msg := xref[idx+1:]
	for _, file := range g.plugin.Files {
		if strings.Compare(pkg, *(file.Proto.Package)) != 0 {
			continue
		}
		if ident, err := getGoIdentByMessage(file.Messages, msg); err == nil {
			return ident, nil
		}
	}
	return protogen.GoIdent{}, fmt.Errorf("not found")
}

func (g *Generator) getMessageByXref(xref string) (*protogen.Message, error) {
	idx := strings.LastIndex(xref, ".")
	pkg := xref[:idx]
	msg := xref[idx+1:]
	for _, file := range g.plugin.Files {
		if strings.Compare(pkg, *(file.Proto.Package)) != 0 {
			continue
		}
		if pmsg, err := getProtoMessage(file.Messages, msg); err == nil {
			return pmsg, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func getProtoMessage(messages []*protogen.Message, msg string) (*protogen.Message, error) {
	for _, message := range messages {
		if strings.Compare(msg, message.GoIdent.GoName) == 0 {
			return message, nil
		}
		if len(message.Messages) > 0 {
			if pmsg, err := getProtoMessage(message.Messages, msg); err == nil {
				return pmsg, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

func getGoIdentByMessage(messages []*protogen.Message, msg string) (protogen.GoIdent, error) {
	for _, message := range messages {
		if strings.Compare(msg, message.GoIdent.GoName) == 0 {
			return message.GoIdent, nil
		}
		if len(message.Messages) > 0 {
			if ident, err := getGoIdentByMessage(message.Messages, msg); err == nil {
				return ident, nil
			}
		}
	}
	return protogen.GoIdent{}, fmt.Errorf("not found")
}

func (g *Generator) generateServiceDesc(gfile *protogen.GeneratedFile, file *protogen.File, service *protogen.Service) {
	serviceName := service.GoName

	gfile.P("// ", serviceName, "_ServiceDesc", " is the ", grpcPackage.Ident("ServiceDesc"), " for ", serviceName, " service.")
	gfile.P("// It's only intended for direct use with ", grpcPackage.Ident("RegisterService"), ",")
	gfile.P("// and not to be introspected or modified (even as a copy)")
	gfile.P("var ", serviceName, "_ServiceDesc", " = ", grpcPackage.Ident("ServiceDesc"), " {")
	gfile.P("ServiceName: ", strconv.Quote(string(service.Desc.FullName())), ",")
	gfile.P("HandlerType: (*", serviceName, "Server)(nil),")
	gfile.P("Methods: []", grpcPackage.Ident("MethodDesc"), "{")
	for _, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue
		}
		gfile.P("{")
		gfile.P("MethodName: ", strconv.Quote(string(method.Desc.Name())), ",")
		gfile.P("Handler: ", method.GoName, ",")
		gfile.P("},")
	}
	gfile.P("},")
	gfile.P("Streams: []", grpcPackage.Ident("StreamDesc"), "{")
	for _, method := range service.Methods {
		if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			continue
		}
		gfile.P("{")
		gfile.P("StreamName: ", strconv.Quote(string(method.Desc.Name())), ",")
		gfile.P("Handler: ", method.GoName, ",")
		if method.Desc.IsStreamingServer() {
			gfile.P("ServerStreams: true,")
		}
		if method.Desc.IsStreamingClient() {
			gfile.P("ClientStreams: true,")
		}
		gfile.P("},")
	}
	gfile.P("},")
	gfile.P("Metadata: \"", file.Desc.Path(), "\",")
	gfile.P("}")
	gfile.P()
}

func (g *Generator) generateServiceName(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("var (")
	gfile.P(serviceName, "Name", "=", `"`, serviceName, `"`)
	gfile.P(")")
}

func (g *Generator) generateServiceEndpoints(gfile *protogen.GeneratedFile, service *protogen.Service, component string) {
	if component != "http" {
		return
	}
	serviceName := service.GoName

	gfile.P("var (")
	gfile.P(serviceName, "ServerEndpoints = []", microServerHttpPackage.Ident("EndpointMetadata"), "{")

	for _, method := range service.Methods {
		if proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
			if endpoints, streaming := generateEndpoints(method); endpoints != nil {
				for _, ep := range endpoints {
					epath, emethod, ebody := getEndpoint(ep)
					gfile.P("{")
					gfile.P(`Name: "`, serviceName+"."+method.GoName, `",`)
					gfile.P(`Path: "`, epath, `",`)
					gfile.P(`Method: "`, emethod, `",`)
					gfile.P(`Body: "`, ebody, `",`)
					gfile.P(`Stream: `, streaming, `,`)
					gfile.P("},")
				}
			}
		}
	}

	gfile.P("}")
	gfile.P(")")
}

func (g *Generator) writeErrors(plugin *protogen.Plugin) error {
	errorsMap := make(map[string]struct{})

	for _, file := range plugin.Files {
		for _, service := range file.Services {
			for _, method := range service.Methods {
				if method.Desc.Options() != nil {
					if proto.HasExtension(method.Desc.Options(), v2.E_Openapiv2Operation) {
						opts := proto.GetExtension(method.Desc.Options(), v2.E_Openapiv2Operation)
						if opts != nil {
							r := opts.(*v2.Operation)
							if r.Responses == nil {
								continue
							}

							for _, rsp := range r.Responses.ResponseCode {
								if schema := rsp.Value.GetJsonReference(); schema != nil {
									xref := schema.XRef
									if xref[0] == '.' {
										xref = xref[1:]
									}
									errorsMap[xref] = struct{}{}
								}
							}
						}
					}
					if proto.HasExtension(method.Desc.Options(), v3.E_Openapiv3Operation) {
						opts := proto.GetExtension(method.Desc.Options(), v3.E_Openapiv3Operation)
						if opts != nil {
							r := opts.(*v3.Operation)
							if r.Responses == nil {
								continue
							}
							resps := r.Responses.ResponseOrReference
							if r.Responses.GetDefault() != nil {
								resps = append(resps, &v3.NamedResponseOrReference{Name: "default", Value: r.Responses.GetDefault()})
							}
							for _, rsp := range resps {
								if schema := rsp.Value.GetReference(); schema != nil {
									xref := schema.XRef
									if xref[0] == '.' {
										xref = xref[1:]
									}
									errorsMap[xref] = struct{}{}
								}
							}
						}
					}
				}
			}
		}
	}

	var gfile *protogen.GeneratedFile
	var importPath protogen.GoImportPath

	if len(errorsMap) > 0 {

		var packageName string

		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}
			if len(file.Services) == 0 {
				continue
			}
			packageName = string(file.GoPackageName)
			importPath = file.GoImportPath
			break
		}

		if g.standalone {
			importPath = "."
		}

		gfile = plugin.NewGeneratedFile("micro_errors.pb.go", importPath)

		gfile.P("// Code generated by protoc-gen-go-micro. DO NOT EDIT.")
		gfile.P("// protoc-gen-go-micro version: " + versionComment)
		gfile.P()
		gfile.P("package ", packageName)
		gfile.P()

		gfile.Import(protojsonPackage)

		gfile.P("var (")
		gfile.P("marshaler = ", protojsonPackage.Ident("MarshalOptions"), "{}")
		gfile.P(")")
	}

	for xref := range errorsMap {
		msg, err := g.getMessageByXref(xref)
		if err != nil {
			return err
		}

		for _, field := range msg.Fields {
			if field.GoName == "Error" {
				return fmt.Errorf("failed generate Error() string interface for %s message %s already have Error field", field.Location.SourceFile, msg.Desc.Name())
			}
		}
		gfile.P(`func (m *`, msg.GoIdent.GoName, `) Error() string {`)
		gfile.P(`buf, _ := marshaler.Marshal(m)`)
		gfile.P("return string(buf)")
		gfile.P(`}`)
		// log.Printf("xref %#+v %v\n", msg.GoIdent.GoName, err)
	}

	return nil
}
