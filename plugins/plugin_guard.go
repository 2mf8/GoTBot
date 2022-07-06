package plugins

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
	. "github.com/2mf8/go-tbot-for-rq/data"
)

type Guard struct {
}

func (guard *Guard) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	rawMsg := strings.TrimSpace(event.RawMessage)
	botId := bot.BotId
	userId := event.UserId
	messageId := event.MessageId
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(101)
	delete := rand.Intn(101) + 200

	s, b := Prefix(rawMsg, ".")
	if !b || !IsAdmin(bot, groupId, botId) {
		return MESSAGE_IGNORE
	}
	ggk, _ := GetJudgeKeys()

	if StartsWith(s, "拦截") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		vocabulary := strings.TrimPrefix(s, "拦截")
		content := strings.Split(vocabulary, " ")
		err := ggk.JudgeKeysUpdate(content...)
		if err != nil {
			log.Panicln(err)
		}
		msg := strconv.Itoa(r) + " （拦截词汇添加成功）"
		replyMsg := pbbot.NewMsg().Text(msg)
		bot.SendGroupMessage(groupId, replyMsg, false)
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return MESSAGE_BLOCK
	}

	if StartsWith(s, "取消拦截") && IsBotAdmin(userId) {
		vocabulary := strings.TrimPrefix(s, "取消拦截")
		content := strings.Split(vocabulary, " ")
		ggk.JudgeKeysDelete(content...)
		msg := strconv.Itoa(delete) + " （拦截词汇删除成功）"
		replyMsg := pbbot.NewMsg().Text(msg)
		bot.SendGroupMessage(groupId, replyMsg, false)
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return MESSAGE_BLOCK
	}

	containsJudgeKeys := Judge(s, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		if IsAdmin(bot, groupId, userId) {
			msg := strconv.Itoa(r) + " （消息触发守卫，已被拦截）"
			replyMsg := pbbot.NewMsg().Text(msg)
			bot.SendGroupMessage(groupId, replyMsg, false)
			log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
			return MESSAGE_BLOCK
		}
		bot.DeleteMsg(messageId)
		bot.SetGroupBan(groupId, userId, int32(60))
		msg := strconv.Itoa(r) + " （消息触发守卫，已撤回消息并禁言该用户一分钟, 请文明发言）"
		replyMsg := pbbot.NewMsg().Text(msg)
		bot.SendGroupMessage(groupId, replyMsg, false)
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return MESSAGE_BLOCK
	}
	return MESSAGE_IGNORE
}

func init() {
	Register("守卫", &Guard{})
}
