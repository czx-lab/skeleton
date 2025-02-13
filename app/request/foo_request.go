package request

import (
	"skeleton/internal/request"

	"github.com/go-playground/validator/v10"
)

type Foo struct {
	Header struct {
		Token string `header:"token" binding:"required"`
	}
	Body struct {
		Name   string `form:"name" binding:"required"`
		Number int    `form:"number" binding:"required,gte=1,lte=100"`
	}
}

func (f Foo) Message() validator.ValidationErrorsTranslations {
	return map[string]string{
		"Header.Token.required": "header 参数 token 必填",
		"Body.Name.required":    "name 必填",
		"Body.Number.required":  "number 必填",
		"Body.Number.gte":       "number >= 1",
		"Body.Number.lte":       "number <= 100",
	}
}

var _ request.IValidator = (*Foo)(nil)

type FooBody struct {
	Name   string `form:"name" binding:"required"`
	Number int    `form:"number" binding:"required,gte=1,lte=100"`
}

func (f FooBody) Message() validator.ValidationErrorsTranslations {
	return map[string]string{
		"Name.required":   "name 必填",
		"Number.required": "number 必填",
		"Number.gte":      "number >= 1",
		"Number.lte":      "number <= 100",
	}
}
