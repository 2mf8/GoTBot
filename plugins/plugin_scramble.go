package plugins

import (
	"context"
	"log"
	"net/url"
	"strings"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
	"github.com/2mf8/GoPbBot/proto_gen/onebot"
)

type ScramblePlugin struct {
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
func (scramble *ScramblePlugin) Do(ctx *context.Context, botId, groupId, userId int64, messageId *onebot.MessageReceipt, rawMsg, card string, botRole, userRole, super bool, rs, rd, rf int) RetStuct {

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}

	ins := Tnoodle(s).Instruction
	shor := Tnoodle(s).ShortName
	show := Tnoodle(s).ShowName
	if ins == s && ins != "instruction" {
		gs := GetScramble(shor)
		if StartsWith(gs, "net") || gs == "获取失败" {
			log.Printf("[INFO] Bot(%v) Group(%v) -> 获取打乱失败", botId, groupId)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
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
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v\n%v<image url=\"%v\"/>", botId, groupId, show, gs, imgUrl)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
				Text: sc,
				Image: imgUrl,
			},
			ReqType: GroupMsg,
		}
	}
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	Register("打乱", &ScramblePlugin{})
}
