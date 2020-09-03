package server

import (
	"context"
	"net"

	ctx "queryprocessor/ctx"
	"queryprocessor/handler"
	grpc_executor "queryprocessor/infuser-protobuf/gen/proto/executor"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

// Server : gRPC 서버 관련 구조체
type Server struct {
	ctx     *ctx.Context
	context context.Context
	grpc    *grpc.Server
}

// New constructor
func New(ctx *ctx.Context) *Server {
	s := new(Server)
	s.ctx = ctx

	s.grpc = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	return s
}

// Run : gRPC 서버 시작 기능
func (s *Server) Run(network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	apiResultHandler := handler.NewApiResultHandler(s.ctx)
	grpc_executor.RegisterApiResultServiceServer(s.grpc, newApiResultServer(apiResultHandler))

	// go func() {
	// 	defer s.grpc.GracefulStop()
	// 	<-s.context.Done()
	// }()

	println("Server is Running")
	return s.grpc.Serve(l)
}
