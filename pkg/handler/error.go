package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gedex/simdoc/pkg/util/upload"

	validator "gopkg.in/validator.v2"
)

type fieldError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

// fieldErrorCode represents error code of a field
type fieldErrorCode int

const (
	// Field related error codes
	ErrorFieldMissing fieldErrorCode = iota
	ErrorFieldInvalid
	ErrorFieldAlreadyExists
	ErrorFieldImmutable
)

// fieldErrorText represents string of fieldErrorCode
var fieldErrorText = [...]string{
	"missing_field",
	"invalid",
	"already_exists",
	"immutable_field",
}

func (fe fieldErrorCode) Error() string {
	return fieldErrorText[fe]
}

func newFieldError(resource, field string, e error) *fieldError {
	return &fieldError{
		Resource: resource,
		Field:    field,
		Code:     e.Error(),
	}
}

type errorResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message,omitempty"`
	Errors  []*fieldError `json:"errors,omitempty"`
}

type errorType int

const (
	ErrorInvalidJSONRequest errorType = iota
	ErrorInternalServerError
	ErrorNotFound
	ErrorRequireAuthentication
	ErrorBadCredentials
	ErrorValidationFailed
	ErrorNotImplemented
	ErrorForbidden
)

var errorText = [...]string{
	"Invalid JSON request",
	"Internal server error",
	"Not found",
	"Requires authentication",
	"Bad credentials",
	"Validation failed",
	"Not implemented",
	"This resource is forbidden",
}

func (e errorType) Error() string {
	return errorText[e]
}

func getValidationErrors(res string, e error) (fe []*fieldError) {
	if v, ok := e.(validator.ErrorMap); ok && len(v) > 0 {
		for k, err := range v {
			field := strings.ToLower(k)
			fe = append(fe, newFieldError(res, field, err))
		}
	}

	return fe
}

func getUploadFilesErrors(files []*upload.File) (fe []map[string]interface{}) {
	for _, f := range files {
		e := map[string]interface{}{
			"name":  f.Name,
			"size":  f.Size,
			"error": f.Error.Error(),
		}
		fe = append(fe, e)
	}
	return
}

func respWithError(w http.ResponseWriter, code int, err error, fieldErrors ...*fieldError) {
	w.WriteHeader(code)

	errResp := &errorResponse{
		Code:    code,
		Message: err.Error(),
	}

	for _, e := range fieldErrors {
		errResp.Errors = append(errResp.Errors, e)
	}

	json.NewEncoder(w).Encode(errResp)
}
