package plugins

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	//. "github.com/2mf8/tbotGo/config"
	. "github.com/2mf8/tbotGo/data"
	. "github.com/2mf8/tbotGo/public"
	. "github.com/2mf8/tbotGo/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
	//"gopkg.in/guregu/null.v3"
)

type Block struct{}

func (block *Block) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	groupId := event.GroupId
	userId := event.UserId
	rawMsg := event.RawMessage
	botId := bot.BotId

	rand.Seed(time.Now().UnixNano())
	//success := rand.Intn(101)
	//delete := rand.Intn(101) + 200
	failure := rand.Intn(101) + 400

	ispblock, err := PBlockGet(userId)
	//fmt.Println(ispblock)
	if err != nil {
		fmt.Println("[INFO] ", err)
	}
	if ispblock.PBlockSync.UserId == userId && !IsBotAdmin(userId) {
		return MESSAGE_BLOCK
	}

	s, b := Prefix(rawMsg, ".")
	if b == false {
		return MESSAGE_IGNORE
	}
	reg1 := regexp.MustCompile("<at qq=\"")
	reg2 := regexp.MustCompile("\"/>")
	reg3 := regexp.MustCompile("  ")

	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, ""))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, " "))

	for Contains(str2, "  ") {
		str2 = strings.TrimSpace(reg3.ReplaceAllString(str2, " "))
	}

	if StartsWith(s, "屏蔽+") && IsBotAdmin(userId) {
		pUserID, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(str2, "屏蔽+")))
		if err != nil {
			replyMsg := strconv.Itoa(failure) + "（用户不存在）"
			reply := pbbot.NewMsg().Text(replyMsg)
			bot.SendGroupMessage(groupId, reply, false)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return MESSAGE_BLOCK
		}
		err = PBlockSave(int64(pUserID), true, userId, time.Now())
		if err != nil {
			replyMsg := "屏蔽" + strconv.Itoa(int(pUserID)) + "失败"
			reply := pbbot.NewMsg().Text(replyMsg)
			bot.SendGroupMessage(groupId, reply, false)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return MESSAGE_BLOCK
		}
		replyMsg := "屏蔽" + strconv.Itoa(int(pUserID)) + "成功"
		reply := pbbot.NewMsg().Text(replyMsg)
		bot.SendGroupMessage(groupId, reply, false)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
		return MESSAGE_BLOCK
	}
	if StartsWith(s, "屏蔽-") && IsBotAdmin(userId) {
		pUserID, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(str2, "屏蔽-")))
		if err != nil {
			replyMsg := strconv.Itoa(failure) + "（用户不存在）"
			reply := pbbot.NewMsg().Text(replyMsg)
			bot.SendGroupMessage(groupId, reply, false)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return MESSAGE_BLOCK
		}
		err = PBlockSave(int64(pUserID), false, userId, time.Now())
		if err != nil {
			replyMsg := "解除屏蔽" + strconv.Itoa(int(pUserID)) + "失败"
			reply := pbbot.NewMsg().Text(replyMsg)
			bot.SendGroupMessage(groupId, reply, false)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
			return MESSAGE_BLOCK
		}
		replyMsg := "解除屏蔽" + strconv.Itoa(int(pUserID)) + "成功"
		reply := pbbot.NewMsg().Text(replyMsg)
		bot.SendGroupMessage(groupId, reply, false)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyMsg)
		return MESSAGE_BLOCK
	}
	return MESSAGE_IGNORE
}

func init() {
	Register("屏蔽", &Block{})
}
