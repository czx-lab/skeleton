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

type IValidator interface {
	Message() validator.ValidationErrorsTranslations
}

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
	messages := make(validator.ValidationErrorsTranslations)
	if ok := errors.As(err, &errs); !ok {
		messages["noValidationErrors"] = err.Error()
		return messages
	}
	for _, fieldError := range err.(validator.ValidationErrors) {
		field := fmt.Sprintf("%s.%s", fieldError.Field(), fieldError.Tag())
		d, ok := param.(IValidator)
		if !ok {
			messages[fieldError.Field()] = fieldError.Translate(r.trans)
			continue
		}
		if message, exist := d.Message()[field]; exist {
			messages[fieldError.Field()] = message
			continue
		}
		messages[fieldError.Field()] = fieldError.Error()
	}

	return messages
}
