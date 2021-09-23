package etcd

import (
	"context"
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

// withAlive 会产生warning ctx cancel
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

func (this *Etcd) GetEntries(listener Listener) (map[string]int, error) {
	if data, err := this.getEntries(listener, 0); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (this *Etcd) getEntries(listener Listener, tag int) (map[string]int, error) {
	resp, err := this.engine.Get(this.context, listener.Name(), clientv3.WithPrefix())
	if err != nil {
		tag += 1
		if tag > 5 {
			return nil, err
		} else {
			time.Sleep(time.Millisecond)
			return this.getEntries(listener, tag)
		}
	}

	entries := make(map[string]int)
	for _, kv := range resp.Kvs {
		rowAddr := strings.TrimPrefix(string(kv.Key), listener.Name())
		addr, weight := listener.GetAddrWithWeight(rowAddr)
		entries[addr] = weight
	}
	return entries, err
}

func (this *Etcd) Watch(listener Listener) error {
	if entries, err := this.GetEntries(listener); err == nil {
		if err := listener.Set(entries); err != nil {
			listener.RecordErr(err)
			return err
		}
		listener.CleanErr()
	} else {
		listener.RecordErr(err)
		return err
	}

	go func() {
		rch := this.engine.Watch(listener.Context(), listener.Name(), clientv3.WithPrefix())
		for _ = range rch {
			if entries, err := this.GetEntries(listener); err == nil {
				if err := listener.Set(entries); err != nil {
					listener.RecordErr(err)
				}
				listener.CleanErr()
			} else {
				listener.RecordErr(err)
			}
			//for _, ev := range n.Events {
			//	rowAddr := strings.TrimPrefix(string(ev.Kv.Key), listener.Name())
			//	addr, weight := listener.GetAddrWithWeight(rowAddr)
			//
			//	switch ev.Type {
			//	case mvccpb.PUT:
			//		if !listener.Exist(addr) {
			//			_ = listener.AddWeight(addr, weight)
			//		}
			//	case mvccpb.DELETE:
			//		_ = listener.Remove(addr)
			//	}
			//}
		}
	}()
	return nil
}

type Listener interface {
	Name() string
	GetAddrWithWeight(row string) (string, int)
	Context() context.Context
	AddWeight(addr string, weight int) error
	Exist(addr string) bool
	Remove(addr string) error
	Set(data map[string]int) error
	RecordErr(err error)
	CleanErr()
}
