package echosec

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"io"
)

type OApiConfig struct {
	basePath   *string
	openapi    []byte
	router     routers.Router
	doc        *openapi3.T
	validators map[string]OApiValidationFunc
}

func NewOApiConfig(openapi []byte, validators map[string]OApiValidationFunc) (OApiConfig, error) {
	cfg := OApiConfig{validators: validators}
	loader := openapi3.Loader{Context: context.Background()}
	r, err := gzip.NewReader(bytes.NewReader(openapi))
	if err != nil {
		return cfg, err
	}
	f, err := io.ReadAll(r)
	if err != nil {
		return cfg, err
	}
	doc, err := loader.LoadFromData(f)
	if err != nil {
		return cfg, err
	}
	cfg.doc = doc
	cfg.router, err = gorillamux.NewRouter(cfg.doc)
	return cfg, err
}

type OApiEchoSec struct {
	Function string   `yaml:"function"`
	Params   []string `yaml:"params"`
}
