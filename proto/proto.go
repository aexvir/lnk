package proto

import _ "embed"

//go:generate buf generate .
//go:generate go run postgen/main.go

//go:embed docs/openapi.yaml
var Schema []byte
