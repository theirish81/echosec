package echosec

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathItem_MatchPattern(t *testing.T) {
	items := PathItems{
		{
			Patterns: Patterns{"/user/:id"},
		},
		{
			Patterns: Patterns{"/foo/bar", "/a/b/c"},
		},
	}
	assert.True(t, items[0].MatchPattern("/user/:id"))
	assert.False(t, items[0].MatchPattern("/user/abc"))

	assert.True(t, items[1].MatchPattern("/foo/bar"))
	assert.True(t, items[1].MatchPattern("/a/b/c"))
	assert.False(t, items[1].MatchPattern("/a/b/c/d"))
}

func TestPathItem_FindMethodValidator(t *testing.T) {
	item := PathItem{
		Methods: ValidationMap{
			"GET": func(c echo.Context) error {
				return nil
			},
			"PUT,PATCH": func(c echo.Context) error {
				return nil
			},
		},
	}
	assert.NotNil(t, item.FindMethodValidator("get"))
	assert.NotNil(t, item.FindMethodValidator("put"))
	assert.NotNil(t, item.FindMethodValidator("patch"))
	assert.Nil(t, item.FindMethodValidator("delete"))
}
