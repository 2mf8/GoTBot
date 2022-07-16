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
	. "github.com/2mf8/GoTBot/utils"
)

type ChannelBlock struct{}

func (block *ChannelBlock) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super bool, rs, rd, rf int) (retStuct RetChannelStuct){

	ispblock, err := PBlockGet(int64(userId))
	//fmt.Println(ispblock)
	if err != nil {
		fmt.Println("[INFO] ", err)
	}
	if ispblock.PBlockSync.UserId == int64(userId) && ispblock.PBlockSync.IsPBlock {
		if !super {
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
			}
		}
	}

	s, b := Prefix(rawMsg, ".")
	if !b {
		return RetChannelStuct{
			RetVal: MESSAGE_IGNORE,
		}
	}
	reg1 := regexp.MustCompile("<at qq=\"")
	reg2 := regexp.MustCompile("\"/>")
	reg3 := regexp.MustCompile("  ")

	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, ""))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, " "))

	for Contains(str2, "  ") {
		str2 = strings.TrimSpace(reg3.ReplaceAllString(str2, " "))
	}

	if StartsWith(s, "屏蔽+") && super {
		pUserID, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(str2, "屏蔽+")))
		if err != nil {
			replyMsg := strconv.Itoa(rf) + "（用户不存在）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyMsg)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyMsg,
				},
				ReqType: GroupMsg,
			}
		}
		err = PBlockSave(int64(pUserID), true, int64(userId), time.Now())
		if err != nil {
			replyMsg := "屏蔽" + strconv.Itoa(int(pUserID)) + "失败"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyMsg)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyMsg,
				},
				ReqType: GroupMsg,
			}
		}
		replyMsg := "屏蔽" + strconv.Itoa(int(pUserID)) + "成功"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyMsg)
		return RetChannelStuct{
			RetVal: MESSAGE_BLOCK,
			ReplyMsg: &Msg{
				Text: replyMsg,
			},
			ReqType: GroupMsg,
		}
	}
	if StartsWith(s, "屏蔽-") && super {
		pUserID, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(str2, "屏蔽-")))
		if err != nil {
			replyMsg := strconv.Itoa(rf) + "（用户不存在）"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyMsg)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyMsg,
				},
				ReqType: GroupMsg,
			}
		}
		err = PBlockSave(int64(pUserID), false, int64(userId), time.Now())
		if err != nil {
			replyMsg := "解除屏蔽" + strconv.Itoa(int(pUserID)) + "失败"
			log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyMsg)
			return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyMsg,
				},
				ReqType: GroupMsg,
			}
		}
		replyMsg := "解除屏蔽" + strconv.Itoa(int(pUserID)) + "成功"
		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) -> %v", botId, guildId, channelId, replyMsg)
		return RetChannelStuct{
				RetVal: MESSAGE_BLOCK,
				ReplyMsg: &Msg{
					Text: replyMsg,
				},
				ReqType: GroupMsg,
			}
	}
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}

func init() {
	ChannelRegister("屏蔽", &ChannelBlock{})
}
