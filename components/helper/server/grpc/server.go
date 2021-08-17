package grpc

import (
	"context"
	"github.com/zander-84/go-libs/components/helper/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

var _ server.Server = new(Server)

// Server is a gRPC server wrapper.
type Server struct {
	server *grpc.Server
	opts   serverOptions

	network string
	addr    string
}
type ServerOption func(o *serverOptions)

type serverOptions struct {
	server           *grpc.Server
}

// ServerHandler with server handler.
func ServerHandler(s *grpc.Server) ServerOption {
	return func(o *serverOptions) {
		o.server = s
	}
}


// NewServer creates a gRPC server by options.
func NewServer(network, addr string, opts ...ServerOption) *Server {
	options := serverOptions{}

	s := &Server{
		network: network,
		addr:    addr,
	}
	for _, o := range opts {
		o(&options)
	}
	s.server = options.server
	s.opts = options
	return s
}

// Start  the gRPC server.
func (s *Server) Start(ctx context.Context) error {

	lis, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	log.Printf("[gRPC] server listening on: %s \n", lis.Addr().String())

	return s.server.Serve(lis)
}

// Stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

