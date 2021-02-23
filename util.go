package main

import (
	"fmt"
	"strings"

	openapi_options "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	api_options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var (
	reflectPackage         = protogen.GoImportPath("reflect")
	stringsPackage         = protogen.GoImportPath("strings")
	fmtPackage             = protogen.GoImportPath("fmt")
	contextPackage         = protogen.GoImportPath("context")
	httpPackage            = protogen.GoImportPath("net/http")
	gorillaMuxPackage      = protogen.GoImportPath("github.com/gorilla/mux")
	chiPackage             = protogen.GoImportPath("github.com/go-chi/chi/v4")
	chiMiddlewarePackage   = protogen.GoImportPath("github.com/go-chi/chi/v4/middleware")
	microApiPackage        = protogen.GoImportPath("github.com/unistack-org/micro/v3/api")
	microClientPackage     = protogen.GoImportPath("github.com/unistack-org/micro/v3/client")
	microServerPackage     = protogen.GoImportPath("github.com/unistack-org/micro/v3/server")
	microClientHttpPackage = protogen.GoImportPath("github.com/unistack-org/micro-client-http/v3")
	deprecationComment     = "// Deprecated: Do not use."
)

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func generateServiceClient(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", unexport(serviceName), "Service struct {")
	gfile.P("c ", microClientPackage.Ident("Client"))
	gfile.P("name string")
	gfile.P("}")

	gfile.P("// New", serviceName, "Service create new service client")
	gfile.P("func New", serviceName, "Service(name string, c ", microClientPackage.Ident("Client"), ") ", serviceName, "Service {")
	gfile.P("return &", unexport(serviceName), "Service{c: c, name: name}")
	gfile.P("}")
	gfile.P()
}

func generateServiceClientMethods(gfile *protogen.GeneratedFile, service *protogen.Service, http bool) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	for _, method := range service.Methods {
		methodName := fmt.Sprintf("%s.%s", serviceName, method.GoName)
		generateClientFuncSignature(gfile, serviceName, method)

		if http && method.Desc.Options() != nil {
			if proto.HasExtension(method.Desc.Options(), openapi_options.E_Openapiv2Operation) {
				opts := proto.GetExtension(method.Desc.Options(), openapi_options.E_Openapiv2Operation)
				if opts != nil {
					r := opts.(*openapi_options.Operation)
					gfile.P("errmap := make(map[string]interface{}, ", len(r.Responses), ")")
					for code, response := range r.Responses {
						if response.Schema != nil && response.Schema.JsonSchema != nil {
							ref := response.Schema.JsonSchema.Ref
							if strings.HasPrefix(ref, "."+string(service.Desc.ParentFile().Package())+".") {
								ref = strings.TrimPrefix(ref, "."+string(service.Desc.ParentFile().Package())+".")
							}
							gfile.P(`errmap["`, code, `"] = &`, ref, "{}")
						}
					}
				}

				gfile.P("opts = append(opts,")
				gfile.P(microClientHttpPackage.Ident("ErrorMap"), "(errmap),")

				if proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
					endpoints, _ := generateEndpoints(method)
					path, method, body := getEndpoint(endpoints[0])
					gfile.P(microClientHttpPackage.Ident("Method"), `("`, method, `"),`)
					gfile.P(microClientHttpPackage.Ident("Path"), `("`, path, `"),`)
					gfile.P(microClientHttpPackage.Ident("Body"), `("`, body, `"),`)
				}

				gfile.P(")")
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
		gfile.P("return &", unexport(serviceName), "Service", method.GoName, "{stream}, nil")
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			gfile.P("type ", unexport(serviceName), "Service", method.GoName, " struct {")
			gfile.P("stream ", microClientPackage.Ident("Stream"))
			gfile.P("}")
		}

		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("func (s *", unexport(serviceName), "Service", method.GoName, ") RecvAndClose() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error) {")
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
		gfile.P("func (s *", unexport(serviceName), "Service", method.GoName, ") Close() error {")
		gfile.P("return s.stream.Close()")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Service", method.GoName, ") Context() ", contextPackage.Ident("Context"), " {")
		gfile.P("return s.stream.Context()")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Service", method.GoName, ") SendMsg(msg interface{}) error {")
		gfile.P("return s.stream.Send(msg)")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", unexport(serviceName), "Service", method.GoName, ") RecvMsg(msg interface{}) error {")
		gfile.P("return s.stream.Recv(msg)")
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingClient() {
			gfile.P("func (s *", unexport(serviceName), "Service", method.GoName, ") Send(msg *", gfile.QualifiedGoIdent(method.Input.GoIdent), ") error {")
			gfile.P("return s.stream.Send(msg)")
			gfile.P("}")
			gfile.P()
		}

		if method.Desc.IsStreamingServer() {
			gfile.P("func (s *", unexport(serviceName), "Service", method.GoName, ") Recv() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error) {")
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

func generateServiceServer(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", unexport(serviceName), "Handler struct {")
	gfile.P(serviceName, "Handler")
	gfile.P("}")
	gfile.P()
}

func generateServiceServerMethods(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	for _, method := range service.Methods {
		generateServerFuncSignature(gfile, serviceName, method, true)

		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			if !method.Desc.IsStreamingClient() {
				gfile.P("msg := &", gfile.QualifiedGoIdent(method.Input.GoIdent), "{}")
				gfile.P("if err := stream.Recv(msg); err != nil {")
				gfile.P("return err")
				gfile.P("}")
				gfile.P("return h.", serviceName, "Handler.", method.GoName, "(ctx, msg, &", unexport(serviceName), method.GoName, "Stream{stream})")
			} else {
				gfile.P("return h.", serviceName, "Handler.", method.GoName, "(ctx, &", unexport(serviceName), method.GoName, "Stream{stream})")
			}
		} else {
			gfile.P("return h.", serviceName, "Handler.", method.GoName, "(ctx, req, rsp)")
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

func generateServiceRegister(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("func Register", serviceName, "Handler(s ", microServerPackage.Ident("Server"), ", sh ", serviceName, "Handler, opts ...", microServerPackage.Ident("HandlerOption"), ") error {")
	gfile.P("type ", unexport(serviceName), " interface {")
	for _, method := range service.Methods {
		generateServerSignature(gfile, serviceName, method, true)
	}
	gfile.P("}")
	gfile.P("type ", serviceName, " struct {")
	gfile.P(unexport(serviceName))
	gfile.P("}")
	gfile.P("h := &", unexport(serviceName), "Handler{sh}")
	gfile.P("for _, endpoint := range New", serviceName, "Endpoints() {")
	gfile.P("opts = append(opts, ", microApiPackage.Ident("WithEndpoint"), "(endpoint))")
	gfile.P("}")
	gfile.P("return s.Handle(s.NewHandler(&", serviceName, "{h}, opts...))")
	gfile.P("}")
}

func generateServerFuncSignature(gfile *protogen.GeneratedFile, serviceName string, method *protogen.Method, private bool) {
	args := append([]interface{}{},
		"func (h *", unexport(serviceName), "Handler) ", method.GoName,
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

func generateClientFuncSignature(gfile *protogen.GeneratedFile, serviceName string, method *protogen.Method) {
	args := append([]interface{}{},
		"func (c *",
		unexport(serviceName),
		"Service) ",
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
		args = append(args, serviceName, "_", method.GoName, "Service")
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
		args = append(args, serviceName, "_", method.GoName, "Service")
	}
	args = append(args, ", error)")
	gfile.P(args...)
}

func generateServiceClientInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", serviceName, "Service interface {")
	for _, method := range service.Methods {
		generateClientSignature(gfile, serviceName, method)
	}
	gfile.P("}")
	gfile.P()
}

func generateServiceServerInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", serviceName, "Handler interface {")
	for _, method := range service.Methods {
		generateServerSignature(gfile, serviceName, method, false)
	}
	gfile.P("}")
	gfile.P()
}

func generateServiceClientStreamInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	for _, method := range service.Methods {
		if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			continue
		}
		methodName := method.GoName
		gfile.P("type ", serviceName, "_", methodName, "Service interface {")
		gfile.P("Context() ", contextPackage.Ident("Context"))
		gfile.P("SendMsg(msg interface{}) error")
		gfile.P("RecvMsg(msg interface{}) error")
		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("RecvAndClose() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error)")
		}
		gfile.P("Close() error")
		if method.Desc.IsStreamingClient() {
			gfile.P("Send(msg *", gfile.QualifiedGoIdent(method.Input.GoIdent), ") error")
		}
		if method.Desc.IsStreamingServer() {
			gfile.P("Recv() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error)")
		}
		gfile.P("}")
		gfile.P()
	}
}

func generateServiceServerStreamInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
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

func generateServiceEndpoints(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("// New", serviceName, "Endpoints provides api endpoints metdata for ", serviceName, " service")
	gfile.P("func New", serviceName, "Endpoints() []*", microApiPackage.Ident("Endpoint"), " {")
	gfile.P("return []*", microApiPackage.Ident("Endpoint"), "{")
	for _, method := range service.Methods {
		if method.Desc.Options() == nil {
			continue
		}
		if proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
			endpoints, streaming := generateEndpoints(method)
			for _, endpoint := range endpoints {
				gfile.P("&", microApiPackage.Ident("Endpoint"), "{")
				generateEndpoint(gfile, serviceName, method.GoName, endpoint, streaming)
				gfile.P("},")
			}
		}
	}
	gfile.P("}")
	gfile.P("}")
	gfile.P()
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
	gfile.P("Method:", fmt.Sprintf(`[]string{"%s"},`, meth))
	if len(rule.GetGet()) == 0 {
		gfile.P("Body:", fmt.Sprintf(`"%s",`, body))
	}
	if streaming {
		gfile.P("Stream: true,")
	}
	gfile.P(`Handler: "rpc",`)

	return
}
