package message

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/lenye/pmsg/pkg/flags"
	"github.com/lenye/pmsg/pkg/http/client"
	"github.com/lenye/pmsg/pkg/weixin"
	"github.com/lenye/pmsg/pkg/weixin/work/token"
)

type CmdWorkSendCustomerParams struct {
	UserAgent   string
	AccessToken string
	CorpID      string
	CorpSecret  string
	ToUser      string
	OpenKfID    string
	MsgID       string
	MsgType     string
	Data        string
}

func (t *CmdWorkSendCustomerParams) Validate() error {
	if t.AccessToken == "" && t.CorpID == "" {
		return flags.ErrWeixinWorkAccessToken
	}

	if err := ValidateCustomerMsgType(t.MsgType); err != nil {
		return fmt.Errorf("invalid flags %s: %v", flags.MsgType, err)
	}

	return nil
}

// CmdWorkSendCustomer 发送微信客服消息
func CmdWorkSendCustomer(arg *CmdWorkSendCustomerParams) error {

	if err := arg.Validate(); err != nil {
		return err
	}

	msg := CustomerMessage{
		ToUser:   arg.ToUser,
		OpenKfID: arg.OpenKfID,
		MsgID:    arg.MsgID,
		MsgType:  arg.MsgType,
	}

	buf := bytes.NewBufferString("")
	buf.WriteString(arg.Data)
	switch arg.MsgType {
	case CustomerMsgTypeText:
		var msgMeta TextMeta
		msgMeta.Content = buf.String()
		msg.Text = &msgMeta
	case CustomerMsgTypeImage:
		var msgMeta ImageMeta
		msgMeta.MediaID = buf.String()
		msg.Image = &msgMeta
	case CustomerMsgTypeVoice:
		var msgMeta VoiceMeta
		msgMeta.MediaID = buf.String()
		msg.Voice = &msgMeta
	case CustomerMsgTypeVideo:
		var msgMeta VideoMeta
		msgMeta.MediaID = buf.String()
		msg.Video = &msgMeta
	case CustomerMsgTypeFile:
		var msgMeta FileMeta
		msgMeta.MediaID = buf.String()
		msg.File = &msgMeta
	case CustomerMsgTypeLink:
		var msgMeta LinkMeta
		if err := json.Unmarshal(buf.Bytes(), &msgMeta); err != nil {
			return fmt.Errorf("invalid json format, %v", err)
		}
		msg.Link = &msgMeta
	case CustomerMsgTypeMiniProgram:
		var msgMeta MiniProgramMeta
		if err := json.Unmarshal(buf.Bytes(), &msgMeta); err != nil {
			return fmt.Errorf("invalid json format, %v", err)
		}
		msg.MiniProgram = &msgMeta
	case CustomerMsgTypeMsgMenu:
		var msgMeta MsgMenuMeta
		if err := json.Unmarshal(buf.Bytes(), &msgMeta); err != nil {
			return fmt.Errorf("invalid json format, %v", err)
		}
		msg.MsgMenu = &msgMeta
	case CustomerMsgTypeLocation:
		var msgMeta LocationMeta
		if err := json.Unmarshal(buf.Bytes(), &msgMeta); err != nil {
			return fmt.Errorf("invalid json format, %v", err)
		}
		msg.Location = &msgMeta
	}

	client.UserAgent = arg.UserAgent

	if arg.AccessToken == "" {
		accessTokenResp, err := token.GetAccessToken(arg.CorpID, arg.CorpSecret)
		if err != nil {
			return err
		}
		arg.AccessToken = accessTokenResp.AccessToken
	}

	if resp, err := SendCustomer(arg.AccessToken, &msg); err != nil {
		return err
	} else {
		fmt.Println(fmt.Sprintf("%v; %v", weixin.MessageOK, resp))
	}

	return nil
}