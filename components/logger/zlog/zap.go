package zlog

import (
	"github.com/zander-84/go-libs/components/helper"
	"github.com/zander-84/go-libs/components/logger"
	"github.com/zander-84/go-libs/think"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var _ logger.Logger = (*ZLog)(nil)

type ZLog struct {
	engine  *zap.Logger
	conf    Conf
	err     error
	lock    sync.Mutex
	once    int64
	time    *helper.Time
	file    *helper.File
	writers []io.Writer
}

func NewZapLog(conf Conf, writers []io.Writer) *ZLog {
	this := new(ZLog)
	this.init(conf, writers)
	return this
}

func (this *ZLog) init(conf Conf, writers []io.Writer) {
	this.conf = conf.SetDefault()
	this.time = helper.NewTime()
	this.file = helper.NewFile()
	this.err = think.ErrInstanceUnDone
	atomic.StoreInt64(&this.once, 0)
	this.writers = writers
}
func (this *ZLog) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if atomic.CompareAndSwapInt64(&this.once, 0, 1) {
		newCore := make([]zapcore.Core, 0)
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = func(i time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(this.time.FormatHyphen())
		}

		logLevel := zap.DebugLevel
		switch this.conf.Level {
		case "debug":
			logLevel = zap.DebugLevel
		case "info":
			logLevel = zap.InfoLevel
		case "warn":
			logLevel = zap.WarnLevel
		case "error":
			logLevel = zap.ErrorLevel
		case "panic":
			logLevel = zap.PanicLevel
		case "fatal":
			logLevel = zap.FatalLevel
		default:
			logLevel = zap.InfoLevel
		}

		priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= logLevel
		})

		//____ 控制台输出
		if this.conf.ConsoleHook.Enable {
			console := zapcore.Lock(os.Stdout)
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
			newCore = append(newCore,
				zapcore.NewCore(consoleEncoder, console, priority),
			)
		}

		//____ 文件写入
		if this.conf.FileHook.Enable {
			this.err = this.file.OpenOrCreateWithAction(this.conf.FileHook.Path, this.conf.Name+".log", func(f *os.File) {
				log.SetOutput(f)
				log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
				log.Println("starting log...")
			})
			if this.err != nil {
				return this.err
			}

			fileHook := lumberjack.Logger{
				Filename:   this.file.GetPrefixPath(this.conf.FileHook.Path, this.conf.Name+".log"), // 日志文件路径
				MaxSize:    this.conf.FileHook.MaxSize,                                              // 每个日志文件保存的最大尺寸 单位：M
				MaxBackups: this.conf.FileHook.MaxBackups,                                           // 日志文件最多保存多少个备份
				MaxAge:     this.conf.FileHook.MaxAge,                                               // 文件最多保存多少天
				Compress:   false,                                                                   // 是否压缩
			}
			fileWriter := zapcore.AddSync(&fileHook)
			encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
			jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
			newCore = append(newCore,
				zapcore.NewCore(jsonEncoder, fileWriter, priority),
			)
		}

		if len(this.writers) > 0 {
			for _, writer := range this.writers {
				encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
				newCore = append(newCore, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), priority))
			}
		}

		if this.conf.AddCaller {
			this.engine = zap.New(zapcore.NewTee(newCore...), zap.AddCaller())
		} else {
			this.engine = zap.New(zapcore.NewTee(newCore...))
		}

		this.err = nil
	}

	return this.err
}

func (this *ZLog) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.engine = nil
	atomic.StoreInt64(&this.once, 0)
	this.err = think.ErrInstanceUnDone
}

func (this *ZLog) Restart(conf Conf, writers []io.Writer) error {
	this.Stop()
	this.init(conf, writers)
	return this.Start()
}

func (this *ZLog) Engine() *zap.Logger {
	return this.engine
}
