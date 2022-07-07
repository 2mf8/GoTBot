package plugins

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
	. "github.com/2mf8/go-tbot-for-rq/data"
	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
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

	if !IsAdmin(bot, groupId, botId) {
		return MESSAGE_IGNORE
	}

	ggg, _ := GetJudgeGroup()

	if StartsWith(rawMsg, ".守卫") && IsBotAdmin(userId) {
		vocabulary := strings.TrimPrefix(rawMsg, ".守卫")
		content := strings.Split(vocabulary, " ")
		err := ggg.JudgeGroupUpdate(ArrayStringToArrayInt64(content)...)
		if err != nil {
			log.Panicln(err)
		}
		msg := strconv.Itoa(r) + " （守卫群添加成功）"
		replyMsg := pbbot.NewMsg().Text(msg)
		bot.SendGroupMessage(groupId, replyMsg, false)
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return MESSAGE_BLOCK
	}

	if StartsWith(rawMsg, ".取消守卫") && IsBotAdmin(userId) {
		vocabulary := strings.TrimPrefix(rawMsg, ".取消守卫")
		content := strings.Split(vocabulary, " ")
		ggg.JudgeGroupDelete(ArrayStringToArrayInt64(content)...)
		msg := strconv.Itoa(delete) + " （守卫群删除成功）"
		replyMsg := pbbot.NewMsg().Text(msg)
		bot.SendGroupMessage(groupId, replyMsg, false)
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return MESSAGE_BLOCK
	}

	containsGroupId := JudgeGroupId(groupId, *ggg.JudgeGroupSync)
	fmt.Println(containsGroupId)
	if containsGroupId == 0 {
		return MESSAGE_IGNORE
	}

	ggk, _ := GetJudgeKeys()

	if StartsWith(rawMsg, ".拦截") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		vocabulary := strings.TrimPrefix(rawMsg, ".拦截")
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

	if StartsWith(rawMsg, ".取消拦截") && IsBotAdmin(userId) {
		vocabulary := strings.TrimPrefix(rawMsg, ".取消拦截")
		content := strings.Split(vocabulary, " ")
		ggk.JudgeKeysDelete(content...)
		msg := strconv.Itoa(delete) + " （拦截词汇删除成功）"
		replyMsg := pbbot.NewMsg().Text(msg)
		bot.SendGroupMessage(groupId, replyMsg, false)
		log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return MESSAGE_BLOCK
	}

	containsJudgeKeys := Judge(rawMsg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		if IsAdmin(bot, groupId, userId) {
			msg := strconv.Itoa(r) + " （消息触发守卫，已被拦截）"
			replyMsg := pbbot.NewMsg().Text(msg)
			bot.SendGroupMessage(groupId, replyMsg, false)
			log.Printf("[守卫] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
			return MESSAGE_BLOCK
		}
		bot.DeleteMsg(messageId)
		bot.SetGroupBan(groupId, userId, int32(120))
		msg := strconv.Itoa(r) + " （消息触发守卫，已撤回消息并禁言该用户两分钟, 请文明发言）"
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
