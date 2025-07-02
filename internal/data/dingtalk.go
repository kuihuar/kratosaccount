package data

import (
	"context"
	"fmt"
	"nancalacc/internal/biz"

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

func (r *dingTalkRepo) FetchAccounts(ctx context.Context, token string) ([]*biz.ThirdPartyAccount, error) {
	// 调用第三方API（示例为REST）
	r.log.Info("FetchAccounts token: %s", token)
	var res []*biz.ThirdPartyAccount
	res = append(res, &biz.ThirdPartyAccount{
		RemoteID: "1",
		Name:     "张三",
		Email:    "zhangsan@example.com",
	})
	return res, nil
	// _, err := r.httpClient.Get(fmt.Sprintf("%s/accounts?token=%s", r.endpoint, token))
	// // 处理响应、JSON解析、错误转换...
	// if err != nil {
	// 	return nil, err
	// }
	// return res, nil

}

func (d *dingTalkRepo) GetAccessToken(ctx context.Context, code string) (string, error) {
	// config := &openapi.Config{
	// 	Protocol: tea.String("https"),
	// 	RegionId: tea.String("central"),
	// }

	// client, err := dingtalkoauth2_1_0.NewClient(config)
	// if err != nil {
	// 	return "", nil
	// }
	d.log.Infof("d.data.DingtalkConf: %v", d.data.thirdParty)
	request := &dingtalkoauth2_1_0.GetAccessTokenRequest{
		AppKey:    tea.String(d.data.thirdParty.AppKey),
		AppSecret: tea.String(d.data.thirdParty.AppSecret),
		//AppKey:    tea.String("dinglz1setxqhrpp7aa0"),
		// AppSecret: tea.String("uHVTlmVFojonyjlBIDbzsxLZ_iJfviqUMpT1LKNxP9P4TYr8LhaiwymiQfb0fjxr"),
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
