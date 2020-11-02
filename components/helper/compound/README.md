go components
===========================

## 定时任务 Example
```go
package main

import (
		"context"
    	"fmt"
    	"github.com/zander-84/go-components"
    	"net/http"
    	"os"
    	"time"
)


const TEXT1 = `甲方（签章）：                           乙方：xxxxxxxxxxxxxxxxxxxxxxx
甲方代表（签字）：                       乙方盖章：
签订日期： {{now}}           签订日期：  {{now}}

基于对xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
{{{img1}}}{{{img2}}}{{{line}}}
一、xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
二、xxxxxxx　{{level}} xxxxxxx　  　，xxxxxxx　{{money}} 元，xxxxxxx　{{money}} 　  _元；（xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx）
1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
2.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
三、xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx：
1.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
2.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
3.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
4.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
5.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
6.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
四、xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx。
用户确认（手写）：以上条款本人已阅读并清楚了解，同意以上约定。
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx：`


func main(){
	c := C.NewComponents("./")
 	comp := NewCompound()
 	comp.Init(1200, 1200, "./src/simsun.ttf", 18, 72, 1.2)
 	if err := comp.AddTitle("移动业务靓号使用协议", 30); err != nil {
 		fmt.Println("add title error:", err.Error())
 	}
 	if err := comp.HandleBody(TEXT1, data, 18, 1100); err != nil {
 		fmt.Println("HandleBody error:", err.Error())
 	}
 	
 data := make(map[string]Data, 0)
 	data["now"] = Data{
 		Value:        "",
 		Typ:          10,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "2006-01-02 15:04:05",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  0,
 		ResizeHeight: 0,
 		PositionX:    0,
 		PositionY:    0,
 		IsRise:       false,
 	}
 
 	data["start"] = Data{
 		Value:        "",
 		Typ:          10,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "2006-01-02 15:04:05",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  0,
 		ResizeHeight: 0,
 		PositionX:    0,
 		PositionY:    0,
 		IsRise:       false,
 	}
 
 	data["end"] = Data{
 		Value:        "",
 		Typ:          10,
 		Days:         0,
 		Months:       0,
 		Years:        5,
 		Format:       "2006-01-02 15:04:05",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  0,
 		ResizeHeight: 0,
 		PositionX:    0,
 		PositionY:    0,
 		IsRise:       false,
 	}
 
 	data["line"] = Data{
 		Value:        "",
 		Typ:          30,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  0,
 		ResizeHeight: 0,
 		PositionX:    0,
 		PositionY:    0,
 		IsRise:       false,
 	}
 	data["img1"] = Data{
 		Value:        "./src/a.png",
 		Typ:          20,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  400,
 		ResizeHeight: 400,
 		PositionX:    140,
 		PositionY:    0,
 		IsRise:       false,
 	}
 	data["img2"] = Data{
 		Value:        "./src/a.png",
 		Typ:          20,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  400,
 		ResizeHeight: 400,
 		PositionX:    600,
 		PositionY:    0,
 		IsRise:       true,
 	}
 	data["img3"] = Data{
 		Value:        "./src/a.png",
 		Typ:          21,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  80,
 		ResizeHeight: 20,
 		PositionX:    150,
 		PositionY:    95,
 		IsRise:       false,
 	}
 	data["level"] = Data{
 		Value:        "1",
 		Typ:          0,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  0,
 		ResizeHeight: 0,
 		PositionX:    0,
 		PositionY:    0,
 		IsRise:       false,
 	}
 
 	data["money"] = Data{
 		Value:        "1000",
 		Typ:          0,
 		Days:         0,
 		Months:       0,
 		Years:        0,
 		Format:       "",
 		MaxWidth:     "",
 		FontSize:     "",
 		ResizeWidth:  0,
 		ResizeHeight: 0,
 		PositionX:    0,
 		PositionY:    0,
 		IsRise:       false,
 	}
 	buf := bytes.NewBufferString("")
 	//file, _ := os.Create("./src/dst2")
 	if err:=comp.Save(buf); err != nil {
 		fmt.Println("buf error:", err.Error())
 	}
}
```
