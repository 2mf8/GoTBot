package plugins

import (
	"context"

	. "github.com/2mf8/GoTBot/utils"
)

type Reply struct {
}

/*
* botId 机器人Id
* groupId 群Id
* userId 用户Id
* messageId 消息Id
* rawMsg 群消息
* card At展示
* userRole 用户角色，是否是管理员
* botRole 机器人角色， 是否是管理员
* retval 返回值，用于判断是否处理下一个插件
* replyMsg 待发送消息
* rs 成功防屏蔽码
* rd 删除防屏蔽码
* rf 失败防屏蔽码
 */
func (rep *Reply) Do(ctx *context.Context, botId, groupId, userId int64, groupName string, messageId int64, rawMsg, card string, botRole, userRole, super bool) RetStuct {
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	Register("回复", &Reply{})
}
