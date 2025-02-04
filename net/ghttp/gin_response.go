package ghttp

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/common/apperror"
)

func JSONSuccess(c *gin.Context, respDTO *ResponseBody) {
	if respDTO.StatusCode == 0 {
		respDTO.StatusCode = http.StatusOK
	}

	c.JSON(respDTO.StatusCode, respDTO)
}

func JSONFail(c *gin.Context, err error) {
	if canResolveErr(c, err) {
		return
	}

	httpError := apperror.ErrInternalServer(err)
	logErr(c, &httpError.BaseError, err)

	c.JSON(httpError.HTTPCode, ResponseBody{
		Success: false,
		Error:   httpError,
		Message: httpError.Error(),
	})
}

func JSONAbort(c *gin.Context, err error) {
	c.Abort()
	JSONFail(c, err)
}

func logErr(c *gin.Context, baseErr *apperror.BaseError, err error) {
	kvs := []any{"err_id", baseErr.ID, "path", c.FullPath()}
	if baseErr.Ancestor() != nil {
		kvs = append(kvs, "ancestor", baseErr.Ancestor())
	}
	slog.ErrorContext(c.Request.Context(), err.Error(), kvs...)
}

func canResolveErr(c *gin.Context, err error) bool {
	var baseErr *apperror.BaseError
	if errors.As(err, &baseErr) {
		httpErr := apperror.NewHTTPError(baseErr, http.StatusBadRequest)
		httpErr.SetErrType("bad_request")
		c.JSON(httpErr.HTTPCode, ResponseBody{
			Success: false,
			Error:   httpErr,
			Message: err.Error(),
		})
		logErr(c, baseErr, err)
		return true
	}

	var httpErr *apperror.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.HTTPCode == 0 {
			httpErr.HTTPCode = http.StatusBadRequest
		}
		c.JSON(httpErr.HTTPCode, ResponseBody{
			Success: false,
			Error:   httpErr,
			Message: httpErr.Error(),
		})
		logErr(c, &httpErr.BaseError, err)
		return true
	}

	return false
}
