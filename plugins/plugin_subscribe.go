package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"

	//. "github.com/2mf8/go-tbot-for-rq/config"
	. "github.com/2mf8/go-tbot-for-rq/data"
	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
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
func (sub *Sub) Do(ctx *context.Context, botId, groupId, userId int64, messageId *onebot.MessageReceipt, rawMsg, card string, botRole, userRole, super bool, rs, rd, rf int) RetStuct {

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}

	if StartsWith(s, "订阅") && (userRole || super) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "订阅"))
		r_groupId, _ := strconv.Atoi(s)
		_ = SubSave(groupId, int64(r_groupId), userId)
		reply := strconv.Itoa(rs) + " （订阅成功）"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
	}
	if StartsWith(s, "取消订阅") && (userRole || super) {
		_ = SubDeleteByGroupId(groupId)
		reply := strconv.Itoa(rd) + " （取消订阅成功）"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
	}
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	Register("订阅", &Sub{})
}
