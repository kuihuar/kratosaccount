package data

import (
	"context"
	"nancalacc/internal/biz"
	"nancalacc/internal/data/models"
	"time"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/clause"
)

type accountRepo struct {
	data *Data
	log  *log.Helper
}

// SaveAccounts implements biz.AccountRepo.
func (r *accountRepo) SaveAccounts(ctx context.Context, accounts []*biz.Account) (int, error) {
	entities := make([]*models.Account, 0, len(accounts))
	for _, acc := range accounts {
		entities = append(entities, &models.Account{
			Username: acc.Username,
			Email:    acc.Email,
			Phone:    acc.Phone,
			Password: acc.Password,
		})
	}

	tx := r.data.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 批量Upsert（依赖ExternalID唯一约束）
	result := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},              // 冲突列
		DoUpdates: clause.AssignmentColumns([]string{"name", "email"}), // 更新字段
	}).Create(&entities)

	if result.Error != nil {
		tx.Rollback()
		return 0, errors.InternalServer("DB_SAVE_FAILED", "账户保存失败").WithCause(result.Error)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, errors.InternalServer("DB_TX_FAILED", "事务提交失败").WithCause(err)
	}

	return int(result.RowsAffected), nil
}

func NewAccountRepo(data *Data, logger log.Logger) biz.AccountRepo {
	return &accountRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accountRepo) Save(ctx context.Context, a *biz.Account) (*biz.Account, error) {
	acc := &models.Account{
		Username:  a.Username,
		Email:     a.Email,
		Phone:     a.Phone,
		Password:  a.Password,
		Status:    a.Status,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	if err := r.data.db.Create(acc).Error; err != nil {
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
