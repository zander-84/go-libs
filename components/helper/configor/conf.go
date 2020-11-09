package configor

import "github.com/jinzhu/configor"

func LoadConf(ptrValue interface{}, files ...string) (interface{}, error) {
	err := configor.Load(ptrValue, files...)
	return ptrValue, err
}
