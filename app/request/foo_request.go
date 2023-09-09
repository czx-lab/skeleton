package request

type Foo struct {
	Name  int    `binding:"required" form:"name" query:"name" json:"name"`
	Token string `header:"token" binding:"required"`
}
