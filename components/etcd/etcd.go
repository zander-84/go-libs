package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/zander-84/go-libs/think"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"sync/atomic"
	"time"
)

type Etcd struct {
	engine  *clientv3.Client
	conf    Conf
	once    int64
	err     error
	lock    sync.Mutex
	context context.Context
}

func (this *Etcd) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.err = think.ErrInstanceUnDone
	this.context = context.Background()
	this.engine = nil
	atomic.StoreInt64(&this.once, 0)
}
func NewEtcd(conf Conf) *Etcd {
	this := new(Etcd)
	this.init(conf)
	return this
}

func (this *Etcd) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if atomic.CompareAndSwapInt64(&this.once, 0, 1) {
		var tlsconf *tls.Config
		if this.conf.TlsCa != "" && this.conf.TlsKey != "" && this.conf.TlsPem != "" {
			ce, err := tls.X509KeyPair([]byte(this.conf.TlsPem), []byte(this.conf.TlsKey))
			if err != nil {
				this.err = err
				return this.err
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM([]byte(this.conf.TlsCa))
			tlsconf = &tls.Config{
				Certificates: []tls.Certificate{ce},
				RootCAs:      pool,
			}
		}
		this.engine, this.err = clientv3.New(clientv3.Config{
			Endpoints:            this.conf.Endpoints,
			AutoSyncInterval:     0,
			DialTimeout:          20 * time.Second,
			DialKeepAliveTime:    0,
			DialKeepAliveTimeout: 0,
			MaxCallSendMsgSize:   0,
			MaxCallRecvMsgSize:   0,
			TLS:                  tlsconf,
			Username:             "",
			Password:             "",
			RejectOldCluster:     false,
			DialOptions:          nil,
			Context:              nil,
			Logger:               nil,
			LogConfig:            nil,
			PermitWithoutStream:  false,
		})
	}
	return this.err
}

func (this *Etcd) Stop() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.engine != nil {
		_ = this.engine.Close()
	}
	this.engine = nil
	atomic.StoreInt64(&this.once, 0)
	this.err = think.ErrInstanceUnDone
	return nil
}

func (this *Etcd) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}
func (this *Etcd) Engine() *clientv3.Client {
	return this.engine
}
