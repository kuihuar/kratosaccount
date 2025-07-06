package data

import (
	"context"
	"encoding/json"
	"fmt"
	"nancalacc/internal/biz"
	"nancalacc/pkg/httputil"
	"sync"

	"time"

	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-kratos/kratos/v2/log"
)

type dingTalkRepo struct {
	data *Data
	log  *log.Helper
}

func NewDingTalkRepo(data *Data, logger log.Logger) biz.DingTalkRepo {

	return &dingTalkRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/dingtalk")),
	}
}

func (r *dingTalkRepo) FetchDepartments(ctx context.Context, token string) ([]*biz.DingtalkDept, error) {
	var deptList []*biz.DingtalkDept

	// 1. 获取子部门ID列表（所有）
	deptIdlist, err := r.getDeptIds(ctx, token)
	if err != nil {
		return nil, err
	}
	r.log.Info("FetchAccounts.deptIdlist: %v", deptIdlist)
	// 2. 获取子部门详情
	deptList, err = r.fetchDeptDetails(ctx, token, deptIdlist, 10)
	if err != nil {
		return nil, err
	}
	return deptList, nil
}

func (r *dingTalkRepo) getDeptIds(ctx context.Context, token string) ([]int64, error) {
	uri := fmt.Sprintf("%s/topapi/v2/department/listsubid?access_token=%s", "https://oapi.dingtalk.com", token)
	input := &biz.ListDeptIDRequest{
		DeptID: 1,
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
	if err != nil {
		return nil, err
	}

	r.log.Info("FetchAccounts.deptList: %v, err: %v", string(bs), err)

	var deptIDResponse *biz.ListDeptIDResponse
	if err = json.Unmarshal(bs, &deptIDResponse); err != nil {
		return nil, err
	}
	if deptIDResponse.Errcode != 0 {
		return nil, fmt.Errorf("钉钉API返回错误: %s, errcode: %d", deptIDResponse.Errmsg, deptIDResponse.Errcode)
	}
	if deptIDResponse.Result.DeptIDList == nil {
		return nil, fmt.Errorf("钉钉API返回错误: %s, errcode: %d", deptIDResponse.Errmsg, deptIDResponse.Errcode)
	}
	deptIdlist := deptIDResponse.Result.DeptIDList
	return deptIdlist, nil
}
func (r *dingTalkRepo) fetchDeptDetails(ctx context.Context, token string, deptIds []int64, maxConcurrent int) ([]*biz.DingtalkDept, error) {
	uriDetail := fmt.Sprintf("%s/topapi/v2/department/get?access_token=%s", "https://oapi.dingtalk.com", token)
	sem := make(chan struct{}, maxConcurrent)
	results := make(chan *biz.DingtalkDept, len(deptIds))
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	for _, deptId := range deptIds {
		wg.Add(1)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			continue
		}

		// 启动goroutine处理任务
		go func(id int64) {
			defer func() {
				<-sem     // 释放信号量
				wg.Done() // 通知任务完成
			}()

			input := &biz.DingtalkDeptRequest{
				DeptID: id,
			}
			jsonData, err := json.Marshal(input)
			if err != nil {
				r.log.Info("fetchDeptDetails.jsonData: %v, err: %v", string(jsonData), err)
				errChan <- err
				return
			}

			bs, err := httputil.PostJSON(uriDetail, jsonData, time.Second*10)
			if err != nil {
				r.log.Info("fetchDeptDetails.PostJSON: %v, err: %v", string(jsonData), err)
				errChan <- err
				return
			}
			var deptResponse *biz.DingtalkDeptResponse
			if err = json.Unmarshal(bs, &deptResponse); err != nil {
				r.log.Info("fetchDeptDetails.Unmarshal: %v, err: %v", string(bs), err)
				errChan <- err
				return
			}
			if deptResponse.Errcode != 0 {
				r.log.Info("fetchDeptDetails.Errcode: %v, err: %v", deptResponse.Errcode, deptResponse.Errmsg)
				errChan <- err
				return
			}
			results <- &deptResponse.Result
		}(deptId) // 传递当前deptId值
	}
	wg.Wait()

	close(results)
	close(errChan)
	var deptList []*biz.DingtalkDept
	for dept := range results {
		deptList = append(deptList, dept)
	}

	return deptList, nil

}

func (r *dingTalkRepo) FetchDepartmentUsers(ctx context.Context, token string, deptIds []int64) ([]*biz.DingtalkDeptUser, error) {
	// 服务端API.通讯录管理.用户管理.获取部门用户详情
	maxConcurrent := 10
	sem := make(chan struct{}, maxConcurrent)
	results := make(chan *biz.DingtalkDeptUser, len(deptIds))
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	for _, deptId := range deptIds {
		wg.Add(1)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			continue
		}

		// 启动goroutine处理任务
		go func(id int64) {
			defer func() {
				<-sem     // 释放信号量
				wg.Done() // 通知任务完成
			}()

			for {
				userList, cursor, err := r.getUserListByDepId(ctx, token, id)
				if err != nil {
					r.log.Info("FetchDepartmentUsers.getUserListByDepId: %v, err: %v", id, err)
					errChan <- err
					return
				}
				for _, user := range userList {
					results <- user
				}
				if cursor == 0 {
					break
				}

			}
		}(deptId)
	}
	wg.Wait()

	close(results)
	close(errChan)
	var userList []*biz.DingtalkDeptUser
	for user := range results {

		userList = append(userList, user)

	}
	for _, user := range userList {
		r.log.Info("FetchDepartmentUsers.userList.user: %v", user)
	}
	return userList, nil
}
func (r *dingTalkRepo) getUserListByDepId(ctx context.Context, token string, deptId int64) ([]*biz.DingtalkDeptUser, int64, error) {
	// 发送post请求
	var cursor int64 = 0
	uri := fmt.Sprintf("%s/topapi/v2/user/list?access_token=%s", "https://oapi.dingtalk.com", token)
	input := &biz.ListDeptUserRequest{
		DeptID: deptId,
		Cursor: cursor,
		Size:   100,
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, 0, err
	}

	r.log.Info("getUserListByDepId.uri: %v, input: %v, jsonData: %v", uri, input, string(jsonData))

	bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
	r.log.Info("getUserListByDepId.body: %v, err: %v", string(bs), err)
	if err != nil {
		return nil, 0, err
	}

	// 打印响应体
	fmt.Println(string(bs))

	var userResponse biz.ListDeptUserResponse
	if err = json.Unmarshal(bs, &userResponse); err != nil {
		return nil, 0, err
	}
	if userResponse.Errcode != 0 {
		return nil, 0, fmt.Errorf("钉钉API返回错误: %s, errcode: %v", userResponse.Errmsg, userResponse.Errcode)
	}

	var userList []*biz.DingtalkDeptUser
	if userResponse.Result.List != nil {
		userList = make([]*biz.DingtalkDeptUser, 0, len(userResponse.Result.List))
		for _, user := range userResponse.Result.List {
			userList = append(userList, &user)
		}
	}
	if userResponse.Result.HasMore {
		return userList, userResponse.Result.NextCursor, nil
	}
	return userList, 0, nil
}
func (d *dingTalkRepo) GetAccessToken(ctx context.Context, code string) (string, error) {
	d.log.Infof("d.data.DingtalkConf: %v", d.data.thirdParty)
	request := &dingtalkoauth2_1_0.GetAccessTokenRequest{
		AppKey:    tea.String(d.data.thirdParty.AppKey),
		AppSecret: tea.String(d.data.thirdParty.AppSecret),
	}

	var accessToken string

	tryErr := func() error {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err := r
				fmt.Printf("恢复的错误: %v\n", err)
			}
		}()

		response, err := d.data.dingtalkCli.GetAccessToken(request)
		if err != nil {
			return err
		}

		accessToken = *response.Body.AccessToken
		return nil
	}()

	if tryErr != nil {
		// 处理错误
		var sdkErr = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			sdkErr = _t
		} else {
			sdkErr.Message = tea.String(tryErr.Error())
		}

		if !tea.BoolValue(util.Empty(sdkErr.Code)) && !tea.BoolValue(util.Empty(sdkErr.Message)) {
			return "", fmt.Errorf("获取access_token失败: [%s] %s", *sdkErr.Code, *sdkErr.Message)
		}
		return "", fmt.Errorf("获取access_token失败: %s", *sdkErr.Message)
	}

	return accessToken, nil
}
