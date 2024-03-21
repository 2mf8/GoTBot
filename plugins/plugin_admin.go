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

	. "github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type Admin struct {
}

/*
* botId 机器人Id
* groupId 群Id
* userId 用户Id
* messageId 消息Id
* rawMsg 群消息
* card At展示
* userRole 用户角色，是否是管理员
* botRole 机器人角色， 是否是管理员
* retval 返回值，用于判断是否处理下一个插件
* replyMsg 待发送消息
* rs 成功防屏蔽码
* rd 删除防屏蔽码
* rf 失败防屏蔽码
 */
func (admin *Admin) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	if botId.Common < 1 {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	if groupId.Common == 560820998 || groupId.Common == 189420325 || groupId.Common == 348591755 || groupId.Common == 481097523 || groupId.Common == 176211061 || groupId.Common == 138080634 {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	s, b := Prefix(rawMsg, ".")
	if !b || !botRole {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
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
	rand.Seed(time.Now().UnixNano())
	jin_duration := 60 + rand.Intn(28740)
	if s == "抽奖禁言" {
		if userRole {
			msg := "失败，您是群主或管理员"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: msg,
				},
				ReqType: utils.GroupBan,
				OfficalMsgId: messageId.Offical,
			}
		}
		msg := " 恭喜你抽中" + convertJinTime(jin_duration) + "禁言套餐，已发放"
		log.Printf("[抽奖禁言] Bot(%v) Group(%v) -> %v", botId, groupId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType:  utils.GroupBan,
			Duration: int64(jin_duration),
			BanId:    userId.Common,
			OfficalMsgId: messageId.Offical,
		}
	}
	if s == "退群" && super {
		return utils.RetStuct{
			RetVal:  utils.MESSAGE_BLOCK,
			ReqType: utils.GroupLeave,
		}
	}
	if StartsWith(str2, "jin") && (super || userRole) {
		str2 = strings.TrimSpace(string([]byte(str2)[len("jin"):]))
		str3 := strings.Split(str2, " ")

		if len(str3) != 2 {
			replyText := "禁言格式错误"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupBan,
				OfficalMsgId: messageId.Offical,
			}
		}
		jinId, err := strconv.ParseInt(str3[0], 10, 64)
		if err != nil {
			replyText := "禁言对象错误"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupBan,
				OfficalMsgId: messageId.Offical,
			}
		}

		duration := convertTime(str3[1])

		if duration <= 0 {
			replyText := "解除 " + strconv.Itoa(int(jinId)) + " 的禁言"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal:   utils.MESSAGE_BLOCK,
				ReqType:  utils.RelieveBan,
				Duration: int64(duration),
				BanId:    jinId,
			}
		}
		if duration < 30*60*60*24 {
			replyText := "禁言 " + strconv.Itoa(int(jinId)) + " " + strconv.Itoa(int(duration)) + "秒"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal:   utils.MESSAGE_BLOCK,
				ReqType:  utils.GroupBan,
				Duration: int64(duration),
				BanId:    jinId,
			}
		} else {
			replyText := "禁言时间超过最大允许范围"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupBan,
				OfficalMsgId: messageId.Offical,
			}
		}
	}
	if (StartsWith(str2, "t") || StartsWith(str2, "T")) && (super || userRole) {
		rejectAddAgain := StartsWith(str2, "T")
		str2 = strings.TrimSpace(string([]byte(strings.ToLower(str2))[len("t"):]))
		tId, err := strconv.ParseInt(str2, 10, 64)
		if err != nil {
			replyText := "踢出对象错误"
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GroupKick,
				OfficalMsgId: messageId.Offical,
			}
		}
		replyText := "踢出 " + strconv.Itoa(int(tId)) + " 成功"
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return utils.RetStuct{
			RetVal:         utils.MESSAGE_BLOCK,
			ReqType:        utils.GroupKick,
			BanId:          tId,
			RejectAddAgain: rejectAddAgain,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func convertTime(str string) int32 {
	var duration int = 0
	reg4 := regexp.MustCompile("天")
	reg5 := regexp.MustCompile("小时")
	reg6 := regexp.MustCompile("时")
	reg7 := regexp.MustCompile("分")
	reg8 := regexp.MustCompile("秒")
	str4 := strings.TrimSpace(reg4.ReplaceAllString(str, "d"))
	str4 = strings.TrimSpace(reg5.ReplaceAllString(str4, "h"))
	str4 = strings.TrimSpace(reg6.ReplaceAllString(str4, "h"))
	str4 = strings.TrimSpace(reg7.ReplaceAllString(str4, "m"))
	str4 = strings.TrimSpace(reg8.ReplaceAllString(str4, "s"))
	str4 = str4 + "s"
	reg9 := regexp.MustCompile(`([0-9]+)(d|h|m|s)`)
	m := reg9.FindAllString(str4, -1)
	for _, v := range m {
		if EndsWith(v, "d") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("d")]))
			duration += num * 60 * 60 * 24
		}
		if EndsWith(v, "h") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("h")]))
			duration += num * 60 * 60
		}
		if EndsWith(v, "m") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("m")]))
			duration += num * 60
		}
		if EndsWith(v, "s") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("s")]))
			duration += num
		}
	}
	return int32(duration)
}

func convertJinTime(i int) string {
	var timeString string
	day := i / 86400
	hour := i % 86400 / 3600
	min := i % 3600 / 60
	sec := i % 60
	if i >= 86400 {
		timeString = fmt.Sprintf("%v 天 %v 小时 %v 分钟", day, hour, min)
		return timeString
	}
	if i <= 3600 {
		timeString = fmt.Sprintf("%v 分钟 %v 秒钟", min, sec)
		return timeString
	}
	timeString = fmt.Sprintf("%v 小时 %v 分钟 %v 秒钟", hour, min, sec)
	return timeString
}

func init() {
	utils.Register("群管", &Admin{})
}
