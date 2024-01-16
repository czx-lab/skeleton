package request

import (
	"github.com/go-playground/validator/v10"
	"skeleton/internal/request"
)

type Foo struct {
	Name  int    `binding:"required" form:"name" query:"name" json:"name"`
	Token string `header:"token" binding:"required"`
}

func (f Foo) Message() validator.ValidationErrorsTranslations {
	return map[string]string{
		"Name.required": "name 必填",
	}
}

var _ request.IValidator = (*Foo)(nil)
