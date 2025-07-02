package biz

import "context"

type ThirdPartyAccount struct {
	RemoteID string                 // 第三方平台ID
	Name     string                 // 账户名
	Email    string                 // 邮箱
	Phone    string                 // 手机号
	RawData  map[string]interface{} // 原始数据（扩展用）
}

type DingTalkRepo interface {
	FetchAccounts(ctx context.Context, token string) ([]*ThirdPartyAccount, error)
	GetAccessToken(ctx context.Context, code string) (string, error)
}
