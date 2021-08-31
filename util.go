package main

import (
	"fmt"
	"strings"

	api_options "github.com/unistack-org/micro-proto/api"
	openapiv2_options "github.com/unistack-org/micro-proto/openapiv2"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var httpMethodMap = map[string]string{
	"GET":     "MethodGet",
	"HEAD":    "MethodHead",
	"POST":    "MethodPost",
	"PUT":     "MethodPut",
	"PATCH":   "MethodPatch",
	"DELETE":  "MethodDelete",
	"CONNECT": "MethodConnect",
	"OPTIONS": "MethodOptions",
	"TRACE":   "MethodTrace",
}

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func generateServiceClient(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	// if rule, ok := getMicroApiService(service); ok {
	//		gfile.P("// client wrappers ", strings.Join(rule.ClientWrappers, ", "))
	//	}
	gfile.P("type ", unexport(serviceName), "Client struct {")
	gfile.P("c ", microClientPackage.Ident("Client"))
	gfile.P("name string")
	gfile.P("}")

	gfile.P("func New", serviceName, "Client(name string, c ", microClientPackage.Ident("Client"), ") ", serviceName, "Client {")
	gfile.P("return &", unexport(serviceName), "Client{c: c, name: name}")
	gfile.P("}")
	gfile.P()
}

func generateServiceClientMethods(gfile *protogen.GeneratedFile, service *protogen.Service, http bool) {
	serviceName := service.GoName
	for _, method := range service.Methods {
		methodName := fmt.Sprintf("%s.%s", serviceName, method.GoName)
		generateClientFuncSignature(gfile, serviceName, method)

		if http && method.Desc.Options() != nil {
			if proto.HasExtension(method.Desc.Options(), openapiv2_options.E_Openapiv2Operation) {
				opts := proto.GetExtension(method.Desc.Options(), openapiv2_options.E_Openapiv2Operation)
				if opts != nil {
					r := opts.(*openapiv2_options.Operation)
					gfile.P("errmap := make(map[string]interface{}, ", len(r.Responses.ResponseCode), ")")
					for _, rsp := range r.Responses.ResponseCode {
						if schema := rsp.Value.GetJsonReference(); schema != nil {
							ref := schema.XRef
							if strings.HasPrefix(ref, "."+string(service.Desc.ParentFile().Package())+".") {
								ref = strings.TrimPrefix(ref, "."+string(service.Desc.ParentFile().Package())+".")
							}
							if ref == "micro.codec.Frame" {
								gfile.P(`errmap["`, rsp.Name, `"] = &`, microCodecPackage.Ident("Frame"), "{}")
							} else {
								gfile.P(`errmap["`, rsp.Name, `"] = &`, ref, "{}")
							}
						}
					}
				}

				gfile.P("opts = append(opts,")
				gfile.P(microClientHttpPackage.Ident("ErrorMap"), "(errmap),")
				gfile.P(")")
			}
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
		}
		if rule, ok := getMicroApiMethod(method); ok {
			if rule.Timeout > 0 {
				gfile.P("opts = append(opts, ", microClientPackage.Ident("WithRequestTimeout"), "(", timePackage.Ident("Second"), "*", rule.Timeout, "))")
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
			gfile.P("func (s *", unexport(serviceName), "Client", method.GoName, ") RecvAndClose() (*", gfile.QualifiedGoIdent(method.Output.GoIdent), ", error) {")
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

func generateServiceServer(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("type ", unexport(serviceName), "Server struct {")
	gfile.P(serviceName, "Server")
	gfile.P("}")
	gfile.P()
}

func generateServiceServerMethods(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	for _, method := range service.Methods {
		generateServerFuncSignature(gfile, serviceName, method, true)
		if rule, ok := getMicroApiMethod(method); ok {
			if rule.Timeout > 0 {
				gfile.P("var cancel ", contextPackage.Ident("CancelFunc"))
				gfile.P("ctx, cancel = ", contextPackage.Ident("WithTimeout"), "(ctx, ", timePackage.Ident("Second"), "*", rule.Timeout, ")")
				gfile.P("defer cancel()")
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

func generateServiceRegister(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("func Register", serviceName, "Server(s ", microServerPackage.Ident("Server"), ", sh ", serviceName, "Server, opts ...", microServerPackage.Ident("HandlerOption"), ") error {")
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
	gfile.P("for _, endpoint := range ", serviceName, "Endpoints {")
	gfile.P("nopts = append(nopts, ", microApiPackage.Ident("WithEndpoint"), "(&endpoint))")
	gfile.P("}")
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

func generateClientFuncSignature(gfile *protogen.GeneratedFile, serviceName string, method *protogen.Method) {
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

func generateServiceClientInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("type ", serviceName, "Client interface {")
	for _, method := range service.Methods {
		generateClientSignature(gfile, serviceName, method)
	}
	gfile.P("}")
	gfile.P()
}

func generateServiceServerInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := service.GoName
	gfile.P("type ", serviceName, "Server interface {")
	for _, method := range service.Methods {
		generateServerSignature(gfile, serviceName, method, false)
	}
	gfile.P("}")
	gfile.P()
}

func generateServiceClientStreamInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
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
	serviceName := service.GoName
	gfile.P("var (")
	gfile.P(serviceName, "Name", "=", `"`, serviceName, `"`)
	gfile.P()
	gfile.P(serviceName, "Endpoints", "=", "[]", microApiPackage.Ident("Endpoint"), "{")
	for _, method := range service.Methods {
		if method.Desc.Options() == nil {
			continue
		}
		if proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
			endpoints, streaming := generateEndpoints(method)
			for _, endpoint := range endpoints {
				gfile.P("{")
				generateEndpoint(gfile, serviceName, method.GoName, endpoint, streaming)
				gfile.P("},")
			}
		}
	}
	gfile.P("}")
	gfile.P(")")
	gfile.P()

	gfile.P("func New", serviceName, "Endpoints()", "[]", microApiPackage.Ident("Endpoint"), "{")
	gfile.P("return ", serviceName, "Endpoints")
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

	return
}
