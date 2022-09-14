package message

import (
	"fmt"

	"github.com/lenye/pmsg/pkg/http/client"
	"github.com/lenye/pmsg/pkg/weixin"
)

/*
url 和 miniprogram 都是非必填字段，若都不传则模板无跳转；
若都传，会优先跳转至小程序。
开发者可根据实际需要选择其中一种跳转方式即可。
当用户的微信客户端版本不支持跳小程序时，将会跳转至url。
*/

// Template 模板消息
type Template struct {
	ClientMsgID string                      `json:"client_msg_id,omitempty"` // 可选, 防重入id。对于同一个openid + client_msg_id, 只发送一条消息,10分钟有效,超过10分钟不保证效果。若无防重入需求，可不填
	ToUser      string                      `json:"touser"`                  // 必须, 接受者OpenID
	TemplateID  string                      `json:"template_id"`             // 必须, 模版ID
	URL         string                      `json:"url,omitempty"`           // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
	MiniProgram *MiniProgram                `json:"miniprogram,omitempty"`   // 可选, 跳小程序所需数据，不需跳小程序可不用传该数据
	Data        map[string]TemplateDataItem `json:"data"`                    // 必须, 模板数据, JSON 格式的 []byte, 满足特定的模板需求
	Color       string                      `json:"color,omitempty"`         // 可选, 模板内容字体颜色，不填默认为黑色
}

// MiniProgram 跳小程序所需数据
type MiniProgram struct {
	AppID    string `json:"appid"`              // 必选, 所需跳转到的小程序appid（该小程序 appid 必须与发模板消息的公众号是绑定关联关系，暂不支持小游戏）
	PagePath string `json:"pagepath,omitempty"` // 可选, 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar），要求该小程序已发布，暂不支持小游戏
}

// TemplateDataItem 模版内某个 .DATA 的值
type TemplateDataItem struct {
	Value string `json:"value"`           // 必选, 模板数据
	Color string `json:"color,omitempty"` // 可选, 模板内容字体颜色，不填默认为黑色
}

// templateSendResponse 发送模板消息的响应
type templateSendResponse struct {
	weixin.ResponseCode
	MsgID int64 `json:"msgid"` // 消息id
}

const templateSendURL = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"

// SendTemplate 发送模板消息
func SendTemplate(accessToken string, msg *Template) (msgID int64, err error) {
	url := fmt.Sprintf(templateSendURL, accessToken)
	var resp templateSendResponse
	_, err = client.PostJSON(url, msg, &resp)
	if err != nil {
		return
	}
	if !resp.Succeed() {
		err = fmt.Errorf("weixin request failed, uri=%q, resp=%+v", url, resp.ResponseCode)
	}
	msgID = resp.MsgID
	return
}
