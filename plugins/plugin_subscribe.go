package plugins

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	//. "github.com/2mf8/GoTBot/config"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type Sub struct {
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
func (sub *Sub) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	s, b := Prefix(rawMsg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	if StartsWith(s, "订阅") && (userRole || super) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "订阅"))
		SubscribeCreate(groupId.Offical, s)
		reply := " 订阅成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply},
			ReqType: utils.GroupMsg,
		}
	}
	if StartsWith(s, "取消订阅") && (userRole || super) {
		SubscribeDelete(groupId.Offical)
		reply := "取消订阅成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply},
			ReqType: utils.GroupMsg,
		}
	}

	reg1 := regexp.MustCompile("<at qq=\"")
	reg2 := regexp.MustCompile("\"/>")
	reg3 := regexp.MustCompile("<@")
	reg4 := regexp.MustCompile(">")
	reg5 := regexp.MustCompile("!")
	str1 := strings.TrimSpace(reg5.ReplaceAllString(reg4.ReplaceAllString(reg3.ReplaceAllString(reg1.ReplaceAllString(s, ""), ""), ""), ""))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, " "))
	if StartsWith(s, "授权") && (userRole || super) {
		as := strings.TrimSpace(strings.TrimPrefix(str2, "授权"))
		if as == "" {
			reply := " 授权失败，对象为空"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply},
				ReqType: utils.GroupMsg,
			}
		} else {
			auser := fmt.Sprintf("%s_%s", groupId.Offical, as)
			AuthCreate(auser, as)
			reply := " 授权成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply},
				ReqType: utils.GroupMsg,
			}
		}
	}
	if StartsWith(s, "取消授权") && (userRole || super) {
		as := strings.TrimSpace(strings.TrimPrefix(str2, "取消授权"))
		auser := fmt.Sprintf("%s_%s", groupId.Offical, as)
		AuthDelete(auser)
		reply := "取消授权成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply},
			ReqType: utils.GroupMsg,
		}
	}

	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}
