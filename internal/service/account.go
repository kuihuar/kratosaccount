package service

import (
	"context"

	accountV1 "nancalacc/api/account/v1"

	"nancalacc/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/go-kratos/kratos/v2/log"
)

type AccountService struct {
	accountV1.UnimplementedAccountServer
	accUsecase *biz.AccountUsecase
	logger     *log.Helper
}

func NewAccountService(accUsecase *biz.AccountUsecase, logger log.Logger) *AccountService {
	return &AccountService{
		accUsecase: accUsecase,
		logger:     log.NewHelper(logger),
	}
}
func (s *AccountService) CreateSyncAccount(ctx context.Context, req *accountV1.CreateSyncAccountRequest) (*accountV1.CreateSyncAccountReply, error) {
	s.logger.Infof("CreateSyncAccount req: %v", req)
	_, err := s.accUsecase.CreateSyncAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return &accountV1.CreateSyncAccountReply{
		TaskId:     "10",
		CreateTime: timestamppb.Now(),
	}, nil
}
func (s *AccountService) GetSyncAccount(ctx context.Context, req *accountV1.GetSyncAccountRequest) (*accountV1.GetSyncAccountReply, error) {
	return &accountV1.GetSyncAccountReply{}, nil
}
func (s *AccountService) CancelSyncTask(ctx context.Context, req *accountV1.CancelSyncAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
