package plugins

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
)

type CPricePlugin struct {
}

func (price *CPricePlugin) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super, userRole bool, rs, rd, rf int) (retStuct RetChannelStuct) {

	reg1 := regexp.MustCompile("％")
	reg2 := regexp.MustCompile("＃")
	reg3 := regexp.MustCompile("＆")
	reg4 := regexp.MustCompile("10001")
	reg5 := regexp.MustCompile("10002")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(rawMsg, "%"))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, "#"))
	str3 := strings.TrimSpace(reg3.ReplaceAllString(str2, "&"))

	s, b := Prefix(str3, "%")
	if !b {
		return RetChannelStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}

	cps := []CuberPrice{}
	ps := ""
	psc := ""
	ic := 0
	from := ""
	sub, _ := GetSubscribe(int64(channelId))
	from = strings.TrimSpace(reg5.ReplaceAllString(reg4.ReplaceAllString(strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), "黄小姐的魔方店"), "奇乐魔方坊"))
	cps, _ = GetItems(strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), s)
	for _, i := range cps {
		if i.Shipping.String == "" {
			ps += "\n" + i.Item + " | " + i.Price.String
		} else {
			ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
		}
		if ic == 19 {
			ps += "\n..."
			break
		}
		ic++
	}
	if len(cps) == 0 {
		replyText := strconv.Itoa(rs) + "（暂无相关记录）"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
				Text: replyText,
			},
			ReqType: GroupMsg,
		}
	} else {
		psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + "\n品名 | 价格 | 备注" + ps + "\n价格源自 " + from
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, psc)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
				Text: psc,
			},
			ReqType: GroupMsg,
		}
	}
}

func init() {
	ChannelRegister("查价", &CPricePlugin{})
}
