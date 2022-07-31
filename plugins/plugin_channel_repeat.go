package plugins

import (
	"context"
	"log"
	"math/rand"
	"time"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
	. "github.com/2mf8/GoTBot/data"
)

type CRepeat struct {
}
/*
* botId 机器人Id
* channelId 群Id
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
func (rep *CRepeat) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super, userRole bool, rs, rd, rf int) (retStuct RetChannelStuct) {

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(101)

	ggk, _ := GetJudgeKeys()
	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		msg := "消息触发守卫，已被拦截"
		log.Printf("[复读守卫] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, msg)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
		}
	}

	if len(rawMsg) < 20 && r%70 == 0 && !(StartsWith(rawMsg, ".") || StartsWith(rawMsg, "%") || StartsWith(rawMsg, "％")) {
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, rawMsg)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
				Text: rawMsg,
			},
			ReqType: GroupMsg,
		}
	}
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	ChannelRegister("复读", &CRepeat{})
}
