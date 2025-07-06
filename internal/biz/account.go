package biz

import (
	"context"
	v1 "nancalacc/api/account/v1"
	"time"

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

	taskId := time.Now().Format("20060102150405")
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
	return &v1.CreateSyncAccountReply{
		TaskId:     taskId,
		CreateTime: timestamppb.Now(),
	}, nil

	// 1. 从第三方获取用户数据

	deptUsers, _, err := uc.thirdPartyRepo.FetchDepartmentUsers(ctx, accessToken, 1, 0)
	if err != nil {
		return nil, err
	}
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: deptUsers: %v, err: %v", deptUsers, err)
	if err != nil {
		return nil, err
		//return errors.New("THIRD_PARTY_FAIL第三方API调用失败")
	}

	// 2. 数据入库
	userCount, err := uc.repo.SaveUsers(ctx, deptUsers)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: userCount: %v, err: %v", userCount, err)
	if err != nil {
		return nil, err
	}
	var deptUserRelations []*DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {
		deptUserRelations = append(deptUserRelations, &DingtalkDeptUserRelation{
			TaskID:         "taskId",
			ThirdCompanyID: "1",
			PlatformID:     "dingtalk",
			Uid:            deptUser.Userid,
			//Did:            strconv.FormatInt(deptUser.DeptID, 10),
			//Order:          sql.NullInt32{Int32: 0, Valid: true},
			Main: 0,
		})
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
