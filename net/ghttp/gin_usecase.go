package ghttp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/gvalidator"
)

type Usecase[Req any] interface {
	Usecase(ctx context.Context, req *Req) (*ResponseBody, error)
}

type Binding[Req any] interface {
	Bind(c *gin.Context) (*Req, error)
}

type Validating[Req any] interface {
	Validate(ctx context.Context, req *Req) error
}

func GinHandleFunc[Req any](uc Usecase[Req]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			req Req
			err error
		)

		defer func() {
			if err != nil {
				// set req, rid into context
				JSONFail(c, err)
			}
		}()

		// bind request
		if binder, isBinder := uc.(Binding[Req]); isBinder {
			bindReq, bindErr := binder.Bind(c)
			if bindErr != nil {
				err = apperror.ErrBindRequest(bindErr)
				return
			}
			req = *bindReq
		} else {
			// TODO: check it
			if bindErr := c.ShouldBindUri(&req); bindErr != nil {
				err = apperror.ErrBindRequest(bindErr)
				return
			}
			if c.Request.Method != http.MethodDelete {
				if bindErr := c.ShouldBind(&req); bindErr != nil {
					err = apperror.ErrBindRequest(bindErr)
					return
				}
			}
		}

		ctx := c.Request.Context()
		// validate req
		if validator, isValidator := uc.(Validating[Req]); isValidator {
			if validateErr := validator.Validate(ctx, &req); validateErr != nil {
				err = apperror.ErrValidation(validateErr)
				return
			}
		}
		// always validate the request
		if validateErr := gvalidator.ValidateStruct(req); validateErr != nil {
			err = apperror.ErrValidation(validateErr)
			return
		}

		resp, usecaseErr := uc.Usecase(ctx, &req)
		if usecaseErr != nil {
			err = usecaseErr
			return
		}

		JSONSuccess(c, resp)
	}
}

type RouteGroup struct {
	router *gin.RouterGroup
}
