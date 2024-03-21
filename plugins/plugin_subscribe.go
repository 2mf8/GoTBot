package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"

	//. "github.com/2mf8/GoTBot/config"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type Sub struct {
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
func (sub *Sub) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	s, b := Prefix(rawMsg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	if StartsWith(s, "订阅") && (userRole || super) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "订阅"))
		r_groupId, _ := strconv.Atoi(s)
		_ = SubSave(groupId.Common, int64(r_groupId), userId.Common)
		reply := " 订阅成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply,
			},
			ReqType: utils.GroupMsg,
			OfficalMsgId: messageId.Offical,
		}
	}
	if StartsWith(s, "取消订阅") && (userRole || super) {
		_ = SubDeleteByGroupId(groupId.Common)
		reply := "取消订阅成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply,
			},
			ReqType: utils.GroupMsg,
			OfficalMsgId: messageId.Offical,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("订阅", &Sub{})
}
