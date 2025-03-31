package apperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ngoctd314/common/gvalidator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sentinel error
var (
	DataNotFound   error = New("the requested resource could not be found")
	FilterRequired error = New("a query filter is required")
	Conflict       error = New("the requested resource already exist")
)

// ErrBindRequest used when the request was malformed or contained invalid parameters
func ErrBindRequest(err error) *HTTPError {
	httpErr := &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: "the request was malformed or contained invalid parameters",
		},
		ErrType:  "bind_fail",
		HTTPCode: http.StatusBadRequest,
	}

	var castErr *json.UnmarshalTypeError
	if errors.As(err, &castErr) {
		httpErr.Validation = []ValidatorField{
			{
				Field:   castErr.Field,
				Message: fmt.Sprintf("cannot cast %s into field %s of type %s", castErr.Value, castErr.Field, castErr.Type.String()),
			},
		}
	} else {
		httpErr.message = err.Error()
	}

	return httpErr
}

type ValidatorField struct {
	Field   string `json:"field"`
	Value   any    `json:"value,omitempty"`
	Message string `json:"message"`
}

func (v ValidatorField) Error() string {
	data, _ := json.Marshal(v)
	return string(data)
}

func ErrValidation(err error) *HTTPError {
	httpErr := &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: "the request was malformed or contained invalid parameters",
		},
		ErrType:  "validate_fail",
		HTTPCode: http.StatusBadRequest,
	}

	var vErr validator.ValidationErrors
	if errors.As(err, &vErr) {
		validationErr := make([]ValidatorField, 0, len(vErr))
		for _, e := range vErr {
			validationErr = append(validationErr, ValidatorField{
				Field:   e.Field(),
				Value:   e.Value(),
				Message: e.Translate(gvalidator.GetTranslator("en")),
			})
		}
		httpErr.Validation = validationErr
	} else {
		httpErr.Validation = err.Error()
	}

	return httpErr
}

func ErrBadRequest(message string) *HTTPError {
	return &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: message,
		},
		ErrType:  "bad_request",
		HTTPCode: http.StatusBadRequest,
	}
}

func ErrNotFound(message string) *HTTPError {
	return &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: message,
		},
		ErrType:  "not_found",
		HTTPCode: http.StatusNotFound,
	}
}

func ErrConflict(message string) *HTTPError {
	return &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: message,
		},
		ErrType:  "conflict",
		HTTPCode: http.StatusConflict,
	}
}

var ErrUnauthorizedAccess error = ErrUnauthorized("you are not authorized to access this resource")

func ErrUnauthorized(message string) *HTTPError {
	return &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: message,
		},
		ErrType:  "unauthorized",
		HTTPCode: http.StatusUnauthorized,
	}
}

var ErrForbiddenAccess error = ErrForbidden("you are not allowed to access this resource")

func ErrForbidden(message string) *HTTPError {
	return &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: message,
		},
		ErrType:  "forbidden",
		HTTPCode: http.StatusForbidden,
	}
}

func ErrInternalServer(err error) *HTTPError {
	return &HTTPError{
		BaseError: BaseError{
			ID:      errID(),
			message: "an internal server error occurred, please contact the system administrator",
		},
		ErrType:  "internal_server",
		HTTPCode: http.StatusInternalServerError,
	}
}

func ErrGRPC(err error) *HTTPError {
	errStatus, ok := status.FromError(err)
	if !ok {
		return ErrInternalServer(err)
	}

	switch errStatus.Code() {
	case codes.InvalidArgument:
		return ErrBadRequest(errStatus.Message())
	case codes.NotFound:
		return ErrNotFound(errStatus.Message())
	case codes.AlreadyExists:
		return ErrConflict(errStatus.Message())
	default:
		return ErrBadRequest(errStatus.Message())
	}
}
