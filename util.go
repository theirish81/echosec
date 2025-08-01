package echosec

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"slices"
)

func GetEchoSecContext(ctx any) (*EchoSecContext, error) {
	switch c := ctx.(type) {
	case context.Context:
		cx := c.Value(EchosecContextAttr)
		if cy, ok := cx.(EchoSecContext); ok {
			return &cy, nil
		}
	case echo.Context:
		return GetEchoSecContext(c.Request().Context())
	}
	return nil, errors.New("EchoSec Context not present")
}

func MustGetEchoSecContext(ctx any) *EchoSecContext {
	cx, _ := GetEchoSecContext(ctx)
	return cx
}

func MustHaveLabels(ctx any, labels ...string) bool {
	c := MustGetEchoSecContext(ctx)
	if c == nil {
		return false
	}
	for _, label := range labels {
		if !slices.Contains(c.Labels, label) {
			return false
		}
	}
	return true
}
