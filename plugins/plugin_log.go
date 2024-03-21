package plugins

import (
	"context"
	"strconv"

	"github.com/2mf8/GoTBot/utils"
	. "github.com/2mf8/GoTBot/utils"
	log "github.com/sirupsen/logrus"
)

type Log struct {
}

func (l *Log) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	gid := ""
	uid := ""
	if groupId.Common > 0 {
		gid = strconv.Itoa(int(groupId.Common))
	} else {
		gid = groupId.Offical
	}
	if userId.Common > 0 {
		uid = strconv.Itoa(int(userId.Common))
	} else {
		uid = userId.Offical
	}
	if botId.Common > 0 {
		log.Infof("BotId(%v) GroupId(%s) UserId(%s) <- %s", botId, gid, uid, rawMsg)
	} else {
		log.Infof("GroupId(%s) UserId(%s) <- %s", gid, uid, rawMsg)
	}
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	Register("Log", &Log{})
}
