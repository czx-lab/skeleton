package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type Decimal = decimal.Decimal

type ModelFieldTime struct {
	time.Time
}

func (t ModelFieldTime) MarshalJSON() ([]byte, error) {
	seconds := t.Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}

func (t ModelFieldTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *ModelFieldTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = ModelFieldTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
