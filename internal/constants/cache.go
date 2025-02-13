package constants

type ICache interface {
	Get(key string) (any, error)
	Set(key string, value any) (bool, error)
	Has(key string) (bool, error)
}
