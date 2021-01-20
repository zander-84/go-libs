package errs

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	var a error
	a = New("test").SetCode("100")
	var d = New("")
	fmt.Println(errors.Is(a, d))

}
