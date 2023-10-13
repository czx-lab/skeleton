package utils

import "strings"

// CaseToCamel 下划线转驼峰(大驼峰)
func CaseToCamel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// LowerCamelCase 转换为小驼峰
func LowerCamelCase(name string) string {
	name = CaseToCamel(name)
	return strings.ToLower(name[:1]) + name[1:]
}
