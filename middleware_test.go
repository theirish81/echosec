package echosec

import (
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const openapi = "b3BlbmFwaTogMy4wLjEKaW5mbzoKICB0aXRsZTogVGVzdAogIHZlcnNpb246IHYxCnNlcnZlcnM6CiAgLSB1cmw6IGh0dHA6Ly9sb2NhbGhvc3Q6ODA4MC9hcGkvdjEKcGF0aHM6CiAgIi9jb21wdXRlL2dyb3VwcyI6CiAgICBwb3N0OgogICAgICBvcGVyYXRpb25JZDogY29tcHV0ZUdyb3VwcwogICAgICByZXNwb25zZXM6CiAgICAgICAgMjAwOgogICAgICAgICAgZGVzY3JpcHRpb246IGdyb3VwcyBjcmVhdGVkCiAgICAgIHgtZWNob3NlYzoKICAgICAgICBmdW5jdGlvbjogZG9fc3R1ZmYKICAgICAgICBwYXJhbXM6CiAgICAgICAgICAtIGdyZWF0"

const openapiGzipped = "H4sIAAAAAAAAA02OvW7DMAyEdz8Fkd2Rki6BXqDo3r0QZLoW4IgESQV9/Mg/ibvxyO/uSIwlcg7wcfbnS5fLSKEDsGwzBvhGtaYeKJqpBHhcOkVZ5AL1UGUOMJlxcG6mFOeJ1MLN37xroa7hHG1a4ZNLdOdq6H6FKutpWQLwYlgnAGKUaK3oawiw058rvAOCylQU9eUAuHp/CIABNUlmW7/diiAJRsNhp/56TO1LTIdtrCVtjoF+1Oo4vk8cJd71f0PfYlveExuvwlI4AQAA"

func TestWithOpenApiConfig(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString(openapi)
	cfg, err := NewOApiConfig(data, map[string]OApiValidationFunc{
		"do_stuff": func(c echo.Context, params []string) error {
			c.Set("caught", true)
			c.Set("param", params[0])
			return nil
		},
	})
	if err != nil {
		t.Error(err)
	}
	middlewareFunc := WithOpenApiConfig(cfg)
	f := middlewareFunc(func(c echo.Context) error {
		assert.Equal(t, true, c.Get("caught"))
		assert.Equal(t, "great", c.Get("param"))
		return nil
	})
	e := echo.New()
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/compute/groups", nil)
	ctx := e.NewContext(req, nil)
	err = f(ctx)
	assert.Nil(t, err)
}

func TestWithOpenApiConfigGzip(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString(openapiGzipped)
	cfg, err := NewOApiConfig(data, map[string]OApiValidationFunc{
		"do_stuff": func(c echo.Context, params []string) error {
			c.Set("caught", true)
			c.Set("param", params[0])
			return nil
		},
	})
	if err != nil {
		t.Error(err)
	}
	middlewareFunc := WithOpenApiConfig(cfg)
	f := middlewareFunc(func(c echo.Context) error {
		assert.Equal(t, true, c.Get("caught"))
		assert.Equal(t, "great", c.Get("param"))
		return nil
	})
	e := echo.New()
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/compute/groups", nil)
	ctx := e.NewContext(req, nil)
	err = f(ctx)
	assert.Nil(t, err)
}
