package plugins

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	. "github.com/2mf8/tbotGo/config"
	. "github.com/2mf8/tbotGo/data"
	. "github.com/2mf8/tbotGo/public"
	. "github.com/2mf8/tbotGo/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
)

type BotSwitch struct {
}

func (botSwitch *BotSwitch) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	rawMsg := strings.TrimSpace(event.RawMessage)
	userId := event.Sender.UserId
	botId := bot.BotId

	rand.Seed(time.Now().UnixNano())
	success := rand.Intn(101)
	//delete := rand.Intn(101) + 200
	failure := rand.Intn(101) + 400

	s, b := Prefix(rawMsg, ".")
	if b == false {
		return MESSAGE_IGNORE
	}

	if StartsWith(s, "开启") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		s = strings.TrimPrefix(s, "开启")
		for _, j := range Conf.Plugins {
			if s == j {
				err := SwitchSave(groupId, j, time.Now(), false)
				if err != nil {
					reply := strconv.Itoa(failure) + " （开启失败）"
					replyMsg := pbbot.NewMsg().Text(reply)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				} else {
					reply := strconv.Itoa(success) + " （开启成功）"
					replyMsg := pbbot.NewMsg().Text(reply)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				}
			}
		}
		reply := strconv.Itoa(failure) + " （功能不存在）"
		replyMsg := pbbot.NewMsg().Text(reply)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}

	if StartsWith(s, "关闭") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		s = strings.TrimPrefix(s, "关闭")
		for _, j := range Conf.Plugins {
			if s == j {
				err := SwitchSave(groupId, j, time.Now(), true)
				if err != nil {
					reply := strconv.Itoa(failure) + " （关闭失败）"
					replyMsg := pbbot.NewMsg().Text(reply)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				} else {
					reply := strconv.Itoa(success) + " （关闭成功）"
					replyMsg := pbbot.NewMsg().Text(reply)
					log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
					_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
					return MESSAGE_BLOCK
				}
			}
		}
		reply := strconv.Itoa(failure) + " （功能不存在）"
		replyMsg := pbbot.NewMsg().Text(reply)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, reply)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	return MESSAGE_IGNORE
}

func init() {
	Register("开关", &BotSwitch{})
}
