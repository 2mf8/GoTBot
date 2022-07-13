package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
	. "github.com/2mf8/go-tbot-for-rq/data"
	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
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
func (guard *Guard) Do(ctx *context.Context, botId, groupId, userId int64, messageId *onebot.MessageReceipt, rawMsg, card string, botRole, userRole, super bool, rs, rd, rf int) RetStuct {
	if !botRole {
		return RetStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}
	guardIntent := int64(PluginGuard)
	sg, _ := SGBGI(groupId)
	isGuard := sg.PluginSwitch.IsCloseOrGuard & guardIntent

	if isGuard > 0 {
		return RetStuct{
			RetVal: MESSAGE_IGNORE,
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
		msg := strconv.Itoa(rs) + " （拦截词汇添加成功）"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return RetStuct{
			RetVal:   MESSAGE_BLOCK,
			ReplyMsg: &Msg{
					Text: msg,
				},
			ReqType:  GroupMsg,
		}
	}

	if StartsWith(rawMsg, ".取消拦截") && super {
		vocabulary := strings.TrimPrefix(rawMsg, ".取消拦截")
		content := strings.Split(vocabulary, " ")
		ggk.JudgeKeysDelete(content...)
		msg := strconv.Itoa(rd) + " （拦截词汇删除成功）"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return RetStuct{
			RetVal:   MESSAGE_BLOCK,
			ReplyMsg: &Msg{
					Text: msg,
				},
			ReqType:  GroupMsg,
		}
	}

	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		if userRole {
			msg := strconv.Itoa(rs) + " （消息触发守卫，已被拦截）"
			log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
			return RetStuct{
				RetVal:   MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: msg,
				},
				ReqType:  GroupMsg,
			}
		}
		msg := strconv.Itoa(rs) + " （消息触发守卫，已撤回消息并禁言该用户两分钟, 请文明发言）"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return RetStuct{
			RetVal:    MESSAGE_BLOCK,
			ReplyMsg:  &Msg{
					Text: msg,
				},
			ReqType:   DeleteMsg,
			Duration:  int32(120),
			MessageId: messageId,
		}
	}
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	Register("守卫", &Guard{})
}
