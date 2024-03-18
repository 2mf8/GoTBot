package plugins

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type Block struct{}

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
func (block *Block) Do(ctx *context.Context, botId, groupId, userId int64, groupName string, messageId int64, rawMsg, card string, botRole, userRole, super bool) utils.RetStuct {
	gid := fmt.Sprintf("%v", groupId)
	uid := fmt.Sprintf("%v", userId)
	ispblock, err := PBlockGet(gid, uid)
	//fmt.Println(ispblock)
	if err != nil {
		fmt.Println("[INFO] ", err)
	}
	if ispblock.UserId == strconv.Itoa(int(userId)) && ispblock.IsPBlock {
		if !super {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
			}
		}
	}

	s, b := Prefix(rawMsg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	reg1 := regexp.MustCompile("<at qq=\"")
	reg2 := regexp.MustCompile("\"/>")
	reg3 := regexp.MustCompile("  ")

	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, ""))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, " "))

	for Contains(str2, "  ") {
		str2 = strings.TrimSpace(reg3.ReplaceAllString(str2, " "))
	}

	if StartsWith(s, "屏蔽+") && super {
		pUserID := strings.TrimSpace(strings.TrimPrefix(str2, "屏蔽+"))
		if err != nil {
			replyMsg := "用户不存在"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyMsg,
				},
				ReqType: utils.GroupMsg,
			}
		}
		pid := fmt.Sprintf("%v", pUserID)
		err = PBlockSave(gid, pid, uid, true, time.Now())
		if err != nil {
			replyMsg := "屏蔽" + pUserID + "失败"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyMsg,
				},
				ReqType: utils.GroupMsg,
			}
		}
		replyMsg := "屏蔽" + pUserID + "成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyMsg,
			},
			ReqType: utils.GroupMsg,
		}
	}
	if StartsWith(s, "屏蔽-") && super {
		pUserID := strings.TrimSpace(strings.TrimPrefix(str2, "屏蔽-"))
		if err != nil {
			replyMsg := "用户不存在"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyMsg,
				},
				ReqType: utils.GroupMsg,
			}
		}
		err = PBlockSave(gid, pUserID, uid, false, time.Now())
		if err != nil {
			replyMsg := "解除屏蔽" + pUserID + "失败"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyMsg,
				},
				ReqType: utils.GroupMsg,
			}
		}
		replyMsg := "解除屏蔽" + pUserID + "成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyMsg,
			},
			ReqType: utils.GroupMsg,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("屏蔽", &Block{})
}
