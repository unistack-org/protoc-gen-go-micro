package main

import (
	"fmt"
	"strings"

	openapi_options "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	api_options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func lowerFirst(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func generateServiceClient(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", lowerFirst(serviceName), "Service struct {")
	gfile.P("c micro_client.Client")
	gfile.P("name string")
	gfile.P("}")

	gfile.P("// New", serviceName, "Service create new service client")
	gfile.P("func New", serviceName, "Service(name string, c micro_client.Client) ", serviceName, "Service {")
	gfile.P("return &", lowerFirst(serviceName), "Service{c: c, name: name}")
	gfile.P("}")
	gfile.P()
}

func generateServiceClientMethods(gfile *protogen.GeneratedFile, service *protogen.Service, http bool) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	for _, method := range service.Methods {
		methodName := fmt.Sprintf("%s.%s", serviceName, method.GoName)
		gfile.P("func (c *", lowerFirst(serviceName), "Service) ", generateClientSignature(serviceName, method), "{")

		if http && method.Desc.Options() != nil {
			if proto.HasExtension(method.Desc.Options(), openapi_options.E_Openapiv2Operation) {
				opts := proto.GetExtension(method.Desc.Options(), openapi_options.E_Openapiv2Operation)
				if opts != nil {
					r := opts.(*openapi_options.Operation)
					gfile.P("errmap := make(map[string]interface{}, ", len(r.Responses), ")")
					for code, response := range r.Responses {
						if response.Schema == nil || response.Schema.JsonSchema == nil {
							continue
						}
						ref := response.Schema.JsonSchema.Ref
						if strings.HasPrefix(ref, "."+string(service.Desc.ParentFile().Package())+".") {
							ref = strings.TrimPrefix(ref, "."+string(service.Desc.ParentFile().Package())+".")
						}
						gfile.P(`errmap["`, code, `"] = &`, ref, "{}")
					}
				}

				gfile.P("opts = append(opts,")
				gfile.P("micro_client_http.ErrorMap(errmap),")

				if proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
					endpoints, _ := generateEndpoints(method)
					path, method, body := getEndpoint(endpoints[0])
					gfile.P(`micro_client_http.Method("`, method, `"),`)
					gfile.P(`micro_client_http.Path("`, path, `"),`)
					gfile.P(`micro_client_http.Body("`, body, `"),`)
				}

				gfile.P(")")
			}
		}

		if !method.Desc.IsStreamingServer() && !method.Desc.IsStreamingClient() {
			gfile.P("rsp := &", method.Output.Desc.Name(), "{}")
			gfile.P(`err := c.c.Call(ctx, c.c.NewRequest(c.name, "`, methodName, `", req), rsp, opts...)`)
			gfile.P("if err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
			gfile.P("return rsp, nil")
			gfile.P("}")
			gfile.P()
			return
		}

		gfile.P(`stream, err := c.c.Stream(ctx, c.c.NewRequest(c.name, "`, methodName, `", &`, method.Input.Desc.Name(), `{}), opts...)`)
		gfile.P("if err != nil {")
		gfile.P("return nil, err")
		gfile.P("}")

		if !method.Desc.IsStreamingClient() {
			gfile.P("if err := stream.Send(req); err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
		}
		gfile.P("return &", lowerFirst(serviceName), "Service", method.GoName, "{stream}, nil")
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			gfile.P("type ", lowerFirst(serviceName), "Service", method.GoName, " struct {")
			gfile.P("stream micro_client.Stream")
			gfile.P("}")
		}

		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("func (s *", lowerFirst(serviceName), "Service", method.GoName, ") RecvAndClose() (*", method.Output.Desc.Name(), ", error) {")
			gfile.P("m := &", method.Output.Desc.Name(), "{}")
			gfile.P("err := s.RecvMsg(m)")
			gfile.P("if err == nil {")
			gfile.P("err = s.Close()")
			gfile.P("}")
			gfile.P("if err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
			gfile.P("return m, nil")
			gfile.P("}")
		}

		gfile.P()
		gfile.P("func (s *", lowerFirst(serviceName), "Service", method.GoName, ") Close() error {")
		gfile.P("return s.stream.Close()")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", lowerFirst(serviceName), "Service", method.GoName, ") Context() context.Context {")
		gfile.P("return s.stream.Context()")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", lowerFirst(serviceName), "Service", method.GoName, ") SendMsg(m interface{}) error {")
		gfile.P("return s.stream.Send(m)")
		gfile.P("}")
		gfile.P()
		gfile.P("func (s *", lowerFirst(serviceName), "Service", method.GoName, ") RecvMsg(m interface{}) error {")
		gfile.P("return s.stream.Recv(m)")
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingClient() {
			gfile.P("func (s *", lowerFirst(serviceName), "Service", method.GoName, ") Send(m *", method.Input.Desc.Name(), ") error {")
			gfile.P("return s.stream.Send(m)")
			gfile.P("}")
			gfile.P()
		}

		if method.Desc.IsStreamingServer() {
			gfile.P("func (s *", lowerFirst(serviceName), "Service", method.GoName, ") Recv() (*", method.Output.Desc.Name(), ", error) {")
			gfile.P("m := &", method.Output.Desc.Name(), "{}")
			gfile.P("if err := s.stream.Recv(m); err != nil {")
			gfile.P("return nil, err")
			gfile.P("}")
			gfile.P("return m, nil")
			gfile.P("}")
			gfile.P()
		}
	}
}

func generateServiceServer(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", lowerFirst(serviceName), "Handler struct {")
	gfile.P(serviceName, "Handler")
	gfile.P("}")
	gfile.P()
}

func generateServiceServerMethods(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	for _, method := range service.Methods {
		//methodName := fmt.Sprintf("%s.%s", serviceName, method.GoName)
		gfile.P("func (h *", lowerFirst(serviceName), "Handler) ", generateServerSignature(serviceName, method), "{")
		generateServerSignature(serviceName, method)

		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			if !method.Desc.IsStreamingClient() {
				gfile.P("m := &", method.Input.Desc.Name(), "{}")
				gfile.P("if err := stream.Recv(m); err != nil {")
				gfile.P("return err")
				gfile.P("}")
				gfile.P("return h.", serviceName, "Handler.", method.GoName, "(ctx, m, &", lowerFirst(serviceName), method.GoName, ",Stream{stream})")
			} else {
				gfile.P("return h.", serviceName, "Handler.", method.GoName, "(ctx, &", lowerFirst(serviceName), method.GoName, "Stream{stream})")
			}
		} else {
			gfile.P("return h.", serviceName, "Handler.", method.GoName, "(ctx, req, rsp)")
		}
		gfile.P("}")
		gfile.P()

		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			gfile.P("type ", lowerFirst(serviceName), method.GoName, "Stream struct {")
			gfile.P("stream micro_server.Stream")
			gfile.P("}")
		}

		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("func (s *", lowerFirst(serviceName), method.GoName, "Stream) SendAndClose(m *", method.Output.Desc.Name(), ") error {")
			gfile.P("err := s.SendMsg(m)")
			gfile.P("if err == nil {")
			gfile.P("err = s.stream.Close()")
			gfile.P("}")
			gfile.P("return err")
			gfile.P("}")
		}

		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			gfile.P("func (s *", lowerFirst(serviceName), method.GoName, "Stream) Close() error {")
			gfile.P("return s.stream.Close()")
			gfile.P("}")
			gfile.P()

			gfile.P("func (s *", lowerFirst(serviceName), method.GoName, "Stream) Context() context.Context {")
			gfile.P("return s.stream.Context()")
			gfile.P("}")
			gfile.P()

			gfile.P("func (s *", lowerFirst(serviceName), method.GoName, "Stream) SendMsg(m interface{}) error {")
			gfile.P("return s.stream.Send(m)")
			gfile.P("}")
			gfile.P()

			gfile.P("func (s *", lowerFirst(serviceName), method.GoName, "Stream) RecvMsg(m interface{}) error {")
			gfile.P("return s.stream.Recv(m)")
			gfile.P("}")
			gfile.P()

			if method.Desc.IsStreamingServer() {
				gfile.P("func (s *", lowerFirst(serviceName), method.GoName, "Stream) Send(m *", method.Output.Desc.Name(), ") error {")
				gfile.P("return s.stream.Send(m)")
				gfile.P("}")
				gfile.P()
			}

			if method.Desc.IsStreamingClient() {
				gfile.P("func (s *", lowerFirst(serviceName), method.GoName, "Stream) Recv() (*", method.Input.Desc.Name(), ", error) {")
				gfile.P("m := &", method.Input.Desc.Name(), "{}")
				gfile.P("if err := s.stream.Recv(m); err != nil {")
				gfile.P("return nil, err")
				gfile.P("}")
				gfile.P("return m, nil")
				gfile.P("}")
				gfile.P()
			}
		}

	}
}

func generateServiceRegister(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("func Register", serviceName, "Handler(s micro_server.Server, sh ", serviceName, "Handler, opts ...micro_server.HandlerOption) error {")
	gfile.P("type ", lowerFirst(serviceName), " interface {")
	for _, method := range service.Methods {
		gfile.P(generateServerSignature(serviceName, method))
	}
	gfile.P("}")
	gfile.P("type ", serviceName, " struct {")
	gfile.P(lowerFirst(serviceName))
	gfile.P("}")
	gfile.P("h := &", lowerFirst(serviceName), "Handler{sh}")
	gfile.P("for _, endpoint := range New", serviceName, "Endpoints() {")
	gfile.P("opts = append(opts, micro_api.WithEndpoint(endpoint))")
	gfile.P("}")
	gfile.P("return s.Handle(s.NewHandler(&", serviceName, "{h}, opts...))")
	gfile.P("}")
}

func generateServerSignature(serviceName string, method *protogen.Method) string {
	methodName := string(method.GoName)
	req := []string{"ctx context.Context"}
	ret := "error"

	if !method.Desc.IsStreamingClient() {
		req = append(req, "req *"+string(method.Input.Desc.Name()))
	}
	if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
		req = append(req, "stream "+serviceName+"_"+methodName+"Stream")
	}
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		req = append(req, "rsp *"+string(method.Output.Desc.Name()))
	}
	return methodName + "(" + strings.Join(req, ", ") + ") " + ret
}

func generateClientSignature(serviceName string, method *protogen.Method) string {
	methodName := string(method.GoName)
	req := ", req *" + string(method.Input.Desc.Name())
	if method.Desc.IsStreamingClient() {
		req = ""
	}
	rsp := "*" + string(method.Output.Desc.Name())
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		rsp = serviceName + "_" + methodName + "Service"
	}
	return fmt.Sprintf("%s(ctx context.Context%s, opts ...micro_client.CallOption) (%s, error)", methodName, req, rsp)
}

func generateServiceClientInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", serviceName, "Service interface {")
	for _, method := range service.Methods {
		gfile.P(generateClientSignature(serviceName, method))
	}
	gfile.P("}")
	gfile.P()
}

func generateServiceServerInterface(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("type ", serviceName, "Handler interface {")
	for _, method := range service.Methods {
		gfile.P(generateServerSignature(serviceName, method))
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
		gfile.P("Context() context.Context")
		gfile.P("SendMsg(msg interface{}) error")
		gfile.P("RecvMsg(msg interface{}) error")
		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("RecvAndClose() (", method.Output.Desc.Name(), ", error)")
		}
		gfile.P("Close() error")
		if method.Desc.IsStreamingClient() {
			gfile.P("Send(msg *", method.Input.Desc.Name(), ") error")
		}
		if method.Desc.IsStreamingServer() {
			gfile.P("Recv() (msg *", method.Output.Desc.Name(), ", error)")
		}
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
		gfile.P("Context() context.Context")
		gfile.P("SendMsg(msg interface{}) error")
		gfile.P("RecvMsg(msg interface{}) error")
		if method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
			gfile.P("SendAndClose(msg *", method.Output.Desc.Name(), ") error")
		}
		gfile.P("Close() error")
		if method.Desc.IsStreamingClient() {
			gfile.P("Send(msg *", method.Output.Desc.Name(), ") error")
		}
		if method.Desc.IsStreamingServer() {
			gfile.P("Recv() (msg *", method.Input.Desc.Name(), ", error)")
		}
	}
}

func generateServiceEndpoints(gfile *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := strings.TrimSuffix(service.GoName, "Service")
	gfile.P("// New", serviceName, "Endpoints provides api endpoints metdata for ", serviceName, " service")
	gfile.P("func New", serviceName, "Endpoints() []*micro_api.Endpoint {")
	gfile.P("return []*", "micro_api.Endpoint{")
	for _, method := range service.Methods {
		if method.Desc.Options() == nil {
			continue
		}
		if proto.HasExtension(method.Desc.Options(), api_options.E_Http) {
			endpoints, streaming := generateEndpoints(method)
			for _, endpoint := range endpoints {
				gfile.P("&", "micro_api.Endpoint{")
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
