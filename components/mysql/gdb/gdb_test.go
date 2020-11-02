package gdb

import (
	"flag"
	"gorm.io/gorm"
	"testing"
)

// go test -v  -run TestGdb  -args 172.16.86.150  3307 zander zander test2
func TestGdb(t *testing.T) {

	if !flag.Parsed() {
		flag.Parse()
	}
	argList := flag.Args()

	gdb := NewGdb(Conf{
		Host:                argList[0],
		Port:                argList[1],
		User:                argList[2],
		Pwd:                 argList[3],
		Database:            argList[4],
		Charset:             "utf8mb4",
		MaxIdleconns:        100,
		MaxOpenconns:        1000,
		ConnMaxLifetime:     300,
		Debug:               true,
		TimeZone:            "",
		RemoveSomeCallbacks: true,
	})

	if err := gdb.Start(); err != nil {
		t.Fatal(err.Error())
	}
	defer gdb.Stop()
	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}

	if err := gdb.Engine().AutoMigrate(&Product{}); err != nil {
		t.Fatal(err.Error())
	}

	if err := gdb.Engine().Create(&Product{Code: "D42", Price: 100}).Error; err != nil {
		t.Fatal(err.Error())
	}

	var product Product
	gdb.Engine().First(&product, 1) // 根据整形主键查找
	t.Log(product)
	gdb.Engine().First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	t.Log(product)

	if err := gdb.Engine().Model(&product).Update("Price", 200).Error; err != nil {
		t.Fatal(err.Error())
	}

	if err := gdb.Engine().Model(&product).Updates(Product{Price: 200, Code: "F42"}).Error; err != nil {
		t.Fatal(err.Error())
	}

	if err := gdb.Engine().Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"}).Error; err != nil {
		t.Fatal(err.Error())
	}

	if err := gdb.Engine().Delete(&product, 1).Error; err != nil {
		t.Fatal(err.Error())
	}

	t.Log("success")
}
