package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"

	//. "github.com/2mf8/GoTBot/config"
	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
)

type CSub struct {
}

func (sub *CSub) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super, userRole bool, rs, rd, rf int) (retStuct RetChannelStuct) {

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetChannelStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}

	if StartsWith(s, "订阅") && super {
		s = strings.TrimSpace(strings.TrimPrefix(s, "订阅"))
		r_channelId, _ := strconv.Atoi(s)
		_ = SubSave(int64(channelId), int64(r_channelId), int64(userId))
		reply := strconv.Itoa(rs) + " （订阅成功）"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, reply)
		return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
	}
	if StartsWith(s, "取消订阅") && super {
		_ = SubDeleteByGroupId(int64(channelId))
		reply := strconv.Itoa(rd) + " （取消订阅成功）"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, reply)
		return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
	}
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	ChannelRegister("订阅", &CSub{})
}
