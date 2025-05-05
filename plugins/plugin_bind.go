package plugins

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	database "github.com/2mf8/GoTBot/data"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
	. "github.com/2mf8/GoTBot/utils"
	"gopkg.in/guregu/null.v3"
)

type Bind struct {
}

func (rep *Bind) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	if !public.Contains(rawMsg, "中") && public.Contains(rawMsg, "bind") && public.Contains(rawMsg, "[CQ:at,qq=2854216320]") {
		ss := strings.Split(rawMsg, "-")
		if len(ss) != 2 {
			replyText := "Bind错误"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		ns := strings.TrimSpace(ss[1])
		if ns != fmt.Sprintf("%v", userId) {
			replyText := "Bind错误"
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		key := fmt.Sprintf("%v_bind", userId)
		time.Sleep(time.Second)
		v, _ := database.RedisGet(key)
		if string(v) != "" {
			err := database.BindUserInfoSave(string(v), null.NewString(card, true), null.NewString("", true), userId.Common)
			if err != nil {
				replyText := "Bind错误"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			replyText := "Bind成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}

	}
	return RetStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

