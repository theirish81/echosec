package echosec

import (
	"bufio"
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
	oApiBytes := make([]byte, 0)

	bytesReader := bufio.NewReader(bytes.NewBuffer(openapi))
	testBytes, err := bytesReader.Peek(2) //read 2 bytes
	if isGzipped(testBytes) {
		r, err := gzip.NewReader(bytesReader)
		if err != nil {
			return cfg, err
		}
		oApiBytes, err = io.ReadAll(r)
		if err != nil {
			return cfg, err
		}
	} else {
		oApiBytes = openapi
	}
	if err != nil {
		return cfg, err
	}
	doc, err := loader.LoadFromData(oApiBytes)
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

func isGzipped(b []byte) bool {
	if len(b) < 2 {
		return false
	}
	return b[0] == 0x1f && b[1] == 0x8b
}
