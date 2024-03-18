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
func (price *PricePlugin) Do(ctx *context.Context, botId, groupId, userId int64, groupName string, messageId int64, rawMsg, card string, botRole, userRole, super bool) utils.RetStuct {
	uid := fmt.Sprintf("%v", userId)
	reg1 := regexp.MustCompile("％")
	reg2 := regexp.MustCompile("＃")
	reg3 := regexp.MustCompile("＆")
	reg4 := regexp.MustCompile("10001")
	reg5 := regexp.MustCompile("10002")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(rawMsg, "%"))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, "#"))
	str3 := strings.TrimSpace(reg3.ReplaceAllString(str2, "&"))
	isMagnetism := strings.Contains(rawMsg, "磁")

	s, b := Prefix(str3, "%")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	ggk, _ := GetJudgeKeys()
	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		msg := " 消息触发守卫，已被拦截"
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType: utils.GroupMsg,
		}
	}

	if StartsWith(s, "#+") && (userRole || super) {
		str4 := strings.TrimSpace(string([]byte(s)[len("#+"):]))
		str5 := strings.Split(str4, "##")
		if len(str5) != 2 {
			if strings.TrimSpace(str5[0]) == "" {
				replyText := "商品名不能为空"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			if Contains(groupName, "黄小姐") {
				err := IDBGAN("10001", "10001", str5[0])
				if err != nil {
					replyText := "删除失败"
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: replyText,
						},
						ReqType: utils.GroupMsg,
					}
				}
				replyText := "删除成功"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			if Contains(groupName, "奇乐") {
				err := IDBGAN("10002", "10002", str5[0])
				if err != nil {
					replyText := "删除失败"
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: replyText,
						},
						ReqType: utils.GroupMsg,
					}
				}
				replyText := "删除成功"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			err := IDBGAN(strconv.Itoa(int(groupId)), strconv.Itoa(int(groupId)), str5[0])
			if err != nil {
				replyText := "删除失败"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			replyText := "删除成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		if strings.TrimSpace(str5[0]) == "" {
			replyText := "商品名不能为空"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		str6 := strings.Split(str5[1], "#&")
		if len(str6) != 2 {
			if Contains(groupName, "黄小姐") {
				_, err := ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
				replyText := "添加成功"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			if Contains(groupName, "奇乐") {
				_, err := ItemSave("10002", "10002", null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
				replyText := "添加成功"
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GroupMsg,
				}
			}
			_, err := ItemSave(strconv.Itoa(int(groupId)), strconv.Itoa(int(groupId)), null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
			replyText := "添加成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		if Contains(groupName, "黄小姐") {
			_, err := ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
			replyText := "添加成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		if Contains(groupName, "奇乐") {
			_, err := ItemSave("10002", "10002", null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
			replyText := "添加成功"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		}
		_, err := ItemSave(strconv.Itoa(int(groupId)), strconv.Itoa(int(groupId)), null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
		replyText := "添加成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyText,
			},
			ReqType: utils.GroupMsg,
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
		fmt.Println(cps)
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
			replyText := "暂无相关记录"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		} else {
			psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + "\n品名 | 价格 | 备注" + ps + "\n价格源自 " + from
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, psc)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: psc,
				},
				ReqType: utils.GroupMsg,
			}
		}
	} else {
		from = strings.TrimSpace(reg5.ReplaceAllString(reg4.ReplaceAllString(strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), "黄小姐的魔方店"), "奇乐魔方坊"))
		cps, _ := GetItems(strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), s)
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
			replyText := "暂无相关记录"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupMsg,
			}
		} else {
			psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + "\n品名 | 价格 | 备注" + ps + "\n价格源自 " + from
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, psc)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: psc,
				},
				ReqType: utils.GroupMsg,
			}
		}
	}
}

func init() {
	utils.Register("查价", &PricePlugin{})
}
