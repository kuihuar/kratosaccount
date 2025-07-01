package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

var (
// ErrUserNotFound is user not found.
// ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Account is a Account model.
type Account struct {
	ID        int64
	Username  string
	Email     string
	Phone     string
	Password  string
	Status    int32
	CreatedAt int64
	UpdatedAt int64
}

// AccountRepo is a Account repo.
type AccountRepo interface {
	Save(context.Context, *Account) (*Account, error)
	Update(context.Context, *Account) (*Account, error)
	FindByID(context.Context, int64) (*Account, error)
	//ListByPage(context.Context, page, page_size int32) ([]*Account, error)
	ListAll(context.Context) ([]*Account, error)
}

// AccountUsecase is a Account usecase.
type AccountUsecase struct {
	repo AccountRepo
	log  *log.Helper
}

// NewAccountUsecase new a Account usecase.
func NewAccountUsecase(repo AccountRepo, logger log.Logger) *AccountUsecase {
	return &AccountUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateAccount creates a Account, and returns the new Account.
func (uc *AccountUsecase) CreateAccount(ctx context.Context, acc *Account) (*Account, error) {
	// TODO: add biz logic
	// 比如密码加密
	uc.log.WithContext(ctx).Infof("CreateAccount: %v", acc.ID)
	return uc.repo.Save(ctx, acc)
}

func (uc *AccountUsecase) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	// TODO: add biz logic

	uc.log.WithContext(ctx).Infof("biz.GetAccountByID: %v", id)
	return uc.repo.FindByID(ctx, id)
}
