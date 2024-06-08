package echosec

import (
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

const openapi = "openapi: 3.0.1\ninfo:\n  title: Test\n  version: v1\nservers:\n    - url: http://localhost:8080/api/v1\n\npaths:\n  \"/compute/groups\":\n    post:\n      operationId: computeGroups\n      parameters: \n        - name: foo\n          in: query\n          required: true\n          schema:\n            type: string\n      requestBody:\n        content:\n          application/json:\n            schema:\n              required:\n              - bar\n              properties:\n                bar:\n                  type: string\n      responses:\n        200:\n          description: groups created\n      x-echosec:\n        function: do_stuff\n        params:\n          - great\n      "

const openapiB64 = "b3BlbmFwaTogMy4wLjEKaW5mbzoKICB0aXRsZTogVGVzdAogIHZlcnNpb246IHYxCnNlcnZlcnM6CiAgICAtIHVybDogaHR0cDovL2xvY2FsaG9zdDo4MDgwL2FwaS92MQoKcGF0aHM6CiAgIi9jb21wdXRlL2dyb3VwcyI6CiAgICBwb3N0OgogICAgICBvcGVyYXRpb25JZDogY29tcHV0ZUdyb3VwcwogICAgICBwYXJhbWV0ZXJzOiAKICAgICAgICAtIG5hbWU6IGZvbwogICAgICAgICAgaW46IHF1ZXJ5CiAgICAgICAgICByZXF1aXJlZDogdHJ1ZQogICAgICAgICAgc2NoZW1hOgogICAgICAgICAgICB0eXBlOiBzdHJpbmcKICAgICAgcmVxdWVzdEJvZHk6CiAgICAgICAgY29udGVudDoKICAgICAgICAgIGFwcGxpY2F0aW9uL2pzb246CiAgICAgICAgICAgIHNjaGVtYToKICAgICAgICAgICAgICByZXF1aXJlZDoKICAgICAgICAgICAgICAtIGJhcgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICBiYXI6CiAgICAgICAgICAgICAgICAgIHR5cGU6IHN0cmluZwogICAgICByZXNwb25zZXM6CiAgICAgICAgMjAwOgogICAgICAgICAgZGVzY3JpcHRpb246IGdyb3VwcyBjcmVhdGVkCiAgICAgIHgtZWNob3NlYzoKICAgICAgICBmdW5jdGlvbjogZG9fc3R1ZmYKICAgICAgICBwYXJhbXM6CiAgICAgICAgICAtIGdyZWF0CiAgICAgIA=="

const openapiGzipped = "H4sIAAAAAAAAA3VRMW7DMAzc/QoiuyOnXQKNXYru3QtVpmMVjsRQVFD/vpJjJErabtLx7ngkA6E35DQ8b7vtrnF+CLoBECcTanjHKPl3Ro4ueA3nXRORy7eQAFpIPGkYRUgrNQVrpjFE0ftu36lsq7KgISPjwt8oG46UBNWBQ6K4uZhQUSwvgEDIRnKvt17Dyn5dyCuBDJsjSkkAK1Ri+AxqGEK4YgAuBz4l5LnCGE/JMWZz4YRVIdoRj0ZXSF7CTNk0Cjt/aG76vJOX0M83rg1e0EstNkSTs8sk6ivm1d0Z/9WsyvaAt/Bp+AEjLpsSh/GRDYX9G/xnmkjBx9rkqetqcY/RsiNZrn+5GlhGI9ivrO8Wbb452ptsSN5eFH34iJKG4VpazneXuc222W9FfgBdLMjrjwIAAA=="

func TestWithOpenApiConfig(t *testing.T) {

	cfg, err := NewOApiConfig([]byte(openapi), map[string]OApiValidationFunc{
		"do_stuff": func(c echo.Context, params []string) error {
			c.Set("caught", true)
			c.Set("param", params[0])
			return nil
		},
	}, false)
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

func TestWithOpenApiConfigB64(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString(openapiB64)
	cfg, err := NewOApiConfig(data, map[string]OApiValidationFunc{
		"do_stuff": func(c echo.Context, params []string) error {
			c.Set("caught", true)
			c.Set("param", params[0])
			return nil
		},
	}, false)
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
	}, false)
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

func TestWithOpenApiConfigValidationEnabled(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString(openapiGzipped)
	cfg, err := NewOApiConfig(data, map[string]OApiValidationFunc{
		"do_stuff": func(c echo.Context, params []string) error {
			c.Set("caught", true)
			c.Set("param", params[0])
			return nil
		},
	}, true)
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

	// Success
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/compute/groups?foo=one", strings.NewReader(`{"bar":"two"}"`))
	req.Header.Set("content-type", "application/json")
	ctx := e.NewContext(req, nil)
	err = f(ctx)
	assert.Nil(t, err)

	// Fail
	req, _ = http.NewRequest("POST", "http://localhost:8080/api/v1/compute/groups?foo=one", strings.NewReader(`{"foobar":"two"}"`))
	req.Header.Set("content-type", "application/json")
	ctx = e.NewContext(req, nil)
	err = f(ctx)
	assert.NotNil(t, err)
}
