package plugins

import (
	"context"
	"log"
	"net/url"
	"strings"

	. "github.com/2mf8/go-tbot-for-rq/data"
	. "github.com/2mf8/go-tbot-for-rq/public"
	. "github.com/2mf8/go-tbot-for-rq/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
)

type ScramblePlugin struct {
}

func (scramble *ScramblePlugin) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	rawMsg := strings.TrimSpace(event.RawMessage)
	botId := bot.BotId

	s, b := Prefix(rawMsg, ".")
	if b == false {
		return MESSAGE_IGNORE
	}

	ins := Tnoodle(s).Instruction
	shor := Tnoodle(s).ShortName
	show := Tnoodle(s).ShowName
	if ins == s && ins != "instruction" {
		gs := GetScramble(shor)
		if StartsWith(gs, "net") || gs == "获取失败" {
			replyMsg := pbbot.NewMsg().Text("获取打乱失败")
			log.Printf("[INFO] Bot(%v) Group(%v) -> 获取打乱失败", botId, groupId)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return MESSAGE_BLOCK
		}
		if shor == "minx" {
			gs = strings.Replace(gs, "U' ", "#\n", -1)
			gs = strings.Replace(gs, "U ", "U\n", -1)
			gs = strings.Replace(gs, "#", "U'", -1)
		}
		imgUrl := "http://localhost:2014/view/" + shor + ".png?scramble=" + url.QueryEscape(strings.Replace(gs, "\n", " ", -1))
		sc := show + "\n" + gs
		replyMsg := pbbot.NewMsg().Text(sc).Image(imgUrl)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v\n%v<image url=\"%v\"/>", botId, groupId, show, gs, imgUrl)
		_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
		return MESSAGE_BLOCK
	}
	return MESSAGE_IGNORE
}

func init() {
	Register("打乱", &ScramblePlugin{})
}
