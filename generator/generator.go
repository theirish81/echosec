package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mitchellh/mapstructure"
	"github.com/theirish81/echosec"
	"os"
	"strings"
	"text/template"
)

var templ = `package {{.PackageName}}

// Generated code. DO NOT EDIT!
// These are the labels found in the OpenAPI document

{{ range $key, $value := .Labels }}
const {{$key}}Label = "{{$value}}"
{{ end }}
`

func main() {
	packageName := flag.String("package", "main", "The package name for the generated code")
	oApi := flag.String("spec", "", "The path to the openapi spec")
	flag.Parse()
	l := openapi3.NewLoader()
	l.IsExternalRefsAllowed = true
	oApiFile, err := l.LoadFromFile(*oApi)
	if err != nil {
		panic(fmt.Sprintf("Error loading openapi spec: %s", err))
	}
	labels := make(map[string]string, 0)
	for _, v := range oApiFile.Paths.Map() {
		for _, o := range v.Operations() {
			extension := o.Extensions["x-echosec"]
			if extension != nil {
				var localConfig echosec.OApiEchoSec
				err = mapstructure.Decode(extension, &localConfig)
				for _, l := range localConfig.Labels {
					labels[snakeToCamel(l.Label)] = l.Label
				}
			}
		}
	}

	t := template.Must(template.New("labels").Parse(templ))
	b := bytes.NewBuffer(make([]byte, 0))
	err = t.Execute(b, map[string]any{"PackageName": packageName, "Labels": labels})
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("echosec_constants.gen.go", b.Bytes(), 0666); err != nil {
		panic(err)
	}

}

func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}
