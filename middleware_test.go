package echosec

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	m := WithManualConfig(ManualConfig{
		PathMapping: PathItems{
			{
				Patterns: Patterns{"/foo", "/bar"},
				Methods: ValidationMap{
					"GET": func(c echo.Context) error {
						return nil
					},
				},
				PathValidation: func(c echo.Context) error {
					return errors.New("no access")
				},
			},
		},
		DefaultValidation: func(c echo.Context) error {
			return errors.New("baseline")
		},
	})
	assert.Nil(t, m(func(c echo.Context) error {
		return nil
	})(newContext("GET", "/foo")))

	assert.Nil(t, m(func(c echo.Context) error {
		return nil
	})(newContext("GET", "/bar")))

	err := m(func(c echo.Context) error {
		return nil
	})(newContext("PUT", "/bar"))
	assert.NotNil(t, err)
	assert.Equal(t, "no access", err.Error())

	err = m(func(c echo.Context) error {
		return nil
	})(newContext("GET", "/gino"))
	assert.NotNil(t, err)
	assert.Equal(t, "baseline", err.Error())
}

func newContext(method string, path string) echo.Context {
	server := echo.New()
	r, _ := http.NewRequest(method, "http://example.com"+path, nil)
	ctx := server.NewContext(r, nil)
	ctx.SetPath(path)
	return ctx
}
