// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	jsonpbcodec "github.com/unistack-org/micro-codec-jsonpb/v3"
	annotations "go.unistack.org/micro-proto/v3/api"
	v3 "go.unistack.org/micro-proto/v3/openapiv3"
	"google.golang.org/protobuf/compiler/protogen"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type openapiv3Generator struct {
	plugin *protogen.Plugin

	requiredSchemas   []string // Names of schemas that need to be generated.
	generatedSchemas  []string // Names of schemas that have already been generated.
	linterRulePattern *regexp.Regexp
	namePattern       *regexp.Regexp
}

func (g *Generator) openapiv3Generate(component string, plugin *protogen.Plugin) error {
	og := &openapiv3Generator{
		plugin:            plugin,
		requiredSchemas:   make([]string, 0),
		generatedSchemas:  make([]string, 0),
		linterRulePattern: regexp.MustCompile(`\(-- .* --\)`),
		namePattern:       regexp.MustCompile("{(.*)=(.*)}"),
	}

	d := og.buildDocumentV3(plugin)

	bytes, err := jsonpbcodec.NewCodec(jsonpbcodec.MarshalOptions(jsonpb.MarshalOptions{Indent: "  "})).Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal: %s", err.Error())
	}

	outputFile := og.plugin.NewGeneratedFile(g.openapiFile, "")
	if _, err := outputFile.Write(bytes); err != nil {
		return fmt.Errorf("failed to write: %s", err.Error())
	}

	return nil
}

// buildDocumentV3 builds an OpenAPIv3 document for a plugin request.
func (g *openapiv3Generator) buildDocumentV3(plugin *protogen.Plugin) *v3.Document {
	d := &v3.Document{}
	d.Openapi = "3.0.3"
	d.Info = &v3.Info{
		Title:       "",
		Version:     "0.0.1",
		Description: "",
	}

	for _, file := range plugin.Files {
		if !proto.HasExtension(file.Desc.Options(), v3.E_Openapiv3Swagger) {
			continue
		}

		ext := proto.GetExtension(file.Desc.Options(), v3.E_Openapiv3Swagger)
		if ext == nil {
			continue
		}

		if doc, ok := ext.(*v3.Document); ok && doc != nil {
			if doc.Openapi != "" {
				d.Openapi = doc.Openapi
			}
			d.Info = proto.Clone(doc.Info).(*v3.Info)
			d.Servers = append(d.Servers, doc.Servers...)
		}

	}

	d.Paths = &v3.Paths{}
	d.Components = &v3.Components{
		Schemas: &v3.SchemasOrReferences{
			AdditionalProperties: []*v3.NamedSchemaOrReference{},
		},
	}
	for _, file := range g.plugin.Files {
		g.addPathsToDocumentV3(d, file)
	}
	for len(g.requiredSchemas) > 0 {
		count := len(g.requiredSchemas)
		for _, file := range g.plugin.Files {
			g.addSchemasToDocumentV3(d, file)
		}
		g.requiredSchemas = g.requiredSchemas[count:len(g.requiredSchemas)]
	}
	// Sort the paths.
	{
		pairs := d.Paths.Path
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Name < pairs[j].Name
		})
		d.Paths.Path = pairs
	}
	// Sort the schemas.
	{
		pairs := d.Components.Schemas.AdditionalProperties
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Name < pairs[j].Name
		})
		d.Components.Schemas.AdditionalProperties = pairs
	}
	return d
}

// filterCommentString removes line breaks and linter rules from comments.
func (g *openapiv3Generator) filterCommentString(c protogen.Comments) string {
	comment := string(c)
	comment = strings.Replace(comment, "\n", "", -1)
	comment = g.linterRulePattern.ReplaceAllString(comment, "")
	return strings.TrimSpace(comment)
}

// addPathsToDocumentV3 adds paths from a specified file descriptor.
func (g *openapiv3Generator) addPathsToDocumentV3(d *v3.Document, file *protogen.File) {
	for _, service := range file.Services {
		comment := g.filterCommentString(service.Comments.Leading)
		if d.Info.Title == "" {
			d.Info.Title = service.GoName
		}
		if d.Info.Description == "" {
			d.Info.Description = comment
		}
		for _, method := range service.Methods {
			comment := g.filterCommentString(method.Comments.Leading)
			inputMessage := method.Input
			outputMessage := method.Output
			operationID := service.GoName + "_" + method.GoName
			eopt := proto.GetExtension(method.Desc.Options(), v3.E_Openapiv3Operation)
			if eopt != nil && eopt != v3.E_Openapiv3Operation.InterfaceOf(v3.E_Openapiv3Operation.Zero()) {
				opt := eopt.(*v3.Operation)
				if opt.OperationId != "" {
					operationID = opt.OperationId
				}
			}
			xt := annotations.E_Http
			extension := proto.GetExtension(method.Desc.Options(), xt)
			var path string
			var methodName string
			var body string
			if extension != nil && extension != xt.InterfaceOf(xt.Zero()) {
				rule := extension.(*annotations.HttpRule)
				body = rule.Body
				switch pattern := rule.Pattern.(type) {
				case *annotations.HttpRule_Get:
					path = pattern.Get
					methodName = "GET"
				case *annotations.HttpRule_Post:
					path = pattern.Post
					methodName = "POST"
				case *annotations.HttpRule_Put:
					path = pattern.Put
					methodName = "PUT"
				case *annotations.HttpRule_Delete:
					path = pattern.Delete
					methodName = "DELETE"
				case *annotations.HttpRule_Patch:
					path = pattern.Patch
					methodName = "PATCH"
				case *annotations.HttpRule_Custom:
					path = "custom-unsupported"
				default:
					path = "unknown-unsupported"
				}
			}
			if methodName != "" {
				op, path2 := g.buildOperationV3(
					file, method, operationID, comment, path, body, inputMessage, outputMessage)
				g.addOperationV3(d, op, path2, methodName)
			}
		}
	}
}

// buildOperationV3 constructs an operation for a set of values.
func (g *openapiv3Generator) buildOperationV3(
	file *protogen.File,
	method *protogen.Method,
	operationID string,
	description string,
	path string,
	bodyField string,
	inputMessage *protogen.Message,
	outputMessage *protogen.Message,
) (*v3.Operation, string) {
	// coveredParameters tracks the parameters that have been used in the body or path.
	coveredParameters := make([]string, 0)
	if bodyField != "" {
		coveredParameters = append(coveredParameters, bodyField)
	}
	// Initialize the list of operation parameters.
	parameters := []*v3.ParameterOrReference{}
	// Build a list of header parameters.
	eopt := proto.GetExtension(method.Desc.Options(), v3.E_Openapiv3Operation)
	if eopt != nil && eopt != v3.E_Openapiv3Operation.InterfaceOf(v3.E_Openapiv3Operation.Zero()) {
		opt := eopt.(*v3.Operation)
		parameters = append(parameters, opt.Parameters...)
	}

	// Build a list of path parameters.
	pathParameters := make([]string, 0)
	if matches := g.namePattern.FindStringSubmatch(path); matches != nil {
		// Add the "name=" "name" value to the list of covered parameters.
		coveredParameters = append(coveredParameters, matches[1])
		// Convert the path from the starred form to use named path parameters.
		starredPath := matches[2]
		parts := strings.Split(starredPath, "/")
		// The starred path is assumed to be in the form "things/*/otherthings/*".
		// We want to convert it to "things/{thing}/otherthings/{otherthing}".
		for i := 0; i < len(parts)-1; i += 2 {
			section := parts[i]
			parameter := singular(section)
			parts[i+1] = "{" + parameter + "}"
			pathParameters = append(pathParameters, parameter)
		}
		// Rewrite the path to use the path parameters.
		newPath := strings.Join(parts, "/")
		path = strings.Replace(path, matches[0], newPath, 1)
	}
	// Add the path parameters to the operation parameters.
	for _, pathParameter := range pathParameters {
		parameters = append(parameters,
			&v3.ParameterOrReference{
				Oneof: &v3.ParameterOrReference_Parameter{
					Parameter: &v3.Parameter{
						Name:        pathParameter,
						In:          "path",
						Required:    true,
						Description: "The " + pathParameter + " id.",
						Schema: &v3.SchemaOrReference{
							Oneof: &v3.SchemaOrReference_Schema{
								Schema: &v3.Schema{
									Type: "string",
								},
							},
						},
					},
				},
			})
	}
	// Add any unhandled fields in the request message as query parameters.
	if bodyField != "*" {
		for _, field := range inputMessage.Fields {
			fieldName := string(field.Desc.Name())
			if !contains(coveredParameters, fieldName) {
				// Get the field description from the comments.
				fieldDescription := g.filterCommentString(field.Comments.Leading)
				parameters = append(parameters,
					&v3.ParameterOrReference{
						Oneof: &v3.ParameterOrReference_Parameter{
							Parameter: &v3.Parameter{
								Name:        fieldName,
								In:          "query",
								Description: fieldDescription,
								Required:    false,
								Schema: &v3.SchemaOrReference{
									Oneof: &v3.SchemaOrReference_Schema{
										Schema: &v3.Schema{
											Type: "string",
										},
									},
								},
							},
						},
					})
			}
		}
	}
	// Create the response.
	responses := &v3.Responses{
		ResponseOrReference: []*v3.NamedResponseOrReference{
			&v3.NamedResponseOrReference{
				Name: "200",
				Value: &v3.ResponseOrReference{
					Oneof: &v3.ResponseOrReference_Response{
						Response: &v3.Response{
							Description: "OK",
							Content:     g.responseContentForMessage(outputMessage),
						},
					},
				},
			},
		},
	}
	// Create the operation.
	op := &v3.Operation{
		Summary:     description,
		OperationId: operationID,
		Parameters:  parameters,
		Responses:   responses,
	}
	// If a body field is specified, we need to pass a message as the request body.
	if bodyField != "" {
		var bodyFieldScalarTypeName string
		var bodyFieldMessageTypeName string
		if bodyField == "*" {
			// Pass the entire request message as the request body.
			bodyFieldMessageTypeName = g.fullMessageTypeName(inputMessage)
		} else {
			// If body refers to a message field, use that type.
			for _, field := range inputMessage.Fields {
				if string(field.Desc.Name()) == bodyField {
					switch field.Desc.Kind() {
					case protoreflect.StringKind:
						bodyFieldScalarTypeName = "string"
					case protoreflect.MessageKind:
						bodyFieldMessageTypeName = g.fullMessageTypeName(field.Message)
					default:
						log.Printf("unsupported field type %+v", field.Desc)
					}
					break
				}
			}
		}
		var requestSchema *v3.SchemaOrReference
		if bodyFieldScalarTypeName != "" {
			requestSchema = &v3.SchemaOrReference{
				Oneof: &v3.SchemaOrReference_Schema{
					Schema: &v3.Schema{
						Type: bodyFieldScalarTypeName,
					},
				},
			}
		} else if bodyFieldMessageTypeName != "" {
			requestSchema = &v3.SchemaOrReference{
				Oneof: &v3.SchemaOrReference_Reference{
					Reference: &v3.Reference{
						XRef: g.schemaReferenceForTypeName(bodyFieldMessageTypeName),
					},
				},
			}
		}
		op.RequestBody = &v3.RequestBodyOrReference{
			Oneof: &v3.RequestBodyOrReference_RequestBody{
				RequestBody: &v3.RequestBody{
					Required: true,
					Content: &v3.MediaTypes{
						AdditionalProperties: []*v3.NamedMediaType{
							&v3.NamedMediaType{
								Name: "application/json",
								Value: &v3.MediaType{
									Schema: requestSchema,
								},
							},
						},
					},
				},
			},
		}
	}
	return op, path
}

// addOperationV3 adds an operation to the specified path/method.
func (g *openapiv3Generator) addOperationV3(d *v3.Document, op *v3.Operation, path string, methodName string) {
	var selectedPathItem *v3.NamedPathItem
	for _, namedPathItem := range d.Paths.Path {
		if namedPathItem.Name == path {
			selectedPathItem = namedPathItem
			break
		}
	}
	// If we get here, we need to create a path item.
	if selectedPathItem == nil {
		selectedPathItem = &v3.NamedPathItem{Name: path, Value: &v3.PathItem{}}
		d.Paths.Path = append(d.Paths.Path, selectedPathItem)
	}
	// Set the operation on the specified method.
	switch methodName {
	case http.MethodGet:
		selectedPathItem.Value.Get = op
	case http.MethodPost:
		selectedPathItem.Value.Post = op
	case http.MethodPut:
		selectedPathItem.Value.Put = op
	case http.MethodDelete:
		selectedPathItem.Value.Delete = op
	case http.MethodPatch:
		selectedPathItem.Value.Patch = op
	}
}

// schemaReferenceForTypeName returns an OpenAPI JSON Reference to the schema that represents a type.
func (g *openapiv3Generator) schemaReferenceForTypeName(typeName string) string {
	if !contains(g.requiredSchemas, typeName) {
		g.requiredSchemas = append(g.requiredSchemas, typeName)
	}
	parts := strings.Split(typeName, ".")
	lastPart := parts[len(parts)-1]
	return "#/components/schemas/" + lastPart
}

// itemsItemForTypeName is a helper constructor.
func (g *openapiv3Generator) itemsItemForTypeName(typeName string) *v3.ItemsItem {
	return &v3.ItemsItem{SchemaOrReference: []*v3.SchemaOrReference{&v3.SchemaOrReference{
		Oneof: &v3.SchemaOrReference_Schema{
			Schema: &v3.Schema{
				Type: typeName,
			},
		},
	}}}
}

// itemsItemForReference is a helper constructor.
func (g *openapiv3Generator) itemsItemForReference(xref string) *v3.ItemsItem {
	return &v3.ItemsItem{SchemaOrReference: []*v3.SchemaOrReference{&v3.SchemaOrReference{
		Oneof: &v3.SchemaOrReference_Reference{
			Reference: &v3.Reference{
				XRef: xref,
			},
		},
	}}}
}

// fullMessageTypeName builds the full type name of a message.
func (g *openapiv3Generator) fullMessageTypeName(message *protogen.Message) string {
	return "." + string(message.Desc.ParentFile().Package()) + "." + string(message.Desc.Name())
}

func (g *openapiv3Generator) responseContentForMessage(outputMessage *protogen.Message) *v3.MediaTypes {
	typeName := g.fullMessageTypeName(outputMessage)

	if typeName == ".google.protobuf.Empty" {
		return &v3.MediaTypes{}
	}

	if typeName == ".google.api.HttpBody" || typeName == ".micro.codec.Frame" {
		return &v3.MediaTypes{
			AdditionalProperties: []*v3.NamedMediaType{
				&v3.NamedMediaType{
					Name:  "application/octet-stream",
					Value: &v3.MediaType{},
				},
			},
		}
	}

	return &v3.MediaTypes{
		AdditionalProperties: []*v3.NamedMediaType{
			&v3.NamedMediaType{
				Name: "application/json",
				Value: &v3.MediaType{
					Schema: &v3.SchemaOrReference{
						Oneof: &v3.SchemaOrReference_Reference{
							Reference: &v3.Reference{
								XRef: g.schemaReferenceForTypeName(g.fullMessageTypeName(outputMessage)),
							},
						},
					},
				},
			},
		},
	}
}

// addSchemasToDocumentV3 adds info from one file descriptor.
func (g *openapiv3Generator) addSchemasToDocumentV3(d *v3.Document, file *protogen.File) {
	// For each message, generate a definition.
	for _, message := range file.Messages {
		typeName := g.fullMessageTypeName(message)
		// Only generate this if we need it and haven't already generated it.
		if !contains(g.requiredSchemas, typeName) ||
			contains(g.generatedSchemas, typeName) {
			continue
		}
		g.generatedSchemas = append(g.generatedSchemas, typeName)
		// Get the message description from the comments.
		messageDescription := g.filterCommentString(message.Comments.Leading)
		// Build an array holding the fields of the message.
		definitionProperties := &v3.Properties{
			AdditionalProperties: make([]*v3.NamedSchemaOrReference, 0),
		}
		for _, field := range message.Fields {
			// Check the field annotations to see if this is a readonly field.
			outputOnly := false
			extension := proto.GetExtension(field.Desc.Options(), annotations.E_FieldBehavior)
			if extension != nil {
				switch v := extension.(type) {
				case []annotations.FieldBehavior:
					for _, vv := range v {
						if vv == annotations.FieldBehavior_OUTPUT_ONLY {
							outputOnly = true
						}
					}
				default:
					log.Printf("unsupported extension type %T", extension)
				}
			}
			// Get the field description from the comments.
			fieldDescription := g.filterCommentString(field.Comments.Leading)
			// The field is either described by a reference or a schema.
			XRef := ""
			fieldSchema := &v3.Schema{
				Description: fieldDescription,
			}
			if outputOnly {
				fieldSchema.ReadOnly = true
			}
			if field.Desc.IsList() {
				fieldSchema.Type = "array"
				switch field.Desc.Kind() {
				case protoreflect.MessageKind:
					fieldSchema.Items = g.itemsItemForReference(
						g.schemaReferenceForTypeName(
							g.fullMessageTypeName(field.Message)))
				case protoreflect.StringKind:
					fieldSchema.Items = g.itemsItemForTypeName("string")
				case protoreflect.Int32Kind,
					protoreflect.Sint32Kind,
					protoreflect.Uint32Kind,
					protoreflect.Int64Kind,
					protoreflect.Sint64Kind,
					protoreflect.Uint64Kind,
					protoreflect.Sfixed32Kind,
					protoreflect.Fixed32Kind,
					protoreflect.Sfixed64Kind,
					protoreflect.Fixed64Kind:
					fieldSchema.Items = g.itemsItemForTypeName("integer")
				case protoreflect.EnumKind:
					fieldSchema.Items = g.itemsItemForTypeName("integer")
				case protoreflect.BoolKind:
					fieldSchema.Items = g.itemsItemForTypeName("boolean")
				case protoreflect.FloatKind, protoreflect.DoubleKind:
					fieldSchema.Items = g.itemsItemForTypeName("number")
				case protoreflect.BytesKind:
					fieldSchema.Items = g.itemsItemForTypeName("string")
				default:
					log.Printf("(TODO) Unsupported array type: %+v", g.fullMessageTypeName(field.Message))
				}
			} else if field.Desc.IsMap() &&
				field.Desc.MapKey().Kind() == protoreflect.StringKind &&
				field.Desc.MapValue().Kind() == protoreflect.StringKind {
				fieldSchema.Type = "object"
			} else {
				k := field.Desc.Kind()
				switch k {
				case protoreflect.MessageKind:
					typeName := g.fullMessageTypeName(field.Message)
					switch typeName {
					case ".google.protobuf.Timestamp":
						// Timestamps are serialized as strings
						fieldSchema.Type = "string"
						fieldSchema.Format = "RFC3339"
					case ".google.type.Date":
						// Dates are serialized as strings
						fieldSchema.Type = "string"
						fieldSchema.Format = "date"
					case ".google.type.DateTime":
						// DateTimes are serialized as strings
						fieldSchema.Type = "string"
						fieldSchema.Format = "date-time"
					default:
						// The field is described by a reference.
						XRef = g.schemaReferenceForTypeName(typeName)
					}
				case protoreflect.StringKind:
					fieldSchema.Type = "string"
				case protoreflect.Int32Kind,
					protoreflect.Sint32Kind,
					protoreflect.Uint32Kind,
					protoreflect.Int64Kind,
					protoreflect.Sint64Kind,
					protoreflect.Uint64Kind,
					protoreflect.Sfixed32Kind,
					protoreflect.Fixed32Kind,
					protoreflect.Sfixed64Kind,
					protoreflect.Fixed64Kind:
					fieldSchema.Type = "integer"
					fieldSchema.Format = k.String()
				case protoreflect.EnumKind:
					fieldSchema.Type = "integer"
					fieldSchema.Format = "enum"
				case protoreflect.BoolKind:
					fieldSchema.Type = "boolean"
				case protoreflect.FloatKind, protoreflect.DoubleKind:
					fieldSchema.Type = "number"
					fieldSchema.Format = k.String()
				case protoreflect.BytesKind:
					fieldSchema.Type = "string"
					fieldSchema.Format = "bytes"
				default:
					log.Printf("(TODO) Unsupported field type: %+v", g.fullMessageTypeName(field.Message))
				}
			}
			var value *v3.SchemaOrReference
			if XRef != "" {
				value = &v3.SchemaOrReference{
					Oneof: &v3.SchemaOrReference_Reference{
						Reference: &v3.Reference{
							XRef: XRef,
						},
					},
				}
			} else {
				value = &v3.SchemaOrReference{
					Oneof: &v3.SchemaOrReference_Schema{
						Schema: fieldSchema,
					},
				}
			}
			definitionProperties.AdditionalProperties = append(
				definitionProperties.AdditionalProperties,
				&v3.NamedSchemaOrReference{
					Name:  string(field.Desc.Name()),
					Value: value,
				},
			)
		}
		// Add the schema to the components.schema list.
		d.Components.Schemas.AdditionalProperties = append(d.Components.Schemas.AdditionalProperties,
			&v3.NamedSchemaOrReference{
				Name: string(message.Desc.Name()),
				Value: &v3.SchemaOrReference{
					Oneof: &v3.SchemaOrReference_Schema{
						Schema: &v3.Schema{
							Description: messageDescription,
							Properties:  definitionProperties,
						},
					},
				},
			},
		)
	}
}

// contains returns true if an array contains a specified string.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// singular produces the singular form of a collection name.
func singular(plural string) string {
	if strings.HasSuffix(plural, "ves") {
		return strings.TrimSuffix(plural, "ves") + "f"
	}
	if strings.HasSuffix(plural, "ies") {
		return strings.TrimSuffix(plural, "ies") + "y"
	}
	if strings.HasSuffix(plural, "s") {
		return strings.TrimSuffix(plural, "s")
	}
	return plural
}
