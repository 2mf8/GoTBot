package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type Guard struct {
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
func (guard *Guard) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	if !botRole {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	gid := ""
	if groupId.Common > 0 {
		gid = strconv.Itoa(int(groupId.Common))
	} else {
		gid = groupId.Offical
	}
	guardIntent := int64(PluginGuard)
	sg, _ := SGBGIACI(gid, gid)
	isGuard := sg.PluginSwitch.IsCloseOrGuard & guardIntent

	if isGuard > 0 {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	ggk, _ := GetJudgeKeys()

	if StartsWith(rawMsg, ".拦截") && (userRole || super) {
		vocabulary := strings.TrimPrefix(rawMsg, ".拦截")
		content := strings.Split(vocabulary, " ")
		err := ggk.JudgeKeysUpdate(content...)
		if err != nil {
			log.Panicln(err)
		}
		msg := "拦截词汇添加成功"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType: utils.GroupMsg,
		}
	}

	if StartsWith(rawMsg, ".取消拦截") && super {
		vocabulary := strings.TrimPrefix(rawMsg, ".取消拦截")
		content := strings.Split(vocabulary, " ")
		ggk.JudgeKeysDelete(content...)
		msg := "拦截词汇删除成功"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType: utils.GroupMsg,
		}
	}

	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		if userRole {
			msg := "消息触发守卫，已被拦截"
			log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: msg,
				},
				ReqType: utils.GroupMsg,
			}
		}
		msg := "消息触发守卫，已撤回消息并禁言该用户两分钟, 请文明发言"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType:  utils.DeleteMsg,
			Duration: int64(120),
			MsgId:    messageId.Common,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("守卫", &Guard{})
}
