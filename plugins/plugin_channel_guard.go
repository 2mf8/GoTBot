package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"
	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
)

type ChannelGuard struct {
}

func (guard *ChannelGuard) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super, userRole bool, rs, rd, rf int) (retStuct RetChannelStuct){

	ggk, _ := GetJudgeKeys()

	if StartsWith(rawMsg, ".拦截") && super {
		vocabulary := strings.TrimPrefix(rawMsg, ".拦截")
		content := strings.Split(vocabulary, " ")
		err := ggk.JudgeKeysUpdate(content...)
		if err != nil {
			log.Panicln(err)
		}
		msg := strconv.Itoa(rs) + " （拦截词汇添加成功）"
		log.Printf("[守卫] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, msg)
		return RetChannelStuct{
			RetVal:   MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
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
		log.Printf("[守卫] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, msg)
		return RetChannelStuct{
			RetVal:   MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
					Text: msg,
				},
			ReqType:  GroupMsg,
		}
	}
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	ChannelRegister("守卫", &ChannelGuard{})
}
