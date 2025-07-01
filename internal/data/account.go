package data

import (
	"context"
	"nancalacc/internal/biz"
	"nancalacc/internal/data/models"

	"github.com/go-kratos/kratos/v2/log"
)

type accountRepo struct {
	data *Data
	log  *log.Helper
}

func NewAccountRepo(data *Data, logger log.Logger) biz.AccountRepo {
	return &accountRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accountRepo) Save(ctx context.Context, a *biz.Account) (*biz.Account, error) {
	acc := &models.Account{
		Username: a.Username,
		Email:    a.Email,
		Phone:    a.Phone,
		Password: a.Password,
	}
	if err := r.data.db.Create(acc).Error; err != nil {
		return nil, err
	}
	return &biz.Account{
		ID:       acc.ID,
		Username: acc.Username,
		Email:    acc.Email,
		Phone:    acc.Phone,
	}, nil
}
func (r *accountRepo) Update(ctx context.Context, a *biz.Account) (*biz.Account, error) {
	return a, nil
}
func (r *accountRepo) FindByID(ctx context.Context, id int64) (*biz.Account, error) {
	r.log.Infof("data.FindByID: %v", id)
	var acc models.Account
	if err := r.data.db.Where("id = ?", id).First(&acc).Error; err != nil {
		return nil, err
	}
	return &biz.Account{
		ID:        acc.ID,
		Username:  acc.Username,
		Email:     acc.Email,
		Phone:     acc.Phone,
		Password:  acc.Password,
		Status:    acc.Status,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}, nil
}
func (r *accountRepo) ListAll(ctx context.Context) ([]*biz.Account, error) {
	return nil, nil
}
