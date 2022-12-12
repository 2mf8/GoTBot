package plugins

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/2mf8/GoPbBot/proto_gen/onebot"
	. "github.com/2mf8/GoTBot/data"
	. "github.com/2mf8/GoTBot/public"
	. "github.com/2mf8/GoTBot/utils"
	"gopkg.in/guregu/null.v3"
)

type PricePlugin struct {
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
func (price *PricePlugin) Do(ctx *context.Context, botId, groupId, userId int64, messageId *onebot.MessageReceipt, rawMsg, card string, botRole, userRole, super bool, rs, rd, rf int) RetStuct {

	reg1 := regexp.MustCompile("％")
	reg2 := regexp.MustCompile("＃")
	reg3 := regexp.MustCompile("＆")
	reg4 := regexp.MustCompile("10001")
	reg5 := regexp.MustCompile("560820998")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(rawMsg, "%"))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, "#"))
	str3 := strings.TrimSpace(reg3.ReplaceAllString(str2, "&"))

	s, b := Prefix(str3, "%")
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

	if StartsWith(s, "#+") && (userRole || super) {
		str4 := strings.TrimSpace(string([]byte(s)[len("#+"):]))
		str5 := strings.Split(str4, "##")
		if len(str5) != 2 {
			if strings.TrimSpace(str5[0]) == "" {
				replyText := strconv.Itoa(rf) + "（商品名不能为空）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
						Text: replyText,
					},
					ReqType: GroupMsg,
				}
			}
			if groupId == 481097523 || groupId == 176211061 || groupId == 138080634 {
				err := IDBGAN("10001", "10001", str5[0])
				if err != nil {
					replyText := strconv.Itoa(rf) + "（删除失败）"
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					return RetStuct{
						RetVal: MESSAGE_BLOCK,
						ReplyMsg: &Msg{
							Text: replyText,
						},
						ReqType: GroupMsg,
					}
				}
				replyText := strconv.Itoa(rs) + "（删除成功）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
						Text: replyText,
					},
					ReqType: GroupMsg,
				}
			}
			if groupId == 560820998 || groupId == 189420325 || groupId == 348591755 {
				err := IDBGAN("10002", "10002", str5[0])
				if err != nil {
					replyText := strconv.Itoa(rf) + "（删除失败）"
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					return RetStuct{
						RetVal: MESSAGE_BLOCK,
						ReplyMsg: &Msg{
							Text: replyText,
						},
						ReqType: GroupMsg,
					}
				}
				replyText := strconv.Itoa(rs) + "（删除成功）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
						Text: replyText,
					},
					ReqType: GroupMsg,
				}
			}
			err := IDBGAN(strconv.Itoa(int(groupId)), strconv.Itoa(int(groupId)), str5[0])
			if err != nil {
				replyText := strconv.Itoa(rf) + "（删除失败）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
						Text: replyText,
					},
					ReqType: GroupMsg,
				}
			}
			replyText := strconv.Itoa(rs) + "（删除成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		if strings.TrimSpace(str5[0]) == "" {
			replyText := strconv.Itoa(rf) + "（商品名不能为空）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		str6 := strings.Split(str5[1], "#&")
		if len(str6) != 2 {
			if groupId == 481097523 || groupId == 176211061 || groupId == 138080634 {
				err := ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, strconv.Itoa(int(userId)), null.NewTime(time.Now(), true))
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
				replyText := strconv.Itoa(rs) + "（添加成功）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
						Text: replyText,
					},
					ReqType: GroupMsg,
				}
			}
			if groupId == 560820998 || groupId == 189420325 || groupId == 348591755 {
				err := ItemSave("10002", "10002", null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, strconv.Itoa(int(userId)), null.NewTime(time.Now(), true))
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
				replyText := strconv.Itoa(rs) + "（添加成功）"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return RetStuct{
					RetVal: MESSAGE_BLOCK,
					ReplyMsg: &Msg{
						Text: replyText,
					},
					ReqType: GroupMsg,
				}
			}
			err := ItemSave(strconv.Itoa(int(groupId)), strconv.Itoa(int(groupId)), null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, strconv.Itoa(int(userId)), null.NewTime(time.Now(), true))
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
			replyText := strconv.Itoa(rs) + "（添加成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		if groupId == 481097523 || groupId == 176211061 || groupId == 138080634 {
			err := ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), strconv.Itoa(int(userId)), null.NewTime(time.Now(), true))
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
			replyText := strconv.Itoa(rs) + "（添加成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		if groupId == 560820998 || groupId == 189420325 || groupId == 348591755 {
			err := ItemSave("10002", "10002", null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), strconv.Itoa(int(userId)), null.NewTime(time.Now(), true))
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
			replyText := strconv.Itoa(rs) + "（添加成功）"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return RetStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyText,
				},
				ReqType: GroupMsg,
			}
		}
		err := ItemSave(strconv.Itoa(int(groupId)), strconv.Itoa(int(groupId)), null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), strconv.Itoa(int(userId)), null.NewTime(time.Now(), true))
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
		replyText := strconv.Itoa(rs) + "（添加成功）"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
				Text: replyText,
			},
			ReqType: GroupMsg,
		}
	}
	cps := []CuberPrice{}
	ps := ""
	psc := ""
	ic := 0
	from := ""
	sub, err := GetSubscribe(groupId)
	if err != nil {
		from = strconv.Itoa(int(groupId))
		cps, _ = GetItems(strconv.Itoa(int(groupId)), strconv.Itoa(int(groupId)), s)
	} else {
		from = strings.TrimSpace(reg5.ReplaceAllString(reg4.ReplaceAllString(strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), "黄小姐的魔方店"), "奇乐魔方坊"))
		cps, _ = GetItems(strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), s)
	}
	for _, i := range cps {
		if i.Shipping.String == "" {
			ps += "\n" + i.Item + " | " + i.Price.String
		} else {
			ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
		}
		if ic == 19 {
			ps += "\n..."
			break
		}
		ic++
	}
	if len(cps) == 0 {
		replyText := strconv.Itoa(rs) + "（暂无相关记录）"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
				Text: replyText,
			},
			ReqType: GroupMsg,
		}
	} else {
		psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + "\n品名 | 价格 | 备注" + ps + "\n价格源自 " + from
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, psc)
		return RetStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
				Text: psc,
			},
			ReqType: GroupMsg,
		}
	}
}

func init() {
	Register("查价", &PricePlugin{})
}
