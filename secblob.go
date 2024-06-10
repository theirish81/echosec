package echosec

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

// SecBlob validates whether the provided dto matches the output expected by OpenAPI specification.
// For this to work, validation of the input must also be enabled
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
	headers := ctx.Response().Header()
	if headers.Get(contentTypeHeaderName) == "" {
		headers.Set(contentTypeHeaderName, contentTypeJsonValue)
		ctx.Response().Header().Set(contentTypeHeaderName, contentTypeJsonValue)
	}

	rvi, ok := ctx.Get(requestValidationInputAttr).(*openapi3filter.RequestValidationInput)
	if ok {
		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: rvi,
			Status:                 status,
			Header:                 headers,
		}
		responseValidationInput.SetBodyBytes(data)
		err := openapi3filter.ValidateResponse(ctx.Request().Context(), responseValidationInput)
		if err != nil {
			return err
		}
	}
	return ctx.JSONBlob(status, data)
}
