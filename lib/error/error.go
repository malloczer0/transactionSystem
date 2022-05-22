package error

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"transactionSystemTestTask/lib/types"
)

// HTTPError is our HTTP error definition.
type HTTPError struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Cause   string `json:"cause,omitempty"`
}

// Error implements a custom echo error handler that will encode errors as JSON objects rather than
// just return a text body. It will also make sure to not have the redundant information given in
// the echo string encoding of HTTP errors.
func Error(err error, context echo.Context) {
	httpErrorObject := HTTPError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	switch err {
	case types.ErrBadRequest:
		httpErrorObject.Code = http.StatusBadRequest
	case types.ErrNotFound:
		httpErrorObject.Code = http.StatusNotFound
	case types.ErrDuplicateEntry, types.ErrConflict:
		httpErrorObject.Code = http.StatusConflict
	case types.ErrForbidden:
		httpErrorObject.Code = http.StatusForbidden
	case types.ErrUnprocessableEntity:
		httpErrorObject.Code = http.StatusUnprocessableEntity
	case types.ErrPartialOk:
		httpErrorObject.Code = http.StatusPartialContent
	case types.ErrGone:
		httpErrorObject.Code = http.StatusGone
	case types.ErrUnauthorized:
		httpErrorObject.Code = http.StatusUnauthorized
	}
	httpError, ok := err.(*echo.HTTPError)
	if ok {
		httpErrorObject.Code = httpError.Code
		httpErrorObject.Message = fmt.Sprintf("%v", httpError.Message)
	}
	httpErrorObject.Name = http.StatusText(httpErrorObject.Code)
	if !context.Response().Committed {
		if context.Request().Method == echo.HEAD {
			context.NoContent(httpErrorObject.Code)
		} else {
			context.JSON(httpErrorObject.Code, httpErrorObject)
		}
	}
}
