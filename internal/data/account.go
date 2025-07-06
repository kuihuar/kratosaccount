package data

import (
	"context"
	"database/sql"
	"fmt"
	"nancalacc/internal/biz"
	"nancalacc/internal/data/models"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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
	ThirdCompanyID = "nancal"
	PlatformID     = "dingtalk"
)

// SaveUsers implements biz.AccountRepo.
func (r *accountRepo) SaveUsers(ctx context.Context, accounts []*biz.DingtalkDeptUser) (int, error) {
	entities := make([]*models.TbLasUser, 0, len(accounts))
	//taskId := time.Now().Format("20060102150405")
	var taskIds []string
	for i := 1; i <= len(accounts); i++ {
		taskId := time.Now().Add(time.Duration(i) * time.Second).Format("20060102150405")
		taskIds = append(taskIds, taskId)
	}

	for _, account := range accounts {
		r.log.Info("SaveUsersReq: account: %+v", account)
		fmt.Printf("account: %+v\n", account)
	}
	for index, account := range accounts {
		entities = append(entities, &models.TbLasUser{
			TaskID:         taskIds[index],
			ThirdCompanyID: ThirdCompanyID,
			PlatformID:     PlatformID,
			Uid:            account.Userid,
			DefDid:         sql.NullString{String: "1", Valid: true},
			DefDidOrder:    0,
			Account:        account.Userid,
			NickName:       account.Nickname,
			Email:          sql.NullString{String: account.Email, Valid: true},
			Phone:          sql.NullString{String: account.Mobile, Valid: true},
			Title:          sql.NullString{String: account.Title, Valid: true},
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:    "sync",
			Ctime:     sql.NullTime{Time: time.Now(), Valid: true},
			Mtime:     time.Now(),
			CheckType: 1,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}

	for _, entitie := range entities {
		r.log.Info("SaveUsersReq: entitie: %+v", entitie)
		fmt.Printf("entitie: %+v\n", entitie)
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

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
func (r *accountRepo) SaveDepartmentUserRelations(ctx context.Context, relations []*biz.DingtalkDeptUserRelation) (int, error) {
	entities := make([]*models.TbLasDepartmentUser, 0, len(relations))
	var taskIds []string
	for i := 1; i <= len(relations); i++ {
		taskId := time.Now().Add(time.Duration(i) * time.Second).Format("20060102150405")
		taskIds = append(taskIds, taskId)
	}
	for index, relation := range relations {
		entities = append(entities, &models.TbLasDepartmentUser{
			Did:            relation.Did,
			TaskID:         taskIds[index],
			ThirdCompanyID: ThirdCompanyID,
			PlatformID:     PlatformID,
			Uid:            relation.Uid,
			Ctime:          time.Now(),
			Order:          sql.NullInt32{Int32: int32(relation.Order), Valid: true},
			CheckType:      1,
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}
