package goredis

import (
	"context"
	"fmt"
	"github.com/zander-84/go-libs/components/helper"
	"github.com/zander-84/go-libs/components/helper/sd"
	"github.com/zander-84/go-libs/think"
	"time"
)

type Typ int32

var (
	UseNode  Typ = 0 // 主节点+从节点
	UseSlave Typ = 1 // 从节点
)

// RedisMasterSalve redis 主从
type RedisMasterSalve struct {
	master        *Rdb
	salves        []*Rdb
	nodes         []*Rdb
	typ           Typ
	nodeBalancer  sd.Balancer
	salveBalancer sd.Balancer
	debug         bool
}

func NewRedisMasterSalve(master *Rdb, salves []*Rdb, typ Typ) *RedisMasterSalve {
	d := new(RedisMasterSalve)
	d.typ = typ
	d.master = master
	d.salves = salves

	d.nodes = make([]*Rdb, 0)
	d.nodes = append(d.nodes, master)
	d.nodes = append(d.nodes, salves...)

	nodeListener := sd.NewListener("nodes")
	nodesTmp := map[string]int{}
	for i := 0; i < len(d.nodes); i++ {
		weight := 2
		if i == 0 {
			weight = 1
		}
		nodesTmp[fmt.Sprintf("%d", i)] = weight
	}
	nodeListener.Set(nodesTmp)
	d.nodeBalancer = sd.NewBalancer(nodeListener, false, sd.WeightRoundRobin)

	salvesTmp := map[string]int{}
	for i := 0; i < len(salves); i++ {
		salvesTmp[fmt.Sprintf("%d", i)] = 1
	}
	salveListener := sd.NewListener("slaves")
	salveListener.Set(salvesTmp)
	d.salveBalancer = sd.NewBalancer(salveListener, false, sd.RoundRobin)
	return d
}

func (this *RedisMasterSalve) SetDebug() {
	this.debug = true
}

func (this *RedisMasterSalve) GetNode() *Rdb {
	key, _ := this.nodeBalancer.Next()
	if this.debug {
		fmt.Println("form node :", key)
	}
	return this.nodes[helper.GetConv().ShouldStoI(key)]
}

func (this *RedisMasterSalve) GetSalve() *Rdb {
	key, _ := this.salveBalancer.Next()
	if this.debug {
		fmt.Println("form salve :", key)
	}
	return this.salves[helper.GetConv().ShouldStoI(key)]
}

func (this *RedisMasterSalve) getRdb() *Rdb {
	if this.typ == UseNode {
		return this.GetNode()
	}
	return this.GetSalve()
}

func (this *RedisMasterSalve) Get(ctx context.Context, key string, toPtr interface{}) error {
	return this.getRdb().Get(ctx, key, toPtr)
}

func (this *RedisMasterSalve) GetFromMaster(ctx context.Context, key string, recPtr interface{}) error {
	return this.master.Get(ctx, key, recPtr)
}

func (this *RedisMasterSalve) GetFromSlave(ctx context.Context, key string, recPtr interface{}) error {
	return this.GetSalve().Get(ctx, key, recPtr)
}

func (this *RedisMasterSalve) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return this.master.Set(ctx, key, value, ttl)
}

func (this *RedisMasterSalve) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return this.master.SetNX(ctx, key, value, ttl)
}

func (this *RedisMasterSalve) GetOrSet(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, f func() (interface{}, error)) error {
	err := this.getRdb().Get(ctx, key, ptrValue)
	if err == nil {
		return err
	}

	if err != nil && err != think.ErrInstanceRecordNotFound {
		return err
	}

	return this.master.GetOrSet(ctx, key, ptrValue, ttl, f)
}

func (this *RedisMasterSalve) Delete(ctx context.Context, key ...string) error {
	return this.master.Delete(ctx, key...)
}

func (this *RedisMasterSalve) Exists(ctx context.Context, key ...string) (bool, error) {
	return this.getRdb().Exists(ctx, key...)
}

func (this *RedisMasterSalve) FlushDB(ctx context.Context) error {
	return this.master.FlushDB(ctx)
}
