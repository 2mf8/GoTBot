package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/2mf8/GoPbBot"
	"github.com/2mf8/GoPbBot/proto_gen/onebot"
	"github.com/2mf8/GoTBot/data"
	_ "github.com/2mf8/GoTBot/plugins"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron"
)

type Push struct {
	Bot     *pbbot.Bot
	GroupId int64
}

var push = Push{}

var pushes = make(map[int64]*Push)

func main() {

	color.Cyan("[INFO] 欢迎您使用GoTBot")

	_, err := os.Stat("conf.toml")
	if err != nil {
		_ = ioutil.WriteFile("conf.toml", []byte("Plugins = [\"守卫\",\"屏蔽\",\"开关\",\"复读\",\"回复\",\"群管\",\"订阅\",\"查价\",\"打乱\",\"学习\"]   #插件管理\nChannelPlugins = [\"守卫\",\"屏蔽\",\"开关\",\"复读\",\"回复\",\"订阅\",\"查价\",\"打乱\",\"学习\"]   #频道插件管理\nAdmins = [2693678434]   #机器人管理员管理\nDatabaseUser = \"sa\"   #MSSQL数据库用户名\nDatabasePassword = \"wr@#kequ5060\"   #MSSQL数据库密码\nDatabasePort = 1433   #MSSQL数据库服务端口\nDatabaseServer = \"127.0.0.1\"   #MSSQL数据库服务网址\nServerPort = 8081   #服务端口\nScrambleServer = \"http://localhost:2014\"   #打乱服务地址"), 0644)
	}

	plugin, _ := public.TbotConf()
	pluginString := fmt.Sprintf("%s", plugin.Conf)
	channelPluginString := fmt.Sprintf("%s", plugin.ChannelConf)
	fmt.Fprintf(color.Output, "%s %s", color.CyanString("[INFO] 已加载插件"), color.HiMagentaString(pluginString))
	fmt.Fprintf(color.Output, "%s %s", color.CyanString("\n[INFO] 已加载频道插件"), color.HiMagentaString(channelPluginString))

	pbbot.HandleConnect = func(bot *pbbot.Bot) {
		fmt.Printf("\n[连接] 新机器人已连接：%d\n", bot.BotId)
		fmt.Println("[已连接] 所有机器人列表：")
		for botId, _ := range pbbot.Bots {
			fmt.Println("[已连接]", botId)
		}
	}

	pbbot.HandleGroupRecallNotice = func(bot *pbbot.Bot, event *onebot.GroupRecallNoticeEvent) {
		groupId := event.GroupId
		msg_id := event.MessageId
		botId := bot.BotId
		log.Printf("[撤回消息] Bot(%v)  Group(%v)  -- MessageID(%v)", botId, groupId, msg_id)
	}

	pbbot.HandleGroupRequest = func(bot *pbbot.Bot, event *onebot.GroupRequestEvent) {
		//groupId := event.GroupId
		//userId := event.UserId
		invitor_uin, _ := strconv.Atoi(event.Extra["invitor_uin"])
		botId := bot.BotId
		if public.IsBotAdmin(int64(invitor_uin)) {
			bot.SetGroupAddRequest(event.Flag, event.SubType, true, "")
			log.Printf("[INFO] Bot(%v) Invitor(%v) -- 机器人加群 %v", botId, invitor_uin, true)
		}
		/*if IsAdmin(bot, groupId, botId) {
			bot.SetGroupAddRequest(event.Flag, event.SubType, false, "")
			log.Printf("[INFO] Bot(%v)  Group(%v)  -- %v 加群", botId, groupId, userId)
		}*/
	}

	pbbot.HandleFriendRequest = func(bot *pbbot.Bot, event *onebot.FriendRequestEvent) {
		bot.SetFriendAddRequest(event.Flag, true, "")
	}

	pbbot.HandleGroupIncreaseNotice = func(bot *pbbot.Bot, event *onebot.GroupIncreaseNoticeEvent) {
		groupId := event.GroupId
		userId := event.UserId
		botId := bot.BotId
		if userId == botId {
			msgPush := pbbot.NewMsg().Text("欢迎使用tbot")
			bot.SendGroupMessage(groupId, msgPush, false)
		}
		//msgPush := pbbot.NewMsg().At(userId, event.)
		//bot.SendGroupMessage(groupId, msgPush, false)
		log.Println(event)
	}

	rand.Seed(time.Now().UnixNano())
	second := rand.Intn(61)
	start := strconv.Itoa(second) + " 0 23 * * ?"
	end := strconv.Itoa(second) + " 59 7 * * ?"
	// 定时消息
	timer := cron.New()
	//cron表达式由6部分组成，从左到右分别表示 秒 分 时 日 月 星期
	//*表示任意值  ？表示不确定值，只能用于星期和日
	//这里表示每天22:32分发送消息
	timer.AddFunc(end, wholeBanRelieve)
	timer.AddFunc(start, wholeBan)
	timer.Start()

	pbbot.HandlePrivateMessage = func(bot *pbbot.Bot, event *onebot.PrivateMessageEvent) {
		if event.RawMessage == "poke" {
			poke := pbbot.NewMsg().Poke(event.UserId)
			bot.SendPrivateMessage(event.UserId, poke, false)
		}
		if event.RawMessage == "延时撤回" {
			resp, err := bot.SendPrivateMessage(event.UserId, pbbot.NewMsg().Text("撤回测试"), false)
			if err != nil {
				fmt.Println(err)
			}
			if resp != nil {
				time.Sleep(10 * time.Second)
			}
			bot.DeleteMsg(resp.MessageId)
		}
	}

	pbbot.HandleChannelMessage = func(bot *pbbot.Bot, event *onebot.ChannelMessageEvent) {
		rand.Seed(time.Now().UnixNano())
		guildId := event.GuildId
		channelId := event.ChannelId
		rawMsg := event.RawMessage
		botId := bot.BotId
		botChannelId := event.SelfId
		userId := event.Sender.TinyId
		card := event.Sender.Nickname
		userRole := public.IsGuildAdmin(event.Sender.RoleNames)
		super := public.IsBotAdmin(int64(event.Sender.TinyId))
		success := rand.Intn(101)
		delete := rand.Intn(101) + 200
		failure := rand.Intn(101) + 400

		log.Printf("[INFO] Bot(%v) GuildId(%v) ChannelId(%v) <- %v", botId, guildId, channelId, rawMsg)

		ctx := context.WithValue(context.Background(), "key", "value")
		sg, _ := data.SGBGI(int64(channelId))
		for _, i := range plugin.ChannelConf {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(data.PluginNameToIntent(i))
			if intent == int64(data.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := utils.ChannelPluginSet[i].ChannelDo(&ctx, botId, botChannelId, guildId, channelId, userId, rawMsg, card, super, userRole, success, delete, failure)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				log.Println(retStuct.ReplyMsg.Text)
				if retStuct.ReplyMsg != nil {
					newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
					if retStuct.ReplyMsg.Image != "" {
						newMsg = newMsg.Image(retStuct.ReplyMsg.Image)
					}
					log.Println(newMsg)
					bot.SendChannelMessage(guildId, channelId, newMsg, false)
				}
				break
			}
		}
	}

	pbbot.HandleGroupMessage = func(bot *pbbot.Bot, event *onebot.GroupMessageEvent) {
		groupId := event.GroupId
		rawMsg := event.RawMessage
		messageId := event.MessageId
		botId := bot.BotId
		userId := event.UserId
		card := event.Sender.Card
		userRole := public.IsAdmin(bot, groupId, userId)
		botRole := public.IsAdmin(bot, groupId, botId)
		super := public.IsBotAdmin(userId)
		rand.Seed(time.Now().UnixNano())
		success := rand.Intn(101)
		delete := rand.Intn(101) + 200
		failure := rand.Intn(101) + 400

		if public.IsBotAdmin(userId) && rawMsg == "撤回打卡" {
			//bot.SetGroupSignIn(groupId)
			reply := pbbot.NewMsg().Text("打卡成功")
			resp, err := bot.SendGroupMessage(groupId, reply, false)
			if err != nil {
				fmt.Println(err)
			}
			if resp != nil {
				time.Sleep(20 * time.Second)
			}
			bot.DeleteMsg(resp.MessageId)
		}
		if rawMsg == "poke me" && super {
			poke := pbbot.NewMsg().Poke(userId)
			bot.SendGroupMessage(groupId, poke, false)
		}

		if groupId == int64(758958532) {
			push = Push{
				Bot:     bot,
				GroupId: groupId,
			}
			if pushes[groupId] == nil {
				pushes[groupId] = &push
			} else {
				pushes[groupId].Bot = bot
			}
		}
		if pushes[int64(758958532)] != nil && pushes[int64(758958532)].Bot.BotId == botId {
			pushes[int64(758958532)].Bot = bot
		}

		log.Printf("[INFO] Bot(%v) Group(%v) <- %v", botId, groupId, rawMsg)
		ctx := context.WithValue(context.Background(), "key", "value")
		sg, _ := data.SGBGI(groupId)
		for _, i := range plugin.Conf {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(data.PluginNameToIntent(i))
			if intent == int64(data.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := utils.PluginSet[i].Do(&ctx, botId, groupId, userId, messageId, rawMsg, card, botRole, userRole, super, success, delete, failure)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					log.Println(retStuct.ReplyMsg.Text)
					if retStuct.ReplyMsg != nil {
						newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
						if retStuct.ReplyMsg.Image != "" {
							newMsg = newMsg.Image(retStuct.ReplyMsg.Image)
						}
						bot.SendGroupMessage(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.GroupBan {
					if retStuct.BanId == 0 {
						if retStuct.ReplyMsg != nil {
							newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
							bot.SendGroupMessage(groupId, newMsg, false)
							break
						}
						break
					} else {
						bot.SetGroupBan(groupId, retStuct.BanId, retStuct.Duration)
						if retStuct.ReplyMsg != nil {
							newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
							bot.SendGroupMessage(groupId, newMsg, false)
							break
						}
						break
					}
				}
				if retStuct.ReqType == utils.RelieveBan {
					if retStuct.BanId == 0 {
						break
					} else {
						bot.SetGroupBan(groupId, retStuct.BanId, retStuct.Duration)
						break
					}
				}
				if retStuct.ReqType == utils.GroupKick {
					bot.SetGroupKick(groupId, userId, retStuct.RejectAddAgain)
					if retStuct.ReplyMsg != nil {
						newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMessage(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.GroupSignIn {
					bot.SetGroupSignIn(groupId)
					if retStuct.ReplyMsg != nil {
						newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMessage(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.GroupLeave {
					bot.SetGroupLeave(groupId, false)
					if retStuct.ReplyMsg != nil {
						newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMessage(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.DeleteMsg {
					bot.DeleteMsg(retStuct.MessageId)
					if retStuct.ReplyMsg != nil {
						newMsg := pbbot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMessage(groupId, newMsg, false)
					}
					break
				}
			}
		}
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/ws/rq/", func(c *gin.Context) {
		if err := pbbot.UpgradeWebsocket(c.Writer, c.Request); err != nil {
			fmt.Println("[失败] 创建机器人失败")
		}
	})

	if err := router.Run(":8081"); err != nil {
		panic(err)
	}
	select {}
}

func activeMsgPush() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1001)
	sendMsg := strconv.Itoa(r) + " 你的夜晚太珍贵，我们不忍心占用\n\n为避免影响大家休息，每晚11点开启全员禁言，次日早晨8点解封[测试]"
	for _, i := range pushes {
		if i.Bot != nil && i.GroupId != 0 {
			reply := pbbot.NewMsg().Text(sendMsg)
			i.Bot.SendGroupMessage(i.GroupId, reply, false)
			log.Printf("[推送] Bot(%v) Group(%v) <- %v", i.Bot, i.GroupId, sendMsg)
		}
	}
}

func wholeBan() {
	for _, i := range pushes {
		if i.Bot != nil && i.GroupId != 0 {
			i.Bot.SetGroupWholeBan(i.GroupId, true)
			log.Printf("[推送] Bot(%v) Group(%v) <- 全员禁言", i.Bot, i.GroupId)
		}
	}
}

func wholeBanRelieve() {
	for _, i := range pushes {
		if i.Bot != nil && i.GroupId != 0 {
			i.Bot.SetGroupWholeBan(i.GroupId, false)
			log.Printf("[推送] Bot(%v) Group(%v) <- 解除全员禁言", i.Bot, i.GroupId)
		}
	}
}
