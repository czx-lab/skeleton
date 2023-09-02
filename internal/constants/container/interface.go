package container

// ContainerInterface 容器相关方法
type ContainerInterface interface {
	Get(key string) any
	Set(key string, value any) bool
	Has(key string) bool
}
