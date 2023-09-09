package request

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type Request struct {
	trans ut.Translator
}

func New(local string) (*Request, error) {
	requestClass := &Request{}
	if err := requestClass.transInit(local); err != nil {
		return nil, err
	}
	return requestClass, nil
}

func (r *Request) transInit(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		var o bool
		r.trans, o = uni.GetTranslator(local)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		switch local {
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, r.trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, r.trans)
		}
		return
	}
	return
}

func (r *Request) Validator(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	err := ctx.ShouldBind(param)
	if err == nil {
		return nil
	}
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		errMap := make(validator.ValidationErrorsTranslations)
		errMap["noValidationErrors"] = err.Error()
		return errMap
	}
	return errs.Translate(r.trans)
}
