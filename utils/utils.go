package utils

import (
	"context"

	"github.com/2mf8/GoPbBot/proto_gen/onebot"
)

type ReqType int

const (
	GroupBan    ReqType = iota // 群禁言
	RelieveBan                 // 禁言解除
	GroupKick                  // 群踢人
	GroupSignIn                // 群打卡
	GroupMsg                   // 群消息
	GroupLeave                 // 退群
	DeleteMsg                  // 消息撤回
	Undefined                  // 未定义
)

type RetStuct struct {
	RetVal         uint
	ReplyMsg       *Msg
	ReqType        ReqType
	Duration       int32
	BanId          int64
	RejectAddAgain bool
	MessageId      *onebot.MessageReceipt
}

type Msg struct {
	Text  string
	At    bool
	Image string
}

/*
* userId 用户Id
* groupId 群Id
* rawMsg 群消息
* userRole 用户角色，是否是管理员
* botRole 机器人角色， 是否是管理员
* retval 返回值，用于判断是否处理下一个插件
* replyMsg 待发送消息
* rs 成功防屏蔽码
* rd 删除防屏蔽码
* rf 失败防屏蔽码
 */
type InputStruct struct {
	BotId     int64
	GroupId   int64
	UserId    int64
	RawMsg    string
	UserRole  bool
	BotRole   bool
	SuperRole bool
	RS        int
	RD        int
	RF        int
}

type Plugin interface {
	Do(ctx *context.Context, botId, groupId, userId int64, messageId *onebot.MessageReceipt, rawMsg, card string, botRole, userRole, super bool, rs, rd, rf int) (retStuct RetStuct)
}

var PluginSet map[string]Plugin

const (
	MESSAGE_BLOCK  uint = 0
	MESSAGE_IGNORE uint = 1
)

func init() {
	PluginSet = make(map[string]Plugin)
}

func Register(k string, v Plugin) {
	PluginSet[k] = v
}
