package plugins

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	//. "github.com/2mf8/tbotGo/config"
	. "github.com/2mf8/tbotGo/data"
	. "github.com/2mf8/tbotGo/public"
	. "github.com/2mf8/tbotGo/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
)

type Sub struct {
}

func (sub *Sub) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	rawMsg := strings.TrimSpace(event.RawMessage)
	botId := bot.BotId
	userId := event.UserId

	rand.Seed(time.Now().UnixNano())
	success := rand.Intn(101)
	delete := rand.Intn(101) + 200
	//failure := rand.Intn(101) + 400

	s, b := Prefix(rawMsg, ".")
	if b == false {
		return MESSAGE_IGNORE
	}

	if StartsWith(s, "订阅") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "订阅"))
		r_groupId, _ := strconv.Atoi(s)
		_ = SubSave(groupId, int64(r_groupId), userId)
		reply := strconv.Itoa(success) + " （订阅成功）"
		replyMsg := pbbot.NewMsg().Text(reply)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	if StartsWith(s, "取消订阅") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		_ = SubDeleteByGroupId(groupId)
		reply := strconv.Itoa(delete) + " （取消订阅成功）"
		replyMsg := pbbot.NewMsg().Text(reply)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	return MESSAGE_IGNORE
}

func init() {
	Register("订阅", &Sub{})
}
