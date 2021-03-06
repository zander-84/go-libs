package gdb

import (
	"database/sql"
	"fmt"
	"github.com/zander-84/go-libs/components/errs"
	"github.com/zander-84/go-libs/components/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Gdb struct {
	engine *gorm.DB
	sqlDB  *sql.DB
	conf   Conf
	once   int64
	err    error
	lock   sync.Mutex
	time   *helper.Time
}

func (this *Gdb) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.err = errs.UninitializedError
	this.time = helper.NewTime(this.conf.TimeZone)
	this.once = 0
	this.engine = nil
	this.sqlDB = nil
}

func NewGdb(conf Conf) *Gdb {
	this := new(Gdb)
	this.init(conf)
	return this
}

func (this *Gdb) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	atomic.AddInt64(&this.once, 1)
	if atomic.LoadInt64(&this.once) != 1 {
		atomic.StoreInt64(&this.once, 2)
		return this.err
	}
	// 时间配置

	// 配置文件
	gormCnf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		AllowGlobalUpdate: false,
		NowFunc: func() time.Time {
			return time.Now().In(this.time.Location())
		},
	}
	// debug
	var LogLevel logger.LogLevel
	if this.conf.Debug {
		LogLevel = logger.Info
	} else {
		LogLevel = logger.Silent
	}
	gormCnf.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      LogLevel,    // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// mysql conf
	mysqlCnf := mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", this.conf.User, this.conf.Pwd, this.conf.Host, this.conf.Port, this.conf.Database, this.conf.Charset), // DSN data source name
		//DefaultStringSize:         256,   // string 类型字段的默认长度
		//DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
	})

	// 开始初始化
	this.engine, this.err = gorm.Open(mysqlCnf, gormCnf)
	if this.err != nil {
		return this.err
	}
	this.sqlDB, this.err = this.engine.DB()
	if this.err != nil {
		return this.err
	}

	this.sqlDB.SetMaxIdleConns(this.conf.MaxIdleconns)
	this.sqlDB.SetMaxOpenConns(this.conf.MaxOpenconns)
	this.sqlDB.SetConnMaxLifetime(time.Duration(this.conf.ConnMaxLifetime) * time.Second)

	if this.conf.RemoveSomeCallbacks {
		_ = this.engine.Callback().Create().Remove("gorm:save_before_associations")
		_ = this.engine.Callback().Create().Remove("gorm:force_reload_after_create")
		_ = this.engine.Callback().Create().Remove("gorm:save_after_associations")
		_ = this.engine.Callback().Update().Remove("gorm:save_before_associations")
		_ = this.engine.Callback().Update().Remove("gorm:save_after_associations")
	}

	return this.err
}

func (this *Gdb) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.sqlDB != nil {
		_ = this.sqlDB.Close()
	}
	this.once = 0
	this.err = errs.UninitializedError
	this.engine = nil
	this.sqlDB = nil
}

func (this *Gdb) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}

func (this *Gdb) Engine() *gorm.DB {
	return this.engine
}

func (this *Gdb) SqlDB() *sql.DB {
	return this.sqlDB
}
