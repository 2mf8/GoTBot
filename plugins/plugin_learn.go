package plugins

import (
	"context"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	//. "github.com/2mf8/go-tbot-for-rq/config"
	. "github.com/2mf8/go-tbot-for-rq/data"
	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
	"gopkg.in/guregu/null.v3"
)

type LearnPlugin struct {
}

func (learnPlugin *LearnPlugin) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	userId := event.Sender.UserId
	rawMsg := strings.TrimSpace(event.RawMessage)
	botId := bot.BotId

	rand.Seed(time.Now().UnixNano())
	success := rand.Intn(101)
	//delete := rand.Intn(101) + 200
	failure := rand.Intn(101) + 400

	s, b := Prefix(rawMsg, ".")
	if b == false {
		return MESSAGE_IGNORE
	}

	reg1 := regexp.MustCompile("＃")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, "#"))
	//if StartsWith(str1, "#+") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
	if StartsWith(str1, "#+") && IsBotAdmin(userId) {
		str2 := strings.TrimSpace(strings.TrimPrefix(str1, "#+"))
		str3 := strings.Split(str2, "##")
		if len(str3) != 2 {
			if strings.TrimSpace(str3[0]) == "" {
				replyText := strconv.Itoa(failure) + "（问指令不能为空）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			err := LDBGAA(groupId, str3[0])
			if err != nil {
				replyText := strconv.Itoa(failure) + "（问答删除失败）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			replyText := strconv.Itoa(success) + "（问答删除成功）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := strconv.Itoa(failure) + "（问指令不能为空）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		err := LearnSave(strings.TrimSpace(str3[0]), groupId, userId, null.NewString(str3[1], true), time.Now())
		if err != nil {
			replyText := strconv.Itoa(failure) + "（添加失败）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		replyText := strconv.Itoa(success) + "（学习已完成，下次触发有效）"
		replyMsg := pbbot.NewMsg().Text(replyText)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	if StartsWith(str1, "++") && IsBotAdmin(userId) {
		str2 := strings.TrimSpace(strings.TrimPrefix(str1, "++"))
		str3 := strings.Split(str2, "##")
		if len(str3) != 2 {
			if strings.TrimSpace(str3[0]) == "" {
				replyText := strconv.Itoa(failure) + "（系统问指令不能为空）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			err := LDBGAA(int64(9999999990), str3[0])
			if err != nil {
				replyText := strconv.Itoa(failure) + "（系统问答删除失败）"
				replyMsg := pbbot.NewMsg().Text(replyText)
				log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
				_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
				return MESSAGE_BLOCK
			}
			replyText := strconv.Itoa(success) + "（系统问答删除成功）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := strconv.Itoa(failure) + "（系统问指令不能为空）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		err := LearnSave(strings.TrimSpace(str3[0]), int64(9999999990), userId, null.NewString(str3[1], true), time.Now())
		if err != nil {
			replyText := strconv.Itoa(failure) + "（系统问答添加失败）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		replyText := strconv.Itoa(success) + "（系统问答学习已完成，下次触发有效）"
		replyMsg := pbbot.NewMsg().Text(replyText)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	if strings.TrimSpace(rawMsg) == "" {
		replyText := strconv.Itoa(failure) + "（指令不能为空）"
		replyMsg := pbbot.NewMsg().Text(replyText)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	learn_get, err := LearnGet(groupId, strings.TrimSpace(s))
	//log.Println(learn_get.LearnSync.Answer.String,"ceshil", err)
	if err != nil || learn_get.LearnSync.Answer.String == "" {
		sys_learn_get, _ := LearnGet(int64(9999999990), strings.TrimSpace(s))
		if sys_learn_get.LearnSync.Answer.String != "" {
			replyMsg := pbbot.NewMsg().Text(sys_learn_get.LearnSync.Answer.String)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, sys_learn_get.LearnSync.Answer.String)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
	}
	if learn_get.LearnSync.Answer.String != "" {
		replyMsg := pbbot.NewMsg().Text(learn_get.LearnSync.Answer.String)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, learn_get.LearnSync.Answer.String)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	return MESSAGE_IGNORE
}

func init() {
	Register("学习", &LearnPlugin{})
}
