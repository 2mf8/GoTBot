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
	"gopkg.in/guregu/null.v3"
)

type LearnPlugin struct {
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
func (learnPlugin *LearnPlugin) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
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
	ggk, _ := GetJudgeKeys()
	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		msg := "消息触发守卫，已被拦截"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType: utils.GroupMsg,
		}
	}
	reg1 := regexp.MustCompile("＃")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, "#"))
	if StartsWith(str1, "#+") && (super || userRole) {
		//if StartsWith(str1, "#+") && super {
		str2 := strings.TrimSpace(strings.TrimPrefix(str1, "#+"))
		str3 := strings.Split(str2, "##")
		if len(str3) != 2 {
			if strings.TrimSpace(str3[0]) == "" {
				replyText := "问指令不能为空"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			err := LDBGAA(gid, gid, str3[0])
			if err != nil {
				replyText := "问答删除失败"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			replyText := "问答删除成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := "问指令不能为空"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		err := LearnSave(strings.TrimSpace(str3[0]), gid, gid, uid, null.NewString(str3[1], true), time.Now(), true)
		if err != nil {
			replyText := "添加失败"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		replyText := "学习已完成，下次触发有效"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyText},
			ReqType: utils.GroupMsg,
		}
	}
	if StartsWith(str1, "++") && super {
		str2 := strings.TrimSpace(strings.TrimPrefix(str1, "++"))
		str3 := strings.Split(str2, "##")
		if len(str3) != 2 {
			if strings.TrimSpace(str3[0]) == "" {
				replyText := "系统问指令不能为空"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			err := LDBGAA("9999999990", "9999999990", str3[0])
			if err != nil {
				replyText := "系统问答删除失败"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			replyText := "系统问答删除成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := "系统问指令不能为空"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		err := LearnSave(strings.TrimSpace(str3[0]), "9999999990", "9999999990", uid, null.NewString(str3[1], true), time.Now(), true)
		fmt.Println(err)
		if err != nil {
			replyText := "系统问答添加失败"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		replyText := "系统问答学习已完成，下次触发有效"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyText},
			ReqType: utils.GroupMsg,
		}
	}
	if strings.TrimSpace(rawMsg) == "" {
		replyText := "指令不能为空"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyText},
			ReqType: utils.GroupMsg,
		}
	}
	learn_get, err := LearnGet(gid, gid, strings.TrimSpace(s))
	//log.Println(learn_get.LearnSync.Answer.String,"ceshil", err)
	if err != nil || learn_get.Answer.String == "" {
		sys_learn_get, _ := LearnGet("9999999990", "9999999990", strings.TrimSpace(s))
		if sys_learn_get.Answer.String != "" {
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, sys_learn_get.Answer.String)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: sys_learn_get.Answer.String,
				},
				ReqType: utils.GroupMsg,
			}
		}
	}
	if learn_get.Answer.String != "" {
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, learn_get.Answer.String)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: learn_get.Answer.String,
			},
			ReqType: utils.GroupMsg,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}
func init() {
	utils.Register("学习", &LearnPlugin{})
}
