package plugins

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"time"

	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
)

type Repeat struct {
}

func (rep *Repeat) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	rawMsg := strings.TrimSpace(event.RawMessage)
	botId := bot.BotId
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(101)

	if len(rawMsg) < 20 && r%70 == 0 && !(StartsWith(rawMsg, ".") || StartsWith(rawMsg, "%") || StartsWith(rawMsg, "％")) {
		replyMsg := pbbot.NewMsg().Text(rawMsg)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, rawMsg)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
	}
	return MESSAGE_IGNORE
}

func init() {
	Register("复读", &Repeat{})
}
