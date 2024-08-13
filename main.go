package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	bot "github.com/2mf8/Go-QQ-Client"
	"github.com/2mf8/Go-QQ-Client/dto"
	"github.com/2mf8/Go-QQ-Client/event"
	"github.com/2mf8/Go-QQ-Client/token"
	"github.com/2mf8/Go-QQ-Client/websocket"
	database "github.com/2mf8/GoTBot/data"
	_ "github.com/2mf8/GoTBot/plugins"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
	gonebot "github.com/2mf8/GoneBot"
	"github.com/2mf8/GoneBot/onebot"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var msgs = make(map[int64]string, 0)

func main() {
	InitLog()
	tomlData := `
	Plugins = ["Log","守卫","开关","Bind","复读","WCA","回复","赛季","查价","打乱","学习","Rank"]   # 插件管理
	AppId = 0 # 机器人AppId
	AccessToken = "" # 机器人AccessToken
	ClientSecret = "" # 机器人ClientSecret
	Admins = [""]   # 机器人管理员管理
	DatabaseUser = "sa"   # MSSQL数据库用户名
	DatabasePassword = ""   # MSSQL数据库密码
	DatabasePort = 1433   # MSSQL数据库服务端口
	DatabaseServer = "127.0.0.1"   # MSSQL数据库服务网址
	DatabaseName = ""  # 数据库名
	ServerPort = 8081   # 服务端口
	ScrambleServer = "http://localhost:2014"   # 打乱服务地址
	RedisServer = "127.0.0.1" # Redis服务网址
	RedisPort = 6379 # Redis端口
	RedisPassword = "" # Redis密码
	RedisTable = 0 # Redis数据表
	RedisPoolSize = 1000 # Redis连接池数量
	JwtKey = ""
	RefreshKey = ""
	`

	log.Infoln("欢迎您使用GoTBot")
	_, err := os.Stat("conf.toml")
	if err != nil {
		_ = os.WriteFile("conf.toml", []byte(tomlData), 0644)
		log.Warn("已生成配置文件 conf.toml ,请修改后重新启动程序。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	allconfig := database.AllConfig
	log.Info("[配置信息]", allconfig)
	pluginString := fmt.Sprintf("%s", allconfig.Plugins)

	log.Infof("已加载插件 %s", pluginString)

	go StartOffical()

	gonebot.HandleConnect = func(bot *gonebot.Bot) {
		log.Infof("\n[连接] 新机器人已连接：%d\n", bot.BotId)
		log.Info("[已连接] 所有机器人列表：")
		for botId, _ := range gonebot.Bots {
			log.Info("[已连接]", botId)
		}
	}

	gonebot.HandlePrivateMessage = func(bot *gonebot.Bot, event *onebot.PrivateMsgEvent) {
		userId := event.Sender.UserId
		rawMsg := event.RawMessage
		if rawMsg == "hello" {
			msg := gonebot.NewMsg().Text("world")
			bot.SendPrivateMsg(userId, msg, false)
		}
	}
	gonebot.HandleGroupMessage = func(bot *gonebot.Bot, ievent *onebot.GroupMsgEvent) {
		sub, _ := database.SubscribeRead()
		groupId := ievent.GroupId
		rawMsg := ievent.RawMessage
		messageId := ievent.MessageId
		botId := bot.BotId
		userId := ievent.UserId
		card := ievent.Sender.Nickname
		gid := fmt.Sprintf("%v", groupId)
		uid := fmt.Sprintf("%v", userId)
		super := public.IsBotAdmin(uid, allconfig.Admins)
		rand.New(rand.NewSource(time.Now().UnixNano()))
		userRole := public.IsAdmin(ievent.Sender.Role)
		gi, _ := bot.GetGroupInfo(groupId, true)
		gmi, err := bot.GetGroupMemberInfo(groupId, bot.BotId, true)
		fmt.Println(gmi.Data.Role, err)
		botIsAdmin := public.IsAdmin(gmi.Data.Role)

		regStr := fmt.Sprintf(`\[CQ:at,qq=%v\]`, bot.BotId)
		reg := regexp.MustCompile(regStr)
		reg1 := regexp.MustCompile(`\[CQ:reply,id=[0-9]+\]`)

		ss := reg.FindAllString(rawMsg, -1)
		s1 := reg1.FindAllString(rawMsg, -1)
		ns := ""

		if super && rawMsg == "订阅" {
			database.SubscribeCreate(fmt.Sprintf("%v", groupId), "offical")
		}
		if super && rawMsg == "取消订阅" {
			database.SubscribeDelete(fmt.Sprintf("%v", groupId))
		}

		l, b := public.Prefix(rawMsg, "++")
		if b && super {
			ctx := context.WithValue(context.Background(), "key", "value")
			botType := utils.BotIdType{
				Common:  botId,
				Offical: "",
			}
			groupIdType := utils.GroupIdType{
				Common:  groupId,
				Offical: "",
			}
			userIdType := utils.UserIdType{
				Common:  userId,
				Offical: "",
			}
			msgIdType := utils.MsgIdType{
				Common:  messageId,
				Offical: "",
			}
			ns := fmt.Sprintf(".++%s", l)
			retStuct := utils.PluginSet["学习"].Do(&ctx, &botType, &groupIdType, &userIdType, gi.Data.GroupName, &msgIdType, ns, card, botIsAdmin, userRole, super)
			if retStuct.ReplyMsg != nil {
				newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
				if retStuct.ReplyMsg.Image != "" {
					newMsg = newMsg.Image(retStuct.ReplyMsg.Image)
				}
				bot.SendGroupMsg(groupId, newMsg, false)
			}
		}

		if len(ss) == 0 {
			ns = rawMsg
		} else {
			if len(s1) > 0 {
				ns = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(rawMsg, ss[0], "."), " ", ""), s1[0], ""), "..", ".")
			} else {
				ns = strings.ReplaceAll(strings.ReplaceAll(rawMsg, ss[0], "."), " ", "")
			}
		}

		if ns == "禁言" && botIsAdmin && (super || userRole) {
			bot.SetGroupWholeBan(groupId, true)
		}

		if ns == "解除" && botIsAdmin && (super || userRole) {
			bot.SetGroupWholeBan(groupId, false)
		}

		ctx := context.WithValue(context.Background(), "key", "value")
		sg, _ := database.SGBGIACI(gid, gid)
		for _, i := range allconfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			botType := utils.BotIdType{
				Common:  botId,
				Offical: "",
			}
			groupIdType := utils.GroupIdType{
				Common:  groupId,
				Offical: "",
			}
			userIdType := utils.UserIdType{
				Common:  userId,
				Offical: "",
			}
			msgIdType := utils.MsgIdType{
				Common:  messageId,
				Offical: "",
			}
			retStuct := utils.PluginSet[i].Do(&ctx, &botType, &groupIdType, &userIdType, gi.Data.GroupName, &msgIdType, ns, card, botIsAdmin, userRole, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						fmt.Println(retStuct.ReplyMsg.Text)
						if !strings.Contains(retStuct.ReplyMsg.Text, "http") && sub[fmt.Sprintf("%v", groupId)] == "offical" {
							_, b := public.Prefix(ievent.RawMessage, ".")
							if b {
								msgs[ievent.MessageId] = ievent.RawMessage
								bot.SendGroupBotCallback(102070767, ievent.GroupId)
							}
							_, p := public.Prefix(ievent.RawMessage, "%")
							if p {
								msgs[ievent.MessageId] = ievent.RawMessage
								bot.SendGroupBotCallback(102070767, ievent.GroupId)
							}
							break
						} else {
							if sub[fmt.Sprintf("%v", groupId)] == "offical" {
								newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
								if retStuct.ReplyMsg.Image != "" {
									newMsg = newMsg.Image(retStuct.ReplyMsg.Image)
								}
								bot.SendGroupMsg(groupId, newMsg, false)
							}
						}
					}
					break
				}
				if retStuct.ReqType == utils.GroupBan {
					if retStuct.BanId == 0 {
						if retStuct.ReplyMsg != nil {
							newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
							bot.SendGroupMsg(groupId, newMsg, false)
							break
						}
						break
					} else {
						bot.SetGroupBan(groupId, retStuct.BanId, retStuct.Duration)
						if retStuct.ReplyMsg != nil {
							newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
							bot.SendGroupMsg(groupId, newMsg, false)
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
					if retStuct.BanId == 0 {
						break
					}
					bot.SetGroupKick(groupId, retStuct.BanId, retStuct.RejectAddAgain)
					if retStuct.ReplyMsg != nil {
						newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMsg(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.GroupLeave {
					bot.SetGroupLeave(groupId, false)
					if retStuct.ReplyMsg != nil {
						newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMsg(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.DeleteMsg {
					bot.DeleteMsg(retStuct.MsgId)
					if retStuct.ReplyMsg != nil {
						newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMsg(groupId, newMsg, false)
					}
					break
				}
			}
		}
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/onebot/v11/ws", func(c *gin.Context) {
		if err := gonebot.UpgradeWebsocket(c.Writer, c.Request); err != nil {
			log.Info("[失败] 创建机器人失败")
		}
	})

	if err := router.Run("127.0.0.1:8080"); err != nil {
		panic(err)
	}
	select {}
}

func InitLog() {
	// 输出到命令行
	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	}
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)

	// 输出到文件
	rl, err := rotatelogs.New(path.Join("logs", "%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join("logs", "latest.log")), // 最新日志软链接
		rotatelogs.WithRotationTime(time.Hour*24),                // 每天一个新文件
		rotatelogs.WithMaxAge(time.Hour*24*3),                    // 日志保留3天
	)
	if err != nil {
		utils.FatalError(err)
		return
	}
	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  rl,
			log.WarnLevel:  rl,
			log.ErrorLevel: rl,
			log.FatalLevel: rl,
			log.PanicLevel: rl,
		},
		&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%time%] [%lvl%]: %msg% \r\n",
		},
	))
}

func StartOffical() {
	botLoginInfo := &public.BotLogin{
		AppId:       database.AllConfig.AppId,
		AccessToken: database.AllConfig.AccessToken,
	}
	token := token.BotToken(botLoginInfo.AppId, botLoginInfo.AccessToken, string(token.TypeBot))
	api := bot.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Warn("登录失败，请检查 appid 和 AccessToken 是否正确。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		log.Warn("登录失败，请检查 appid 和 AccessToken 是否正确。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
	}
	var guildMsg event.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
		guildId := data.GuildID
		channelId := data.ChannelID
		userId := data.Author.ID
		content := strings.TrimSpace(data.Content)
		msgId := data.ID
		reg4 := regexp.MustCompile("/")
		content = strings.TrimSpace(reg4.ReplaceAllString(content, ""))
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)
		ctx := context.WithValue(context.Background(), "key", "value")
		me, _ := api.Me(ctx)
		reg7 := regexp.MustCompile(fmt.Sprintf("<@!%s> ", me.ID))
		content = reg7.ReplaceAllString(content, ".")
		sg, _ := database.SGBGIACI(guildId, channelId)
		botType := utils.BotIdType{
			Common:  0,
			Offical: "",
		}
		groupIdType := utils.GroupIdType{
			Common:  0,
			Offical: channelId,
		}
		userIdType := utils.UserIdType{
			Common:  0,
			Offical: userId,
		}
		msgIdType := utils.MsgIdType{
			Common:  0,
			Offical: msgId,
		}
		for _, i := range database.AllConfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := utils.PluginSet[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &msgIdType, content, "", true, false, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						msg := fmt.Sprintf("%s", strings.TrimSpace(retStuct.ReplyMsg.Text))
						if retStuct.ReplyMsg.Image != "" {
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								Image:   retStuct.ReplyMsg.Image,
								MsgID:   data.ID,
							})
						} else {
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								MsgID:   data.ID,
							})
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								Image:   retStuct.ReplyMsg.Images[1],
								MsgID:   data.ID,
							})
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								Image:   retStuct.ReplyMsg.Images[1],
								MsgID:   data.ID,
							})
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								Image:   retStuct.ReplyMsg.Images[2],
								MsgID:   data.ID,
							})
						}
					}
					break
				}
			}
		}
		return nil
	}
	var recall event.InteractionEventHandler = func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		/* {
		"op": 0,
		"s": 2,
		"t": "INTERACTION_CREATE",
		"id": "INTERACTION_CREATE:5fd23877-f35a-42c0-bd7c-d39d22f459a5",
		-"d": {
		"application_id": "101981675",
		"chat_type": 1,
		-"data": {
		"resolved": { },
		"type": 11
		},
		"group_member_openid": "1973D8F78F51DE1E4C8ED4E54E1FB2F8",
		"group_openid": "2622F289E5391B88684D0C46AABBBC40",
		"id": "c899f03f-7cae-431b-9982-c4fc8b28de21",
		"scene": "group",
		"timestamp": "2024-08-11T14:02:51+08:00",
		"type": 11,
		"version": 1
		}
		} */
		b, _ := json.Marshal(data)
		fmt.Println("\n\n\n", string(b))
		sg, _ := database.SGBGIACI(data.GroupOpenID, data.GroupOpenID)
		botType := utils.BotIdType{
			Common:  0,
			Offical: "",
		}
		groupIdType := utils.GroupIdType{
			Common:  0,
			Offical: data.GroupOpenID,
		}
		userIdType := utils.UserIdType{
			Common:  0,
			Offical: data.GroupMemberOpenID,
		}
		for i, v := range msgs {
			defer delete(msgs, i)
			for _, i := range database.AllConfig.Plugins {
				intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
				if intent == int64(database.PluginReply) {
					break
				}
				if intent > 0 {
					continue
				}
				retStuct := utils.PluginSet[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &utils.MsgIdType{}, v, "", true, false, false)
				if retStuct.RetVal == utils.MESSAGE_BLOCK {
					if retStuct.ReqType == utils.GroupMsg {
						if retStuct.ReplyMsg != nil {
							msg := strings.TrimSpace(retStuct.ReplyMsg.Text)
							if retStuct.ReplyMsg.Image != "" {
								resp, _ := api.PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Image, SrvSendMsg: false})
								if resp != nil {
									newMsg := &dto.GroupMessageToCreate{
										Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
										Media: &dto.FileInfo{
											FileInfo: resp.FileInfo,
										},
										EventID: data.ID,
										MsgType: 7,
										MsgReq:  1,
									}
									api.PostGroupMessage(ctx, data.GroupOpenID, newMsg)
								}
							} else {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									EventID: data.ID,
									MsgType: 0,
								}
								api.PostGroupMessage(ctx, data.GroupOpenID, newMsg)
							}
							if len(retStuct.ReplyMsg.Images) == 2 {
								resp, _ := api.PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
								if resp != nil {
									newMsg := &dto.GroupMessageToCreate{
										Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
										Media: &dto.FileInfo{
											FileInfo: resp.FileInfo,
										},
										EventID: data.ID,
										MsgType: 7,
										MsgReq:  2,
									}
									api.PostGroupMessage(ctx, data.GroupOpenID, newMsg)
								}
							}
							if len(retStuct.ReplyMsg.Images) >= 3 {
								resp, _ := api.PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
								if resp != nil {
									newMsg := &dto.GroupMessageToCreate{
										Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
										Media: &dto.FileInfo{
											FileInfo: resp.FileInfo,
										},
										EventID: data.ID,
										MsgType: 7,
										MsgReq:  2,
									}
									api.PostGroupMessage(ctx, data.GroupOpenID, newMsg)
								}
								resp1, _ := api.PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[2], SrvSendMsg: false})
								if resp1 != nil {
									newMsg := &dto.GroupMessageToCreate{
										Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
										Media: &dto.FileInfo{
											FileInfo: resp1.FileInfo,
										},
										EventID: data.ID,
										MsgType: 7,
										MsgReq:  3,
									}
									api.PostGroupMessage(ctx, data.GroupOpenID, newMsg)
								}
							}
						}
						break
					}
				}
			}
		}
		fmt.Println(err)
		return nil
	}
	var groupMessage event.GroupAtMessageEventHandler = func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		groupId := data.GroupId
		userId := data.Author.UserId
		content := strings.TrimSpace(data.Content)
		msgId := data.MsgId
		reg4 := regexp.MustCompile("/")
		content = strings.TrimSpace(reg4.ReplaceAllString(content, ""))
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)
		content = fmt.Sprintf(".%s", content)
		ctx := context.WithValue(context.Background(), "key", "value")
		sg, _ := database.SGBGIACI(groupId, groupId)
		botType := utils.BotIdType{
			Common:  0,
			Offical: "",
		}
		groupIdType := utils.GroupIdType{
			Common:  0,
			Offical: groupId,
		}
		userIdType := utils.UserIdType{
			Common:  0,
			Offical: userId,
		}
		msgIdType := utils.MsgIdType{
			Common:  0,
			Offical: msgId,
		}
		for _, i := range database.AllConfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := utils.PluginSet[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &msgIdType, content, "", true, false, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						msg := fmt.Sprintf("\n%s", strings.TrimSpace(retStuct.ReplyMsg.Text))
						if retStuct.ReplyMsg.Image != "" {
							resp, _ := api.PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Image, SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									MsgID:   data.MsgId,
									MsgType: 7,
									MsgReq:  1,
								}
								api.PostGroupMessage(ctx, groupId, newMsg)
							}
						} else {
							newMsg := &dto.GroupMessageToCreate{
								Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
								MsgID:   data.MsgId,
								MsgType: 0,
							}
							api.PostGroupMessage(ctx, groupId, newMsg)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							resp, _ := api.PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									MsgID:   data.MsgId,
									MsgType: 7,
									MsgReq:  2,
								}
								api.PostGroupMessage(ctx, groupId, newMsg)
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							resp, _ := api.PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									MsgID:   data.MsgId,
									MsgType: 7,
									MsgReq:  2,
								}
								api.PostGroupMessage(ctx, groupId, newMsg)
							}
							resp1, _ := api.PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[2], SrvSendMsg: false})
							if resp1 != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp1.FileInfo,
									},
									MsgID:   data.MsgId,
									MsgType: 7,
									MsgReq:  3,
								}
								api.PostGroupMessage(ctx, groupId, newMsg)
							}
						}
					}
					break
				}
			}
		}
		return nil
	}
	var c2cMessage event.C2CMessageEventHandler = func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		userId := data.Author.UserOpenId
		msgId := data.Id
		content := data.Content
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)
		sg, _ := database.SGBGIACI("c2c", "c2c")

		botType := utils.BotIdType{
			Common:  0,
			Offical: "",
		}
		groupIdType := utils.GroupIdType{
			Common:  0,
			Offical: "c2c",
		}
		userIdType := utils.UserIdType{
			Common:  0,
			Offical: userId,
		}
		msgIdType := utils.MsgIdType{
			Common:  0,
			Offical: msgId,
		}

		for _, i := range database.AllConfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := utils.PluginSet[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &msgIdType, content, "", true, false, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						if retStuct.ReplyMsg.Image != "" {
							resp, err := api.PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Image, SrvSendMsg: false})
							log.Info(err)
							if resp != nil {
								newMsg := &dto.C2CMessageToCreate{
									Content: strings.TrimSpace(retStuct.ReplyMsg.Text),
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									MsgID:   data.Id,
									MsgType: 7,
									MsgReq:  1,
								}
								_, err := api.PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
						} else {
							newMsg := &dto.C2CMessageToCreate{
								Content: strings.TrimSpace(retStuct.ReplyMsg.Text),
								MsgID:   data.Id,
								MsgType: 0,
								MsgReq:  1,
							}
							_, err := api.PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
							log.Info(err)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							resp, err := api.PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
							log.Info(err)
							if resp != nil {
								newMsg := &dto.C2CMessageToCreate{
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									MsgID:   data.Id,
									MsgType: 7,
									MsgReq:  1,
								}
								_, err := api.PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							resp, err := api.PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
							log.Info(err)
							if resp != nil {
								newMsg := &dto.C2CMessageToCreate{
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									MsgID:   data.Id,
									MsgType: 7,
									MsgReq:  1,
								}
								_, err := api.PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
							resp1, err := api.PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[2], SrvSendMsg: false})
							log.Info(err)
							if resp1 != nil {
								newMsg := &dto.C2CMessageToCreate{
									Media: &dto.FileInfo{
										FileInfo: resp1.FileInfo,
									},
									MsgID:   data.Id,
									MsgType: 7,
									MsgReq:  1,
								}
								_, err := api.PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
						}
					}
					break
				}
			}
		}

		if data.Content == "测试" {
			resp, err := api.PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: "https://www.2mf8.cn/static/image/cube3/b1.png", SrvSendMsg: false})
			log.Info(err)
			if resp != nil {
				newMsg := &dto.C2CMessageToCreate{
					Content: "msg",
					Media: &dto.FileInfo{
						FileInfo: resp.FileInfo,
					},
					MsgID:   data.Id,
					MsgType: 7,
					MsgReq:  1,
				}
				_, err := api.PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
				log.Info(err)
			}
			return nil
		}
		return nil
	}
	intent := websocket.RegisterHandlers(groupMessage, recall, c2cMessage, guildMsg)
	bot.NewSessionManager().Start(ws, token, &intent)
}
