package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type BotSwitch struct {
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
func (botSwitch *BotSwitch) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	s, b := Prefix(rawMsg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	gid := ""
	uid := ""
	if groupId.Common > 0 {
		gid = strconv.Itoa(int(groupId.Common))
	} else {
		gid = groupId.Offical
	}
	if userId.Common > 0 {
		uid = strconv.Itoa(int(userId.Common))
	} else {
		uid = userId.Offical
	}

	if StartsWith(s, "开启") && (userRole || super) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "开启"))
		if s == "开关" {
			log.Println("[开关] 不支持开启或关闭")
			return utils.RetStuct{
				RetVal:  utils.MESSAGE_BLOCK,
				ReqType: utils.GroupMsg,
			}
		}
		i := PluginNameToIntent(s)
		if i == 0 {
			reply := " 功能不存在"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GroupMsg,
			}
		}
		err := SwitchSave(gid, gid, uid, int64(i), time.Now(), false)
		if err != nil {
			reply := "开启失败"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GroupMsg,
			}
		} else {
			reply := "开启成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GroupMsg,
			}
		}
	}

	if StartsWith(s, "关闭") && (userRole || super) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "关闭"))
		if s == "开关" {
			log.Println("[开关] 不支持开启或关闭")
			return utils.RetStuct{
				RetVal:  utils.MESSAGE_BLOCK,
				ReqType: utils.GroupMsg,
			}
		}
		i := PluginNameToIntent(s)
		if i == 0 {
			reply := "功能不存在"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GroupMsg,
			}
		}
		err := SwitchSave(gid, gid, uid, int64(i), time.Now(), true)
		if err != nil {
			reply := "关闭失败"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GroupMsg,
			}
		} else {
			reply := "关闭成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GroupMsg,
			}
		}

	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("开关", &BotSwitch{})
}
