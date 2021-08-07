package etcd

import (
	"context"
	"github.com/zander-84/go-libs/components/helper/sd"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

func (this *Etcd) RegisterServer(s *Server) error {
	go func() {
		for {
			select {
			case <-s.Context().Done():
				return
			case <-time.After(s.TTlSecond()):
				if getResp, err := this.get(s.Key()); err == nil && getResp.Count == 0 {
					_ = this.withAlive(s.Context(), s.Key(), s.Val(), s.TTl())
				}
			}
		}
	}()
	return nil
}

func (this *Etcd) Deregister(s *Server) error {
	s.Deregister()
	_, err := this.delete(s.Key(), clientv3.WithIgnoreLease())
	return err
}

func (this *Etcd) get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return this.engine.Get(ctx, key, opts...)
}
func (this *Etcd) delete(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return this.engine.Delete(ctx, key, opts...)
}
func (this *Etcd) withAlive(ctx context.Context, key string, val string, ttl int64) error {
	leaseResp, err := this.engine.Grant(ctx, ttl)
	if err != nil {
		return err
	}
	_, err = this.engine.Put(ctx, key, val, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	ch, err := this.engine.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			r, ok := <-ch
			if ok {
				if r == nil {
					return
				}
			} else {
				return
			}
		}
	}()
	return nil
}

func (this *Etcd) Watch(listener *sd.Listener) {
	getResp, err := this.engine.Get(listener.Context(), listener.Name(),
		clientv3.WithPrefix())
	if err != nil {

	} else {
		for i := range getResp.Kvs {
			rowAddr := strings.TrimPrefix(string(getResp.Kvs[i].Key), listener.Name())
			addr, weight := listener.SpiltByHyphen(rowAddr)
			listener.AddWeight(addr, weight)
		}
	}

	go func() {
		rch := this.engine.Watch(listener.Context(), listener.Name(), clientv3.WithPrefix())
		for n := range rch {
			for _, ev := range n.Events {
				rowAddr := strings.TrimPrefix(string(ev.Kv.Key), listener.Name())
				addr, weight := listener.SpiltByHyphen(rowAddr)

				switch ev.Type {
				case mvccpb.PUT:
					if !listener.Exist(addr) {
						listener.AddWeight(addr, weight)
					}
				case mvccpb.DELETE:
					listener.Remove(addr)
				}
			}
		}
	}()
}
