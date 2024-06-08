package echosec

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const openApiWithResponse = "openapi: 3.0.1\ninfo:\n  title: Test\n  version: v1\nservers:\n    - url: http://localhost:8080/api/v1\n\npaths:\n  \"/compute/groups\":\n    post:\n      operationId: computeGroups\n      parameters: \n        - name: foo\n          in: query\n          required: true\n          schema:\n            type: string\n      requestBody:\n        content:\n          application/json:\n            schema:\n              required:\n              - bar\n              properties:\n                bar:\n                  type: string\n      responses:\n        200:\n          description: groups created\n          content:\n            application/json:\n              schema:\n                required:\n                - a\n                - b\n                properties:\n                  a:\n                    type: string\n                  b:\n                    type: integer\n            \n      x-echosec:\n        function: do_stuff\n        params:\n          - great\n      "

type mockResponseWriter struct {
}

func (w *mockResponseWriter) Header() http.Header {
	return http.Header{}
}
func (w *mockResponseWriter) WriteHeader(_ int) {

}
func (w *mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func TestSecBlob(t *testing.T) {
	cfg, err := NewOApiConfig([]byte(openApiWithResponse), map[string]OApiValidationFunc{
		"do_stuff": func(c echo.Context, params []string) error {
			return nil
		},
	}, true)
	if err != nil {
		t.Error(err)
	}
	middlewareFunc := WithOpenApiConfig(cfg)
	f := middlewareFunc(func(c echo.Context) error {
		return nil
	})
	e := echo.New()
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/compute/groups?foo=bar", nil)
	ctx := e.NewContext(req, nil)
	res := echo.NewResponse(&mockResponseWriter{}, e)
	ctx.SetResponse(res)
	err = f(ctx)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, SecBlob(ctx, http.StatusOK, map[string]any{
		"a": "foobar",
		"b": 1,
	}))
	assert.NotNil(t, SecBlob(ctx, http.StatusOK, map[string]any{
		"a": 27,
		"b": "foobar",
	}))
}
