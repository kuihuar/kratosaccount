package data

import (
	"context"
	"database/sql"
	"nancalacc/internal/biz"
	"nancalacc/internal/data/models"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/clause"
)

type accountRepo struct {
	data *Data
	log  *log.Helper
}

func NewAccountRepo(data *Data, logger log.Logger) *accountRepo {
	return &accountRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

var (
	ThirdCompanyID = "10"
	PlatformID     = "1010"
)

// SaveUsers implements biz.AccountRepo.
func (r *accountRepo) SaveUsers(ctx context.Context, accounts []*biz.DingtalkDeptUser) (int, error) {
	entities := make([]*models.TbLasUser, 0, len(accounts))
	for _, acc := range accounts {
		entities = append(entities, &models.TbLasUser{
			Uid:      acc.Userid,
			NickName: acc.Name,
			Email:    sql.NullString{String: acc.Email, Valid: acc.Email != ""},
			Phone:    sql.NullString{String: acc.Mobile, Valid: acc.Mobile != ""},
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

// ID             uint           `gorm:"primaryKey;autoIncrement;column:id;type:int unsigned;comment:主键id" json:"id"`
//
//	Did            string         `gorm:"not null;column:did;type:varchar(255);comment:部门id" json:"did"`
//	TaskID         string         `gorm:"not null;column:task_id;type:varchar(20);comment:任务id" json:"task_id"`
//	ThirdCompanyID string         `gorm:"not null;column:third_company_id;type:varchar(20);comment:租户id" json:"third_company_id"`
//	PlatformID     string         `gorm:"not null;column:platform_id;type:varchar(60);comment:平台id" json:"platform_id"`
//	Pid            sql.NullString `gorm:"column:pid;type:varchar(255);comment:父部门id" json:"pid"`
//	Name           string         `gorm:"not null;column:name;type:varchar(255);comment:部门名称" json:"name"`
//	Order          int            `gorm:"column:order;type:int;default:0;comment:排序" json:"order"`
//	Source         string         `gorm:"column:source;type:varchar(20);default:sync;comment:来源" json:"source"`
//	Ctime          sql.NullTime   `gorm:"column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"ctime"`
//	Mtime          time.Time      `gorm:"not null;column:mtime;type:timestamp;comment:修改时间" json:"mtime"`
//	CheckType      int8           `gorm:"not null;column:check_type;type:tinyint;default:0;comment:1-勾选 0-未勾选" json:"check_type"`
//	Type           sql.NullString `gorm:"column:type;type:varchar(255);comment:类型" json:"type"`
//
// SaveDepartments implements biz.AccountRepo.
func (r *accountRepo) SaveDepartments(ctx context.Context, departments []*biz.DingtalkDept) (int, error) {
	entities := make([]*models.TbLasDepartment, 0, len(departments))
	//taskId := time.Now().Format("20060102150405")
	var taskIds []string
	for i := 1; i <= len(departments); i++ {
		taskId := time.Now().Add(time.Duration(i) * time.Second).Format("20060102150405")
		taskIds = append(taskIds, taskId)
	}
	for index, dep := range departments {
		entities = append(entities, &models.TbLasDepartment{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			TaskID:         taskIds[index],
			Name:           dep.Name,
			ThirdCompanyID: ThirdCompanyID,
			PlatformID:     PlatformID,
			Pid:            sql.NullString{String: strconv.FormatInt(dep.ParentID, 10), Valid: true},
			Order:          int(dep.Order),
			Source:         "sync",
			Ctime:          sql.NullTime{Time: time.Now(), Valid: true},
			Mtime:          time.Now(),
			CheckType:      1,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// SaveDepartmentUserRelations implements biz.AccountRepo.
func (r *accountRepo) SaveDepartmentUserRelations(ctx context.Context, departmentUsers []*biz.DingtalkDeptUserRelation) (int, error) {
	return 0, nil
}
