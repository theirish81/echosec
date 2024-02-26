package echosec

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

// ValidationFunc is any function meant to validate access to a path or method
type ValidationFunc func(c echo.Context) error

type OApiValidationFunc func(c echo.Context, params []string) error

// WithManualConfig returns an echo.MiddlewareFunc for the echo server
// cfg is the ManualConfig of the middleware
func WithManualConfig(cfg ManualConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, p := range cfg.PathMapping {
				if p.MatchPattern(c.Path(), cfg.BasePath) {
					if validator := p.FindMethodValidator(c.Request().Method); validator != nil {
						if err := validator(c); err != nil {
							return err
						}
						return next(c)
					}
					if p.PathValidation != nil {
						if err := p.PathValidation(c); err != nil {
							return err
						}
						return next(c)
					}
				}
			}
			if cfg.DefaultValidation != nil {
				if err := cfg.DefaultValidation(c); err != nil {
					return err
				}
			}
			return next(c)
		}
	}
}

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
