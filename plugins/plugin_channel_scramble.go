package plugins

import (
	"context"
	"log"
	"net/url"
	"strings"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
)

type CScramblePlugin struct {
}

func (scramble *CScramblePlugin) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super, userRole bool, rs, rd, rf int) (retStuct RetChannelStuct) {

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetChannelStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}

	ins := Tnoodle(s).Instruction
	shor := Tnoodle(s).ShortName
	show := Tnoodle(s).ShowName
	if ins == s && ins != "instruction" {
		gs := GetScramble(shor)
		if StartsWith(gs, "net") || gs == "获取失败" {
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> 获取打乱失败", botId, guildId, channelId)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: "获取打乱失败",
				},
				ReqType: GroupMsg,
			}
		}
		if shor == "minx" {
			gs = strings.Replace(gs, "U' ", "#\n", -1)
			gs = strings.Replace(gs, "U ", "U\n", -1)
			gs = strings.Replace(gs, "#", "U'", -1)
		}
		imgUrl := "http://www.2mf8.cn:2014/view/" + shor + ".png?scramble=" + url.QueryEscape(strings.Replace(gs, "\n", " ", -1))
		sc := show + "\n" + gs
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v\n%v\n<guild_image file=\"%v\" url=\"%v\" />", botId, guildId, channelId, show, gs, shor+".png", imgUrl)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
				Text:  sc,
				Image: imgUrl,
				File:  shor + ".png",
			},
			ReqType: GroupMsg,
		}
	}
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	ChannelRegister("打乱", &CScramblePlugin{})
}
