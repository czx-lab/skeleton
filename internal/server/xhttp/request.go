package xhttp

import (
	"errors"
	"fmt"
	"reflect"
	"czx/internal/utils"
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
	headerName = "Header"
	bodyName   = "Body"
	uriName    = "Uri"
)

type IValidator interface {
	Message() validator.ValidationErrorsTranslations
}

type Request struct {
	trans ut.Translator
}

func NewReq(local string) (*Request, error) {
	req := &Request{}
	trans, err := trans(local)
	if err != nil {
		return nil, err
	}
	req.trans = trans
	return req, nil
}

func (r *Request) Validator(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	var err error
	checkHeader := utils.CheckFieldExistence(param, headerName)
	if !checkHeader {
		goto CheckUriBlock
	}
	if err := r.ValiHeader(ctx, param); err != nil {
		return err
	}

CheckUriBlock:
	checkUri := utils.CheckFieldExistence(param, uriName)
	if !checkUri {
		goto CheckBodyBlock
	}
	if err := r.ValiUri(ctx, param); err != nil {
		return err
	}

CheckBodyBlock:
	field := bodyName
	checkBody := utils.CheckFieldExistence(param, field)
	if checkBody {
		if err := r.ValiBody(ctx, param); err != nil {
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
		return r.valiErr(field, param, err)
	}
	return nil
}

func (r *Request) ValiHeader(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	headerVal := reflectValue(param).FieldByName(headerName)
	if headerVal.CanInterface() {
		if err := ctx.ShouldBindHeader(headerVal.Addr().Interface()); err != nil {
			return r.valiErr(headerName, param, err)
		}
	}
	return nil
}

func (r *Request) ValiBody(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	bodyVal := reflectValue(param).FieldByName(bodyName)
	if bodyVal.CanInterface() {
		if err := ctx.ShouldBind(bodyVal.Addr().Interface()); err == nil {
			return nil
		}
	}
	return nil
}

func (r *Request) ValiUri(ctx *gin.Context, param any) validator.ValidationErrorsTranslations {
	uriVal := reflectValue(param).FieldByName(uriName)
	if uriVal.CanInterface() {
		if err := ctx.ShouldBindUri(uriVal.Addr().Interface()); err != nil {
			return r.valiErr(headerName, param, err)
		}
	}
	return nil
}

func (r *Request) valiErr(field string, param any, err error) validator.ValidationErrorsTranslations {
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

func reflectValue(param any) reflect.Value {
	value := reflect.ValueOf(param)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func trans(local string) (ut.Translator, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		trans, ok := uni.GetTranslator(local)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		var err error
		switch local {
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return trans, err
	}
	return nil, nil
}
