package plugins

import (
	"context"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/2mf8/tbotGo/data"
	. "github.com/2mf8/tbotGo/public"
	. "github.com/2mf8/tbotGo/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
	"gopkg.in/guregu/null.v3"
)

type PricePlugin struct {
}

func (price *PricePlugin) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	rawMsg := strings.TrimSpace(event.RawMessage)
	botId := bot.BotId
	userId := event.UserId

	rand.Seed(time.Now().UnixNano())
	success := rand.Intn(101)
	//delete := rand.Intn(101) + 200
	failure := rand.Intn(101) + 400

	reg1 := regexp.MustCompile("％")
	reg2 := regexp.MustCompile("＃")
	reg3 := regexp.MustCompile("＆")
	reg4 := regexp.MustCompile("10001")
	reg5 := regexp.MustCompile("560820998")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(rawMsg, "%"))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, "#"))
	str3 := strings.TrimSpace(reg3.ReplaceAllString(str2, "&"))

	s, b := Prefix(str3, "%")
	if b == false {
		return MESSAGE_IGNORE
	}

	if StartsWith(s, "#+") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		str4 := strings.TrimSpace(string([]byte(s)[len("#+"):]))
		str5 := strings.Split(str4, "##")
		if len(str5) != 2 {
			if strings.TrimSpace(str5[0]) == "" {
				replyText := strconv.Itoa(failure) + "（商品名不能为空）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			if groupId == 481097523 || groupId == 176211061 || groupId == 138080634 {
				err := IDBGAN(int64(10001), str5[0])
				if err != nil {
					replyText := strconv.Itoa(failure) + "（删除失败）"
					replyMsg := pbbot.NewMsg().Text(replyText)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				}
				replyText := strconv.Itoa(success) + "（删除成功）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			if groupId == 560820998 || groupId == 189420325 || groupId == 348591755 {
				err := IDBGAN(int64(560820998), str5[0])
				if err != nil {
					replyText := strconv.Itoa(failure) + "（删除失败）"
					replyMsg := pbbot.NewMsg().Text(replyText)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				}
				replyText := strconv.Itoa(success) + "（删除成功）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			err := IDBGAN(groupId, str5[0])
			if err != nil {
				replyText := strconv.Itoa(failure) + "（删除失败）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			replyText := strconv.Itoa(success) + "（删除成功）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		if strings.TrimSpace(str5[0]) == "" {
			replyText := strconv.Itoa(failure) + "（商品名不能为空）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		str6 := strings.Split(str5[1], "#&")
		if len(str6) != 2 {
			if groupId == 481097523 || groupId == 176211061 || groupId == 138080634 {
				err := ItemSave(int64(10001), null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, userId, null.NewTime(time.Now(), true))
				if err != nil {
					replyText := strconv.Itoa(failure) + "（添加失败）"
					replyMsg := pbbot.NewMsg().Text(replyText)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				}
				replyText := strconv.Itoa(success) + "（添加成功）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			if groupId == 560820998 || groupId == 189420325 || groupId == 348591755 {
				err := ItemSave(int64(560820998), null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, userId, null.NewTime(time.Now(), true))
				if err != nil {
					replyText := strconv.Itoa(failure) + "（添加失败）"
					replyMsg := pbbot.NewMsg().Text(replyText)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				}
				replyText := strconv.Itoa(success) + "（添加成功）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			err := ItemSave(groupId, null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, userId, null.NewTime(time.Now(), true))
			if err != nil {
				replyText := strconv.Itoa(failure) + "（添加失败）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			replyText := strconv.Itoa(success) + "（添加成功）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		if groupId == 481097523 || groupId == 176211061 || groupId == 138080634 {
			err := ItemSave(int64(10001), null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), userId, null.NewTime(time.Now(), true))
			if err != nil {
				replyText := strconv.Itoa(failure) + "（添加失败）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			replyText := strconv.Itoa(success) + "（添加成功）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		if groupId == 560820998 || groupId == 189420325 || groupId == 348591755 {
			err := ItemSave(int64(560820998), null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), userId, null.NewTime(time.Now(), true))
			if err != nil {
				replyText := strconv.Itoa(failure) + "（添加失败）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			replyText := strconv.Itoa(success) + "（添加成功）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		err := ItemSave(groupId, null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), userId, null.NewTime(time.Now(), true))
		if err != nil {
			replyText := strconv.Itoa(failure) + "（添加失败）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		replyText := strconv.Itoa(success) + "（添加成功）"
		replyMsg := pbbot.NewMsg().Text(replyText)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	cps := []CuberPrice{}
	ps := ""
	psc := ""
	ic := 0
	from := ""
	sub, err := GetSubscribe(groupId)
	if err != nil {
		from = strconv.Itoa(int(groupId))
		cps, _ = GetItems(groupId, s)
	} else {
		from = strings.TrimSpace(reg5.ReplaceAllString(reg4.ReplaceAllString(strconv.Itoa(int(sub.SubSync.ReplaceGroupId)), "黄小姐的魔方店"), "奇乐魔方坊"))
		cps, _ = GetItems(sub.SubSync.ReplaceGroupId, s)
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
		replyText := strconv.Itoa(success) + "（暂无相关记录）"
		replyMsg := pbbot.NewMsg().Text(replyText)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	} else {
		psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + "\n品名 | 价格 | 备注" + ps + "\n价格源自 " + from
		replyMsg := pbbot.NewMsg().Text(psc)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, psc)
		bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
}

func init() {
	Register("魔友价", &PricePlugin{})
}
