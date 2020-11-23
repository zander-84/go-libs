package mysql

import (
	"errors"
	"strings"
)

type Mysql interface {
	Engine() interface{}
	Transaction(f func(tx interface{}) (int, error)) (int, error)
}

type BasicModel struct {
	DisplayFields []string `json:"-" gorm:"-"`
	IsCache       bool     `json:"-" gorm:"-"`
}

func (this *BasicModel) MarshalJSON() ([]byte, error) {

	return nil, errors.New("You Should Do MarshalJSON !!! ")
}

func (this *BasicModel) FilterFields(data map[string]interface{}) map[string]interface{} {
	if this.DisplayFields == nil {
		return data
	} else {
		for key, _ := range data {
			exist := false
			for _, key2 := range this.DisplayFields {
				if key == key2 {
					exist = true
					break
				}
			}
			if !exist {
				delete(data, key)
			}
		}
		return data
	}
}

func (this *BasicModel) Fields(tableName string, fileds []string) []string {
	for k, v := range fileds {
		if !strings.Contains(v, ".") {
			fileds[k] = "`" + tableName + "`" + "." + v
		}
	}

	return fileds
}
