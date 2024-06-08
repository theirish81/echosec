package echosec

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SecBlob(ctx echo.Context, status int, dto any) error {
	var data []byte
	switch t := dto.(type) {
	case string:
		data = []byte(t)
	case []byte:
		data = t
	default:
		data, _ = json.Marshal(dto)
	}
	rvi, ok := ctx.Get(requestValidationInputAttr).(*openapi3filter.RequestValidationInput)
	h := http.Header{}
	h.Set(contentTypeHeaderName, contentTypeJsonValue)
	if ok {
		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: rvi,
			Status:                 status,
			Header:                 h,
		}
		responseValidationInput.SetBodyBytes(data)
		err := openapi3filter.ValidateResponse(ctx.Request().Context(), responseValidationInput)
		if err != nil {
			return err
		}
	}
	return ctx.JSONBlob(status, data)
}
