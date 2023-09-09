package request

import (
	AppRequest "github.com/czx-lab/skeleton/internal/request"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
