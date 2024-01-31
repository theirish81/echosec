package echosec

import "github.com/labstack/echo/v4"

// ValidationFunc is any function meant to validate access to a path or method
type ValidationFunc func(c echo.Context) error

// Middleware returns an echo.MiddlewareFunc for the echo server
// cfg is the Config of the middleware
func Middleware(cfg Config) echo.MiddlewareFunc {
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
