package plugins

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
	"gopkg.in/guregu/null.v3"
)

type CLearnPlugin struct {
}

func (learnPlugin *CLearnPlugin) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super, userRole bool, rs, rd, rf int) (retStuct RetChannelStuct) {

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetChannelStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}
	ggk, _ := GetJudgeKeys()
	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		msg := strconv.Itoa(rf) + " （消息触发守卫，已被拦截）"
		log.Printf("[守卫] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, msg)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
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
				log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)		
				return RetChannelStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			err := LDBGAA(int64(channelId), str3[0])
			if err != nil {
				replyText := strconv.Itoa(rf) + "（问答删除失败）"
				log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
				return RetChannelStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			replyText := strconv.Itoa(rs) + "（问答删除成功）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := strconv.Itoa(rf) + "（问指令不能为空）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		err := LearnSave(strings.TrimSpace(str3[0]), int64(channelId), int64(userId), null.NewString(str3[1], true), time.Now())
		if err != nil {
			replyText := strconv.Itoa(rf) + "（添加失败）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		replyText := strconv.Itoa(rs) + "（学习已完成，下次触发有效）"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
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
				log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
				return RetChannelStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			err := LDBGAA(int64(9999999990), str3[0])
			if err != nil {
				replyText := strconv.Itoa(rf) + "（系统问答删除失败）"
				log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
				return RetChannelStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
					ReqType: GroupMsg,
				}
			}
			replyText := strconv.Itoa(rs) + "（系统问答删除成功）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := strconv.Itoa(rf) + "（系统问指令不能为空）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		err := LearnSave(strings.TrimSpace(str3[0]), int64(9999999990), int64(userId), null.NewString(str3[1], true), time.Now())
		if err != nil {
			replyText := strconv.Itoa(rf) + "（系统问答添加失败）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		replyText := strconv.Itoa(rs) + "（系统问答学习已完成，下次触发有效）"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
			ReqType: GroupMsg,
		}
	}
	if strings.TrimSpace(rawMsg) == "" {
		replyText := strconv.Itoa(rf) + "（指令不能为空）"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyText)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
					Text: replyText,
				},
			ReqType: GroupMsg,
		}
	}
	learn_get, err := LearnGet(int64(channelId), strings.TrimSpace(s))
	//log.Println(learn_get.LearnSync.Answer.String,"ceshil", err)
	if err != nil || learn_get.LearnSync.Answer.String == "" {
		sys_learn_get, _ := LearnGet(int64(9999999990), strings.TrimSpace(s))
		if sys_learn_get.LearnSync.Answer.String != "" {
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, sys_learn_get.LearnSync.Answer.String)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &ChannelMsg{
					Text: sys_learn_get.LearnSync.Answer.String,
				},
				ReqType: GroupMsg,
			}
		}
	}
	if learn_get.LearnSync.Answer.String != "" {
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, learn_get.LearnSync.Answer.String)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &ChannelMsg{
					Text: learn_get.LearnSync.Answer.String,
				},
			ReqType: GroupMsg,
		}
	}
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}
func init() {
	ChannelRegister("学习", &CLearnPlugin{})
}
