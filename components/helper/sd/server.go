package sd

import (
	"context"
	"time"
)

type Server struct {
	key        string
	val        string
	ttl        int64
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func RegisterServer(key string, val string, ttl int64) *Server {
	s := new(Server)
	s.key = key
	s.val = val
	s.ttl = ttl
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	return s
}

func (s *Server) Deregister() {
	s.cancelFunc()
}

func (s *Server) Key() string {
	return s.key
}

func (s *Server) Val() string {
	return s.val
}

func (s *Server) TTlSecond() time.Duration {
	return time.Second * time.Duration(s.ttl)
}

func (s *Server) TTl() int64 {
	return s.ttl
}

func (s *Server) Context() context.Context {
	return s.ctx
}
