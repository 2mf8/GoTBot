package utils

import (
	"context"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

type ReqType int

const (
	GroupBan   ReqType = iota // 频道禁言
	RelieveBan                // 禁言解除
	GroupKick                 // 频道踢人
	GroupMsg                  // 频道消息
	GroupLeave                // 退频道
	DeleteMsg                 // 消息撤回
	Undefined                 // 未定义
)

type RetStuct struct {
	RetVal         uint
	ReplyMsg       *Msg
	ReqType        ReqType
	Duration       int64
	BanId          int64
	RejectAddAgain bool
	Retract        int
	MsgId          int64
	OfficalMsgId   string
}

type Msg struct {
	Text   string
	At     bool
	Image  string
	Images []string
}

type GroupIdType struct {
	Common  int64
	Offical string
}

type UserIdType struct {
	Common  int64
	Offical string
}

type MsgIdType struct {
	Common  int64
	Offical string
}

type BotIdType struct {
	Common  int64
	Offical string
}

type Plugin interface {
	Do(ctx *context.Context, botId *BotIdType, groupId *GroupIdType, userId *UserIdType, groupName string, messageId *MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct RetStuct)
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

func FatalError(err error) {
	log.Errorf(err.Error())
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	sBuf := string(buf)
	log.Errorf(sBuf)
	time.Sleep(5 * time.Second)
	panic(err)
}
