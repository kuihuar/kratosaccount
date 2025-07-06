package biz

import (
	"context"
	v1 "nancalacc/api/account/v1"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
// ErrUserNotFound is user not found.
// ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// AccountRepo is a Account repo.
type AccountRepo interface {
	SaveUsers(context.Context, []*DingtalkDeptUser) (int, error)
	SaveDepartments(context.Context, []*DingtalkDept) (int, error)
	SaveDepartmentUserRelations(context.Context, []*DingtalkDeptUserRelation) (int, error)
	SaveCompanyCfg(context.Context) error
}

// AccountUsecase is a Account usecase.
type AccountUsecase struct {
	repo           AccountRepo
	thirdPartyRepo DingTalkRepo
	log            *log.Helper
}

// NewAccountUsecase new a Account usecase.
func NewAccountUsecase(repo AccountRepo, dingTalkRepo DingTalkRepo, logger log.Logger) *AccountUsecase {
	return &AccountUsecase{
		repo:           repo,
		thirdPartyRepo: dingTalkRepo,
		log:            log.NewHelper(logger),
	}
}

var (
// isSyncing atomic.Bool
// taskId    uint64 = 10
)

// CreateSyncAccount creates a Account, and returns the new Account.
func (uc *AccountUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {

	// if isSyncing.CompareAndSwap(false, true) {
	// 	return nil, status.Errorf(codes.AlreadyExists, "同步任务已存在")

	// }
	// taskId := strconv.FormatUint(atomic.AddUint64(&taskId, 1), 10)
	// defer isSyncing.Store(false) // 确保锁释放

	// 0. 保存公司配置
	err := uc.repo.SaveCompanyCfg(ctx)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: err: %v", err)
	if err != nil {
		return nil, err
	}
return nil, err
	uc.log.WithContext(ctx).Infof("CreateSyncAccount: %v", req)

	// 1. 获取access_token
	accessToken, err := uc.thirdPartyRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return nil, err
	}

	// 1. 从第三方获取部门和用户数据
	depts, err := uc.thirdPartyRepo.FetchDepartments(ctx, accessToken)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: depts: %v, err: %v", depts, err)
	if err != nil {
		return nil, err
	}
	// 2. 数据入库
	deptCount, err := uc.repo.SaveDepartments(ctx, depts)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: deptCount: %v, err: %v", deptCount, err)
	if err != nil {
		return nil, err
	}
	var deptIds []int64
	for _, dept := range depts {
		deptIds = append(deptIds, dept.DeptID)
	}

	// 1. 从第三方获取用户数据
	deptUsers, err := uc.thirdPartyRepo.FetchDepartmentUsers(ctx, accessToken, deptIds)
	if err != nil {
		return nil, err
	}
	// 2. 数据入库
	userCount, err := uc.repo.SaveUsers(ctx, deptUsers)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: userCount: %v, err: %v", userCount, err)
	if err != nil {
		return nil, err
	}

	// 2. 关系数据入库
	var deptUserRelations []*DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {
		for _, depId := range deptUser.DeptIDList {
			deptUserRelations = append(deptUserRelations, &DingtalkDeptUserRelation{
				Uid:   deptUser.Userid,
				Did:   strconv.FormatInt(depId, 10),
				Order: deptUser.DeptOrder,
			})
		}

	}
	// 3. 数据入库
	relationCount, err := uc.repo.SaveDepartmentUserRelations(ctx, deptUserRelations)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: relationCount: %v, err: %v", relationCount, err)
	if err != nil {
		return nil, err
	}
	return &v1.CreateSyncAccountReply{
		TaskId:     "taskId",
		CreateTime: timestamppb.Now(),
	}, nil
}
func (s *AccountUsecase) GetSyncAccount(ctx context.Context, req *v1.GetSyncAccountRequest) (*v1.GetSyncAccountReply, error) {
	return &v1.GetSyncAccountReply{}, nil
}

func (s *AccountUsecase) CancelSyncTask(ctx context.Context, req *v1.CancelSyncAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
