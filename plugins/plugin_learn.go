package plugins

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
	. "github.com/2mf8/go-tbot-for-rq/data"
	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
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
func (learnPlugin *LearnPlugin) Do(ctx *context.Context, botId, groupId, userId int64, messageId *onebot.MessageReceipt, rawMsg, card string, botRole, userRole, super bool, rs, rd, rf int) RetStuct {

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}
	ggk, _ := GetJudgeKeys()
	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		msg := strconv.Itoa(rf) + " （消息触发守卫，已被拦截）"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
					Text: msg,
				},
			ReqType: GroupMsg,
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
				replyText := strconv.Itoa(rf) + "（问指令不能为空）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)			
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			err := LDBGAA(groupId, str3[0])
			if err != nil {
				replyText := strconv.Itoa(rf) + "（问答删除失败）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			replyText := strconv.Itoa(rs) + "（问答删除成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := strconv.Itoa(rf) + "（问指令不能为空）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		err := LearnSave(strings.TrimSpace(str3[0]), groupId, userId, null.NewString(str3[1], true), time.Now())
		if err != nil {
			replyText := strconv.Itoa(rf) + "（添加失败）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		replyText := strconv.Itoa(rs) + "（学习已完成，下次触发有效）"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
					Text: replyText,
				},
			ReqType: GroupMsg,
		}
	}
	if StartsWith(str1, "++") && super {
		str2 := strings.TrimSpace(strings.TrimPrefix(str1, "++"))
		str3 := strings.Split(str2, "##")
		if len(str3) != 2 {
			if strings.TrimSpace(str3[0]) == "" {
				replyText := strconv.Itoa(rf) + "（系统问指令不能为空）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			err := LDBGAA(int64(9999999990), str3[0])
			if err != nil {
				replyText := strconv.Itoa(rf) + "（系统问答删除失败）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			replyText := strconv.Itoa(rs) + "（系统问答删除成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := strconv.Itoa(rf) + "（系统问指令不能为空）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		err := LearnSave(strings.TrimSpace(str3[0]), int64(9999999990), userId, null.NewString(str3[1], true), time.Now())
		if err != nil {
			replyText := strconv.Itoa(rf) + "（系统问答添加失败）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		replyText := strconv.Itoa(rs) + "（系统问答学习已完成，下次触发有效）"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
					Text: replyText,
				},
			ReqType: GroupMsg,
		}
	}
	if strings.TrimSpace(rawMsg) == "" {
		replyText := strconv.Itoa(rf) + "（指令不能为空）"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
					Text: replyText,
				},
			ReqType: GroupMsg,
		}
	}
	learn_get, err := LearnGet(groupId, strings.TrimSpace(s))
	//log.Println(learn_get.LearnSync.Answer.String,"ceshil", err)
	if err != nil || learn_get.LearnSync.Answer.String == "" {
		sys_learn_get, _ := LearnGet(int64(9999999990), strings.TrimSpace(s))
		if sys_learn_get.LearnSync.Answer.String != "" {
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, sys_learn_get.LearnSync.Answer.String)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: sys_learn_get.LearnSync.Answer.String,
				},
				ReqType: GroupMsg,
			}
		}
	}
	if learn_get.LearnSync.Answer.String != "" {
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, learn_get.LearnSync.Answer.String)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
					Text: learn_get.LearnSync.Answer.String,
				},
			ReqType: GroupMsg,
		}
	}
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}
func init() {
	Register("学习", &LearnPlugin{})
}
