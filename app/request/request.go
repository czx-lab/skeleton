package request

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	AppRequest "skeleton/internal/request"
)

type Request struct {
	validator *AppRequest.Request
}

func New() (*Request, error) {
	appRequest, err := AppRequest.New("zh")
	if err != nil {
		return nil, err
	}
	return &Request{
		validator: appRequest,
	}, nil
}

func (r Request) Validator(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	return r.validator.Validator(ctx, param)
}
