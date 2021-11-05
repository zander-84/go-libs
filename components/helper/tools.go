package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var tool = NewTool()

type Tool struct{}

func NewTool() *Tool { return new(Tool) }

func GetTool() *Tool { return tool }

func PrettyPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return
	}

	fmt.Println(out.String())
}
