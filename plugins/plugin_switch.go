package plugins

import (
	"context"
	"log"
	"strconv"
	"strings"
	"github.com/2mf8/GoPbBot/proto_gen/onebot"
	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
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
func (botSwitch *BotSwitch) Do(ctx *context.Context, botId, groupId, userId int64, messageId *onebot.MessageReceipt, rawMsg, card string, botRole, userRole, super bool, rs, rd, rf int) RetStuct {

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}

	if StartsWith(s, "开启") && (userRole || super) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "开启"))
		if s == "开关" {
			log.Println("[开关] 不支持开启或关闭")
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReqType: GroupMsg,
			}
		}
		i := PluginNameToIntent(s)
		if i == 0 {
			reply := strconv.Itoa(rf) + " （功能不存在）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}
		err := SwitchSave(groupId, int64(i), false)
		if err != nil {
			reply := strconv.Itoa(rf) + " （开启失败）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		} else {
			reply := strconv.Itoa(rs) + " （开启成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}
	}

	if StartsWith(s, "关闭") && (userRole || super) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "关闭"))
		if s == "开关" {
			log.Println("[开关] 不支持开启或关闭")
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReqType: GroupMsg,
			}
		}
		i := PluginNameToIntent(s)
		if i == 0 {
			reply := strconv.Itoa(rf) + " （功能不存在）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}
		err := SwitchSave(groupId, int64(i), true)
		if err != nil {
			reply := strconv.Itoa(rf) + " （关闭失败）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		} else {
			reply := strconv.Itoa(rs) + " （关闭成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
			
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: reply,
				},
				ReqType: GroupMsg,
			}
		}

	}
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	Register("开关", &BotSwitch{})
}
