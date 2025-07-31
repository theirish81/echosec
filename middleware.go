package echosec

import (
	"context"
	"errors"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

// OApiValidationFunc is a function meant to validate whether the requesting entity has the permission to access
// a certain resource
type OApiValidationFunc func(c echo.Context, params []string) error

// WithOpenApiConfig creates an echo.MiddlewareFunc function given an OApiConfig object
func WithOpenApiConfig(cfg OApiConfig) echo.MiddlewareFunc {
	evaluator := newEvaluator(cfg.vars)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			route, pathParams, err := cfg.router.FindRoute(c.Request())
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
			ectx := EchoSecContext{Config: localConfig}
			setDoubleContext(c, EchosecContextAttr, ectx)
			params := localConfig.Params
			if params == nil {
				params = make([]string, 0)
			}
			err = f(c, params)
			if err != nil {
				return err
			}
			labels, err := evaluator.Eval(c, route.Operation.OperationID, localConfig.Labels)
			if err != nil {
				return err
			}
			ectx.Labels = labels
			setDoubleContext(c, EchosecContextAttr, ectx)

			if cfg.openApiValidationEnabled {
				requestValidationInput := &openapi3filter.RequestValidationInput{
					Request:    c.Request(),
					PathParams: pathParams,
					Route:      route,
					Options: &openapi3filter.Options{
						AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
							return nil
						},
					},
				}
				c.Set(requestValidationInputAttr, requestValidationInput)
				err = openapi3filter.ValidateRequest(c.Request().Context(), requestValidationInput)
				if err != nil {
					return err
				}
				return next(c)
			}
			return next(c)
		}
	}
}

type Evaluator struct {
	compiled map[string]*vm.Program
	vars     []string
}

func newEvaluator(vars []string) *Evaluator {
	return &Evaluator{
		compiled: make(map[string]*vm.Program),
		vars:     vars,
	}
}

func (c *Evaluator) scope(ctx echo.Context) map[string]any {
	base := map[string]any{"ctx": ctx}
	for _, v := range c.vars {
		base[v] = ctx.Get(v)
	}
	return base
}

func (c *Evaluator) Eval(ctx echo.Context, operationId string, labels []Label) ([]string, error) {
	out := make([]string, 0)
	if labels == nil {
		return out, nil
	}
	for _, l := range labels {
		if v, ok := c.compiled[operationId+"_"+l.Label]; ok {
			r, err := expr.Run(v, c.scope(ctx))
			if err != nil {
				return nil, err
			}
			if v, ok := r.(bool); ok && v {
				out = append(out, l.Label)
			}
		} else {
			cx, err := expr.Compile(l.Condition)
			if err != nil {
				return nil, err
			}
			c.compiled[operationId+"_"+l.Label] = cx
			r, err := expr.Run(cx, c.scope(ctx))
			if err != nil {
				return nil, err
			}
			if v, ok := r.(bool); ok && v {
				out = append(out, l.Label)
			}
		}
	}
	return out, nil
}

// setDoubleContext sets a context k/v to both the echo and Golang context
func setDoubleContext(ectx echo.Context, key string, val any) {
	ectx.Set(key, val)
	ctx := ectx.Request().Context()
	ctx = context.WithValue(ctx, key, val)
	ectx.SetRequest(ectx.Request().WithContext(ctx))
}
