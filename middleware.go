package echosec

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

// OApiValidationFunc is a function meant to validate whether the requesting entity has the permission to access
// a certain resource
type OApiValidationFunc func(c echo.Context, params []string) error

// WithOpenApiConfig creates an echo.MiddlewareFunc function given an OApiConfig object
func WithOpenApiConfig(cfg OApiConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			route, _, err := cfg.router.FindRoute(c.Request())
			if err != nil {
				return err
			}
			sec := route.Operation.Extensions["x-echosec"]
			var localConfig OApiEchoSec
			err = mapstructure.Decode(sec, &localConfig)
			if err != nil {
				return errors.New("error in local configuration")
			}
			f := cfg.validators[localConfig.Function]
			if f == nil {
				return errors.New("validator function not found")
			}
			params := localConfig.Params
			if params == nil {
				params = make([]string, 0)
			}
			err = f(c, params)
			if err != nil {
				return err
			}
			return next(c)
		}
	}
}
