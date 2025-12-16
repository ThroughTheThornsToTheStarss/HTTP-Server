package grpc

import (
	"net"

	proto "git.amocrm.ru/ilnasertdinov/http-server-go/proto"
	"google.golang.org/grpc"
)

type Handler struct {
	srv *grpc.Server
	lis net.Listener
}

func NewHandler(addr string, accountSrv proto.AccountServiceServer) (*Handler, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	proto.RegisterAccountServiceServer(s, accountSrv)

	return &Handler{srv: s, lis: lis}, nil
}

func (h *Handler) Run() error { return h.srv.Serve(h.lis) }

func (h *Handler) Stop() {
	h.srv.GracefulStop()
	_ = h.lis.Close()
}
