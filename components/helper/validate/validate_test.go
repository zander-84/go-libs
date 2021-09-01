package validate

import (
	"fmt"
	"testing"
)

type user2 struct {
	UserName string `form:"name" json:"name"  validate:"required,min=5"  comment:"用户名2"`
	Password string `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"密码2"`
}

// go test -v  -run TestValidate
func TestValidate(t *testing.T) {
	v := NewValidate(Conf{
		Locale:      "zh",
		ValidateTag: "validate",
	})

	type User struct {
		UserName string   `form:"name" json:"name"  validate:"required,min=5"  comment:"用户名"`
		Password string   `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"密码"`
		U        *[]user2 `form:"u" json:"u" xml:"u" validate:"required,dive,required" comment:"密码"`
	}

	u := User{
		UserName: "z",
		Password: "aaaaaaa",
	}
	//test1(v,u)
	test2(v, u)
}

func test2(v *Validate, u interface{}) {
	if err := v.ValidateStruct(u); err != nil {
		fmt.Println(err.Error())
	}
}
func test1(v *Validate, u interface{}) {
	if err := v.ValidateStruct(u); err != nil {
		fmt.Println(err.Error())
	}
}

//go test -v  -bench=BenchmarkValidate1
// BenchmarkValidate1-8      275316              4080 ns/op
func BenchmarkValidate1(b *testing.B) {
	v := NewValidate(Conf{
		Locale:      "zh",
		ValidateTag: "validate",
	})

	type User struct {
		Name     string `form:"name" json:"name"  validate:"required,min=5"  comment:"name:用户名"`
		Password string `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"password:密码"`
	}
	u := User{
		Name:     "z",
		Password: "123",
	}

	for i := 0; i < b.N; i++ {
		test1(v, u)
	}
}

//go test -v  -bench=BenchmarkValidate2   性能损耗在翻译
// BenchmarkValidate2-8     583309               1782 ns/op   一个错误    原生：   BenchmarkValidate2-8     1442572               796 ns/op
// BenchmarkValidate2-8     409200               2848 ns/op   两个错误            BenchmarkValidate2-8      975138              1098 ns/op
func BenchmarkValidate2(b *testing.B) {
	v := NewValidate(Conf{
		Locale:      "zh",
		ValidateTag: "validate",
	})

	type User struct {
		Name     string `form:"name" json:"name"  validate:"required,min=5"  comment:"name:用户名"`
		Password string `form:"password" json:"password" xml:"password" validate:"required,max=4" comment:"password:密码"`
	}
	u := User{
		Name:     "aaaaaaaa",
		Password: "aaaaaaaa",
	}

	for i := 0; i < b.N; i++ {
		test2(v, u)
	}
}
