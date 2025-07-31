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
	"unicode"
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
					labels[SnakeToCamel(l.Label)] = l.Label
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

func SnakeToCamel(s string) string {
	if !strings.Contains(s, "_") {
		// If it's not snake_case, just capitalize the first rune
		runes := []rune(s)
		if len(runes) == 0 {
			return s
		}
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	}

	parts := strings.Split(s, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}
