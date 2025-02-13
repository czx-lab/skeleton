package request

import (
	"errors"
	"fmt"
	"reflect"
	"skeleton/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

const (
	headerFieldName = "Header"
	bodyFieldName   = "Body"
	uriFieldName    = "Uri"
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

func paramReflectValue(param any) reflect.Value {
	paramValue := reflect.ValueOf(param)
	if paramValue.Kind() == reflect.Ptr {
		paramValue = paramValue.Elem()
	}
	return paramValue
}

func (r *Request) Validator(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	var err error
	checkHeader := utils.CheckFieldExistence(param, headerFieldName)
	if !checkHeader {
		goto CheckUriBlock
	}
	if err := r.valiHeader(ctx, param); err != nil {
		return err
	}

CheckUriBlock:
	checkUri := utils.CheckFieldExistence(param, uriFieldName)
	if !checkUri {
		goto CheckBodyBlock
	}
	if err := r.valiUri(ctx, param); err != nil {
		return err
	}

CheckBodyBlock:
	field := bodyFieldName
	checkBody := utils.CheckFieldExistence(param, bodyFieldName)
	if checkBody {
		if err := r.valiBody(ctx, param); err != nil {
			return err
		}
	}
	if !checkBody && !checkHeader && !checkUri {
		field = ""
		if err = ctx.ShouldBind(param); err == nil {
			return nil
		}
	}
	if err != nil {
		return r.valiError(field, param, err)
	}
	return nil
}

func (r *Request) valiBody(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	bodyVal := paramReflectValue(param).FieldByName(bodyFieldName)
	if bodyVal.CanInterface() {
		if err := ctx.ShouldBind(bodyVal.Addr().Interface()); err == nil {
			return nil
		}
	}
	return nil
}

func (r *Request) valiHeader(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	headerVal := paramReflectValue(param).FieldByName(headerFieldName)
	if headerVal.CanInterface() {
		if err := ctx.ShouldBindHeader(headerVal.Addr().Interface()); err != nil {
			return r.valiError(headerFieldName, param, err)
		}
	}
	return nil
}

func (r *Request) valiUri(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	uriVal := paramReflectValue(param).FieldByName(uriFieldName)
	if uriVal.CanInterface() {
		if err := ctx.ShouldBindUri(uriVal.Addr().Interface()); err != nil {
			return r.valiError(headerFieldName, param, err)
		}
	}
	return nil
}

func (r *Request) valiError(field string, param any, err error) validator.ValidationErrorsTranslations {
	var errs validator.ValidationErrors
	messages := make(validator.ValidationErrorsTranslations)
	if ok := errors.As(err, &errs); !ok {
		messages["noValidationErrors"] = err.Error()
		return messages
	}
	for _, fieldError := range err.(validator.ValidationErrors) {
		stringBuilder := strings.Builder{}
		if len(field) > 0 {
			stringBuilder.WriteString(fmt.Sprintf("%s.", field))
		}
		stringBuilder.WriteString(fmt.Sprintf("%s.%s", fieldError.Field(), fieldError.Tag()))
		d, ok := param.(IValidator)
		if !ok {
			messages[fieldError.Field()] = fieldError.Translate(r.trans)
			continue
		}
		if message, exist := d.Message()[stringBuilder.String()]; exist {
			messages[fieldError.Field()] = message
			continue
		}
		messages[fieldError.Field()] = fieldError.Error()
	}

	return messages
}
