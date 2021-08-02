package zlog

import (
	"github.com/zander-84/go-libs/think"
)

type Conf struct {
	Level       string //{debug info warn error panic fatal}
	Name        string //日志名称
	TimeZone    string
	AddCaller   bool
	ConsoleHook struct {
		Enable bool //是否启用
	}
	FileHook struct {
		Enable     bool   //是否启用
		Path       string //地址
		MaxAge     int    //日志最大的保存时间，单位天
		MaxBackups int    //最大旧文件数量
		MaxSize    int    //日志分割的尺寸，单位MB
	}
	//EmailHook struct {
	//	Enable   bool   //是否启用
	//	Level    int    //触发的级别  只在某个级别
	//	Host     string //
	//	Port     int
	//	User     string
	//	Password string
	//	To       []string
	//}
	//MysqlHook struct {
	//	Enable    bool //是否启用
	//	TableName string
	//}
	//
	//MongoHook struct {
	//	Enable    bool //是否启用
	//	TableName string
	//}
}

func (c *Conf) SetDefault() Conf {
	c.SetDefaultBasic()
	c.SetDefaultFileHook()
	return *c
}

func (c *Conf) SetDefaultBasic() {
	if c.Name == "" {
		c.Name = "log"
	}

	if c.TimeZone == "" {
		c.TimeZone = think.DefaultTimeZone
	}
}

func (c *Conf) SetDefaultFileHook() {
	if c.FileHook.Enable {
		if c.FileHook.Path == "" {
			c.FileHook.Path = "./"
		}
		if c.FileHook.MaxAge == 0 {
			c.FileHook.MaxAge = 30
		}
		if c.FileHook.MaxBackups == 0 {
			c.FileHook.MaxBackups = 30
		}
		if c.FileHook.MaxSize == 0 {
			c.FileHook.MaxSize = 100
		}
	}
}
