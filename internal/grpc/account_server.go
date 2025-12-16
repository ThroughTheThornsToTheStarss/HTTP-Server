package grpc

import (
    "context"

    proto "git.amocrm.ru/ilnasertdinov/http-server-go/proto"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type accountDeleter interface {
    DeleteAccount(accountID uint64) error
}

type AccountServer struct {
    proto.UnimplementedAccountServiceServer
    uc accountDeleter
}

func NewAccountServer(uc accountDeleter) *AccountServer {
    return &AccountServer{uc: uc}
}

func (s *AccountServer) DisableAccount(ctx context.Context, req *proto.DisableAccountRequest) (*proto.DisableAccountResponse, error) {
    if req == nil || req.AccountId == 0 {
        return nil, status.Error(codes.InvalidArgument, "account_id is required")
    }

    if err := s.uc.DeleteAccount(req.AccountId); err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &proto.DisableAccountResponse{
        Ok:   true,
        Info: "Account disabled",
    }, nil
}

