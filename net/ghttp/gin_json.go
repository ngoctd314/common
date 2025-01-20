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
	if isInternalErr(c, err) {
		return
	}

	httpError := apperror.ErrInternalServer(err)
	logBaseErr(&httpError.BaseError)

	c.JSON(httpError.HTTPCode, ResponseBody{
		Success: false,
		Message: httpError.Error(),
	})
}

func JSONAbort(c *gin.Context, err error) {
	c.Abort()
	JSONFail(c, err)
}

func logBaseErr(baseErr *apperror.BaseError) {
	kvs := []any{"err_id", baseErr.ID}
	if baseErr.Ancestor() != nil {
		kvs = append(kvs, "err", baseErr.Ancestor())
	}
	slog.Error(baseErr.Error(), kvs...)
}

func isInternalErr(c *gin.Context, err error) bool {
	var baseErr *apperror.BaseError
	if errors.As(err, &baseErr) {
		httpErr := apperror.NewHTTPError(err, http.StatusBadRequest)
		c.JSON(httpErr.HTTPCode, ResponseBody{
			Success: false,
			Error:   err,
			Message: baseErr.Error(),
		})
		logBaseErr(baseErr)
		return true
	}

	var httpErr *apperror.HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.HTTPCode, ResponseBody{
			Success: false,
			Error:   err,
			Message: httpErr.Error(),
		})
		logBaseErr(&httpErr.BaseError)
		return true
	}

	return false
}
