package main

import v3 "go.unistack.org/micro-proto/v3/openapiv3"

func getMediaType(eopt interface{}) string {
	ctype := "application/json"

	if eopt == nil {
		return ctype
	}

	if eopt == v3.E_Openapiv3Operation.InterfaceOf(v3.E_Openapiv3Operation.Zero()) {
		return ctype
	}

	opt, ok := eopt.(*v3.Operation)
	if !ok || opt.RequestBody == nil {
		return ctype
	}

	if opt.GetRequestBody() == nil {
		return ctype
	}

	if opt.GetRequestBody().GetRequestBody() == nil {
		return ctype
	}

	c := opt.GetRequestBody().GetRequestBody().GetContent()
	if c == nil {
		return ctype
	}

	for _, prop := range c.GetAdditionalProperties() {
		ctype = prop.Name
	}

	return ctype
}
