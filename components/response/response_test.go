package response

import (
	"errors"
	"fmt"
	"github.com/zander-84/go-libs/components/errs"
	"testing"
)

// go test -v  -run TestResponse
func TestResponse(t *testing.T) {
	a := errs.NewCustomError("cuowu", errs.ErrorLevel)
	f := errors.New("fffff")
	testResponse(a)
	testResponse(f)

}

func testResponse(err error) {
	if a, ok := err.(errs.CustomError); ok {
		fmt.Println(a.IsSuccess())
	} else {
		fmt.Println("error")
	}
}
