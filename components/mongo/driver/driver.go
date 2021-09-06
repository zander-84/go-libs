package driver

import (
	"context"
	"fmt"
	"github.com/zander-84/go-libs/think"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"sync/atomic"
)

type Mongo struct {
	engine  *mongo.Client
	conf    Conf
	once    int64
	err     error
	lock    sync.Mutex
	context context.Context
}

func (this *Mongo) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.err = think.ErrInstanceUnDone
	atomic.StoreInt64(&this.once, 0)
	this.engine = nil
	this.context = context.Background()
}

func NewMongo(conf Conf) *Mongo {
	this := new(Mongo)
	this.init(conf)
	return this
}

func (this *Mongo) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if atomic.CompareAndSwapInt64(&this.once, 0, 1) {

		dns := fmt.Sprintf("mongodb://%s:%s", this.conf.Host, this.conf.Port)

		mongoOptions := new(options.ClientOptions)
		mongoOptions.ApplyURI(dns)

		if this.conf.User != "" && this.conf.Pwd != "" {
			mongoOptions.SetAuth(options.Credential{
				AuthMechanism:           "",
				AuthMechanismProperties: nil,
				AuthSource:              this.conf.Database,
				Username:                this.conf.User,
				Password:                this.conf.Pwd,
				PasswordSet:             false,
			})
		}
		MaxPoolSize := this.conf.MaxPoolSize
		MinPoolSize := this.conf.MinPoolSize
		mongoOptions.MaxPoolSize = &MaxPoolSize
		mongoOptions.MinPoolSize = &MinPoolSize

		this.engine, this.err = mongo.Connect(context.Background(), mongoOptions)
		if this.err != nil {
			return this.err
		}

		if this.err = this.engine.Ping(context.Background(), nil); this.err != nil {
			return this.err
		}
	}
	return this.err
}

func (this *Mongo) Stop() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.engine != nil {
		_ = this.engine.Disconnect(this.context)
	}
	this.engine = nil
	atomic.StoreInt64(&this.once, 0)
	this.err = think.ErrInstanceUnDone
	return nil
}

func (this *Mongo) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}

func (this *Mongo) Engine() *mongo.Client {
	return this.engine
}

func (this *Mongo) DB() *mongo.Database {
	return this.engine.Database(this.conf.Database)
}

func (this *Mongo) Collection(collection string) *mongo.Collection {
	return this.engine.Database(this.conf.Database).Collection(collection)
}

func (this *Mongo) GetDB(dbname string) *mongo.Database {
	return this.engine.Database(dbname)
}
func (this *Mongo) GetCollection(dbname string, collection string) *mongo.Collection {
	return this.GetDB(dbname).Collection(collection)
}
