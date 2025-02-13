package memo

import (
	"czx/internal/constants"
	"strings"
	"sync"
)

var smap sync.Map

type IMemo interface {
	FuzzyDel(string) error
}

type Memo struct {
}

func New() *Memo {
	return &Memo{}
}

// Get implements constants.ICache.
func (m *Memo) Get(key string) (any, error) {
	value, exist := smap.Load(key)
	if exist {
		return value, nil
	}
	return nil, nil
}

// Has implements constants.ICache.
func (m *Memo) Has(key string) (bool, error) {
	_, ok := smap.Load(key)
	return ok, nil
}

// Set implements constants.ICache.
func (m *Memo) Set(key string, value any) (res bool, err error) {
	ok, err := m.Has(key)
	if err != nil {
		return
	}
	if !ok {
		smap.Store(key, value)
		res = true
	}
	return
}

func (m *Memo) FuzzyDel(prefix string) error {
	smap.Range(func(key, value any) bool {
		if key, ok := key.(string); ok {
			if strings.HasPrefix(key, prefix) {
				smap.Delete(key)
			}
		}
		return true
	})
	return nil
}

var _ constants.ICache = (*Memo)(nil)
