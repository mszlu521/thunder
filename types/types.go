package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type ArrayString []string

func (a *ArrayString) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var as ArrayString
	err := json.Unmarshal(bytes, &as)
	*a = as
	return err
}
func (a ArrayString) Value() (driver.Value, error) {
	if a == nil || len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

func (*ArrayString) GormDataType() string {
	return "varchar"
}

type Int64String int64

// MarshalJSON 实现 MarshalJSON 接口
func (i Int64String) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(i), 10) + `"`), nil
}

// UnmarshalJSON 实现 UnmarshalJSON 接口
func (i *Int64String) UnmarshalJSON(data []byte) error {
	strData := string(data)
	if strings.Contains(strData, "\"") {
		strData = strings.Replace(strData, "\"", "", -1)
	}
	parsedInt, err := strconv.ParseInt(strData, 10, 64)
	if err != nil {
		return err
	}
	*i = Int64String(parsedInt)
	return nil
}
