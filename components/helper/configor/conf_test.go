package configor

import (
	"testing"
)

// go test -v  -run TestNewConf
func TestNewConf(t *testing.T) {
	var Config = struct {
		APPName string `default:"app name"`

		DB struct {
			Name     string `default:"root"`
			User     string `default:"root"`
			Password string `required:"true" env:"DBPassword"`
			Port     uint   `default:"3306"`
		}

		Contacts []struct {
			Name  string
			Email string `required:"true"`
		}
	}{}

	if _, err := LoadConf(&Config, "conf1.yml", "conf2.yml"); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(Config)
	}

}
