package ghttp

import (
	"context"

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
		var req Req
		var err error

		defer func() {
			if err != nil {
			}
		}()

		// bind request
		if binder, isBinder := uc.(Binding[Req]); isBinder {
			bindReq, bindErr := binder.Bind(c)
			if bindErr != nil {
				err = apperror.ErrBindRequest(err)
				return
			}
			req = *bindReq
		} else {
			if bindErr := c.ShouldBind(&req); bindErr != nil {
				err = apperror.ErrBindRequest(err)
				return
			}
			if bindErr := c.ShouldBindUri(&req); bindErr != nil {
				err = apperror.ErrBindRequest(err)
				return
			}
		}

		ctx := c.Request.Context()
		// validate req
		if validator, isValidator := uc.(Validating[Req]); isValidator {
			if validateErr := validator.Validate(ctx, &req); validateErr != nil {
				err = apperror.ErrValidation(validateErr)
				return
			}
		} else {
			if validateErr := gvalidator.ValidateStruct(req); validateErr != nil {
				err = apperror.ErrValidation(validateErr)
				return
			}
		}

		resp, usecaseErr := uc.Usecase(c, &req)
		if usecaseErr != nil {
			err = usecaseErr
			return
		}

		JSONSuccess(c, resp)
	}
}
