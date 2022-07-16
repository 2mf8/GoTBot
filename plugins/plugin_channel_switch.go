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

type CBotSwitch struct {
}

func (botSwitch *CBotSwitch) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super bool, rs, rd, rf int) (retStuct RetChannelStuct){

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetChannelStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}

	if StartsWith(s, "开启") && super {
		s = strings.TrimSpace(strings.TrimPrefix(s, "开启"))
		if s == "开关" {
			log.Println("[开关] 不支持开启或关闭")
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReqType: GroupMsg,
			}
		}
		i := PluginNameToIntent(s)
		if i == 0 {
			reply := strconv.Itoa(rf) + " （功能不存在）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, rawMsg)
			
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}
		err := SwitchSave(int64(channelId), int64(i), false)
		if err != nil {
			reply := strconv.Itoa(rf) + " （开启失败）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, rawMsg)
			
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		} else {
			reply := strconv.Itoa(rs) + " （开启成功）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, rawMsg)
			
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}
	}

	if StartsWith(s, "关闭") && super {
		s = strings.TrimSpace(strings.TrimPrefix(s, "关闭"))
		if s == "开关" {
			log.Println("[开关] 不支持开启或关闭")
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReqType: GroupMsg,
			}
		}
		i := PluginNameToIntent(s)
		if i == 0 {
			reply := strconv.Itoa(rf) + " （功能不存在）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, rawMsg)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}
		err := SwitchSave(int64(channelId), int64(i), true)
		if err != nil {
			reply := strconv.Itoa(rf) + " （关闭失败）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, rawMsg)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		} else {
			reply := strconv.Itoa(rs) + " （关闭成功）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, rawMsg)
			
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}

	}
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	ChannelRegister("开关", &CBotSwitch{})
}
