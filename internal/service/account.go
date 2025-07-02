package service

import (
	"context"
	"time"

	//pb "nancalacc/api/account/v1"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type AccountService struct {
	v1.UnimplementedAccountServer
	accUsecase *biz.AccountUsecase
	logger     *log.Helper
}

func NewAccountService(accUsecase *biz.AccountUsecase, logger log.Logger) *AccountService {
	return &AccountService{
		accUsecase: accUsecase,
		logger:     log.NewHelper(logger),
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *v1.CreateAccountRequest) (*v1.CreateAccountReply, error) {
	_, err := s.accUsecase.CreateAccount(ctx, &biz.Account{
		Username: req.Username,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: "123456",
	})

	if err != nil {
		return nil, err
	}

	return &v1.CreateAccountReply{}, nil
	//return &pb.CreateAccountReply{}, nil
}
func (s *AccountService) UpdateAccount(ctx context.Context, req *v1.UpdateAccountRequest) (*v1.UpdateAccountReply, error) {
	s.logger.WithContext(ctx).Infof("req: %v", req)
	num, err := s.accUsecase.SyncAccount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateAccountReply{
		Num: int32(num),
	}, nil
}
func (s *AccountService) DeleteAccount(ctx context.Context, req *v1.DeleteAccountRequest) (*v1.DeleteAccountReply, error) {
	return &v1.DeleteAccountReply{}, nil
}
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {

	s.logger.WithContext(ctx).Infof("req: %v", req)

	acc, err := s.accUsecase.GetAccountByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetAccountReply{
		User: &v1.UserInfo{
			Id:        acc.ID,
			Username:  acc.Username,
			Phone:     acc.Phone,
			Email:     acc.Email,
			Status:    acc.Status,
			CreatedAt: time.Unix(acc.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt: time.Unix(acc.UpdatedAt, 0).Format(time.RFC3339),
		},
	}, nil
}
func (s *AccountService) ListAccount(ctx context.Context, req *v1.ListAccountRequest) (*v1.ListAccountReply, error) {
	return &v1.ListAccountReply{}, nil
}
