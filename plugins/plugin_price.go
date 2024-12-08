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
func (price *PricePlugin) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
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
	is, p := Prefix(str3, ".%")
	if !(b || p) {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	if p {
		s = is
	}
	if gid == "c2c" && !super {
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

	auser := fmt.Sprintf("%s_%s", groupId.Offical, userId.Offical)
	fmt.Println(auser, userId.Offical)
	fmt.Println(isAuth(auser, userId.Offical))
	if StartsWith(s, "#+") && (userRole || super || isAuth(auser, userId.Offical)) {
		str4 := strings.TrimSpace(string([]byte(s)[len("#+"):]))
		str5 := strings.Split(str4, "##")
		sub, err := SubscribeRead()
		fmt.Println("suberr", err)
		if len(str5) != 2 {
			if err != nil {
				err := IDBGAN(gid, gid, strings.TrimSpace(str5[0]))
				fmt.Println(err, str5[0])
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
			} else {
				err := IDBGAN(strings.TrimSpace(sub[gid]), strings.TrimSpace(sub[gid]), strings.TrimSpace(str5[0]))
				fmt.Println(sub[gid], gid, err, str5[0], sub[gid])
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
			if err != nil {
				_, err := ItemSave(gid, gid, null.String{}, strings.TrimSpace(str5[0]), null.NewString(str6[0], true), null.String{}, null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
			} else {
				_, err := ItemSave(sub[gid], sub[gid], null.String{}, strings.TrimSpace(str5[0]), null.NewString(str6[0], true), null.String{}, null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
				fmt.Println(err)
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
		}
		if err != nil {
			_, err := ItemSave(gid, gid, null.String{}, strings.TrimSpace(str5[0]), null.NewString(str6[0], true), null.NewString(str6[1], true), null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
		} else {
			_, err := ItemSave(sub[gid], sub[gid], null.String{}, strings.TrimSpace(str5[0]), null.NewString(str6[0], true), null.NewString(str6[1], true), null.NewString(uid, true), time.Now().Unix(), isMagnetism, null.NewString("", true))
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
					Text: replyText},
				ReqType: utils.GroupMsg,
			}
		}
	}
	cps := []CuberPrice{}
	ps := ""
	psc := ""
	ic := 0
	from := ""
	sub, err := SubscribeRead()
	if err != nil {
		from = gid
		index := 0
		ss := strings.Split(s, "#")
		if len(ss) > 1 {
			i, err := strconv.ParseInt(ss[1], 10, 64)
			if err == nil {
				index = int(i)
			}
		}
		cps, _ = GetItems(groupId.Offical, groupId.Offical, ss[0])
		for id, i := range cps {
			if index > 0 {
				if id < index-1 {
					continue
				} else {
					if i.Shipping.String == "" {
						ps += "\n" + i.Item + " | " + i.Price.String
					} else {
						ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
					}
					if ic == 19 {
						ps += "\n..." + "\n\n翻页请使用\n%[品名]#[序号]\n指令。例如：\n%三#21\n"
						break
					}
					ic++
				}
			} else {
				if i.Shipping.String == "" {
					ps += "\n" + i.Item + " | " + i.Price.String
				} else {
					ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
				}
				if ic == 19 {
					ps += "\n..." + "\n\n翻页请使用\n%[品名]#[序号]\n指令。例如：\n%三#21\n"
					break
				}
				ic++
			}
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
		from = strings.TrimSpace(reg5.ReplaceAllString(reg4.ReplaceAllString(sub[groupId.Offical], "黄小姐的魔方店"), "奇乐魔方坊"))
		if from == "" {
			from = gid
		}
		index := 0
		ss := strings.Split(s, "#")
		if len(ss) > 1 {
			i, err := strconv.ParseInt(ss[1], 10, 64)
			if err == nil {
				index = int(i)
			}
		}
		cps, _ := GetItems(sub[groupId.Offical], sub[groupId.Offical], ss[0])
		for id, i := range cps {
			if index > 0 {
				if id < index-1 {
					continue
				} else {
					if i.Shipping.String == "" {
						ps += "\n" + i.Item + " | " + i.Price.String
					} else {
						ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
					}
					if ic == 19 {
						ps += "\n..." + "\n\n翻页请使用\n%[品名]#[序号]\n指令。例如：\n%三#21\n"
						break
					}
					ic++
				}
			} else {
				if i.Shipping.String == "" {
					ps += "\n" + i.Item + " | " + i.Price.String
				} else {
					ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
				}
				if ic == 19 {
					ps += "\n..." + "\n\n翻页请使用\n%[品名]#[序号]\n指令。例如：\n%三#21\n"
					break
				}
				ic++
			}
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

func isAuth(k string, u string) bool {
	ia, err := AuthRead()
	if err != nil {
		return false
	}
	if ia[k] == u {
		return true
	}
	return false
}

func init() {
	utils.Register("查价", &PricePlugin{})
}
