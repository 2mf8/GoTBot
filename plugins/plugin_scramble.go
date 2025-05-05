package plugins

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
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
func (scramble *ScramblePlugin) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	s, b := Prefix(rawMsg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	ins := Tnoodle(s).Instruction
	shor := Tnoodle(s).ShortName
	show := Tnoodle(s).ShowName
	fmt.Println(ins, shor, show, s)
	if ins != "instruction" {
		gs := GetScramble(shor)
		if StartsWith(gs, "net") || gs == "获取失败" {
			log.Printf("[INFO] Bot(%v) Group(%v) -> 获取打乱失败", botId, groupId)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "获取打乱失败",
				},
				ReqType: utils.GroupMsg,
			}
		}
		if shor == "minx" {
			gs = strings.Replace(gs, "U' ", "#\n", -1)
			gs = strings.Replace(gs, "U ", "U\n", -1)
			gs = strings.Replace(gs, "#", "U'", -1)
		}
		imgUrl := fmt.Sprintf("%s/view/", AllConfig.ScrambleServer) + shor + ".png?scramble=" + url.QueryEscape(strings.Replace(gs, "\n", " ", -1))
		sc := show + "\n" + gs
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v\n%v<image url=\"%v\"/>", botId, groupId, show, gs, imgUrl)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text:  sc,
				Image: imgUrl,
			},
			ReqType: utils.GroupMsg,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

