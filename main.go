package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	database "github.com/2mf8/GoTBot/data"
	_ "github.com/2mf8/GoTBot/plugins"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
	gonebot "github.com/2mf8/GoneBot"
	"github.com/2mf8/GoneBot/keyboard"
	"github.com/2mf8/GoneBot/onebot"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

func main() {
	InitLog()
	tomlData := `
	Plugins = ["守卫","开关","复读","服务号","WCA","回复","频道管理","赛季","查价","打乱","学习"]   # 插件管理
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
		fmt.Printf("\n[连接] 新机器人已连接：%d\n", bot.BotId)
		fmt.Println("[已连接] 所有机器人列表：")
		for botId, _ := range gonebot.Bots {
			fmt.Println("[已连接]", botId)
		}
	}
	/*gonebot.HandleLifeTime = func(bot *gonebot.Bot, event *onebot.LifeTime) {
		fmt.Println("生命周期", event.SelfId, event.PostType, event.Time, event.MetaEventType, event.SubType)
	}
	gonebot.HandleHeartBeat = func(bot *gonebot.Bot, event *onebot.BotHeartBeat) {
		fmt.Println("心跳", event.SelfId, event.PostType, event.Time, event.MetaEventType, event.Status.Online, event.Status.AppEnabled, event.Status.AppGood, event.Status.AppInitialized, event.Status.Good)
	}*/
	gonebot.HandleGroupMessage = func(bot *gonebot.Bot, event *onebot.GroupMsgEvent) {
		groupId := event.GroupId
		rawMsg := event.RawMessage
		messageId := event.MessageId
		botId := bot.BotId
		userId := event.UserId
		card := event.Sender.Nickname
		gid := fmt.Sprintf("%v", groupId)
		uid := fmt.Sprintf("%v", userId)
		super := public.IsBotAdmin(uid, allconfig.Admins)
		rand.New(rand.NewSource(time.Now().UnixNano()))
		userRole := public.IsAdmin(event.Sender.Role)
		fmt.Println(messageId, card, super, event.Sender.Nickname)
		gi, _ := bot.GetGroupInfo(groupId, true)
		gmi, _ := bot.GetGroupMemberInfo(groupId, bot.BotId, true)
		botIsAdmin := public.IsAdmin(gmi.Data.Role)
		log.Infof("[INFO] BotId(%v) GroupId(%v) UserId(%v) <- %s", botId, groupId, userId, rawMsg)

		fmt.Println("权限测试", super, botIsAdmin, userRole, gi.Data.GroupName)
		regStr := fmt.Sprintf(`\[CQ:at,qq=%v\]`, bot.BotId)
		reg := regexp.MustCompile(regStr)
		reg1 := regexp.MustCompile(`\[CQ:reply,id=[0-9]+\]`)

		ss := reg.FindAllString(rawMsg, -1)
		s1 := reg1.FindAllString(rawMsg, -1)
		ns := ""
		if len(ss) == 0 {
			ns = rawMsg
		} else {
			if len(s1) > 0 {
				ns = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(rawMsg, ss[0], "."), " ", ""), s1[0], ""), "..", ".")
			} else {
				ns = strings.ReplaceAll(strings.ReplaceAll(rawMsg, ss[0], "."), " ", "")
			}
		}
		fmt.Println(ns)

		if ns == "mk" && super {
			kc := []*keyboard.Row{
				{
					Buttons: []*keyboard.Button{
						{
							ID: "1",
							RenderData: &keyboard.RenderData{
								Label:        "3",
								VisitedLabel: "3",
								Style:        0,
							},
							Action: &keyboard.Action{
								Type: 2,
								Permission: &keyboard.Permission{
									Type: keyboard.PermissionTypAll,
								},
								Data:                 "3",
								Reply:                true,
								Enter:                true,
								AtBotShowChannelList: true,
							},
						},
						{
							ID: "2",
							RenderData: &keyboard.RenderData{
								Label:        "4",
								VisitedLabel: "4",
								Style:        0,
							},
							Action: &keyboard.Action{
								Type: 2,
								Permission: &keyboard.Permission{
									Type: keyboard.PermissionTypAll,
								},
								Data:                 "4",
								Reply:                true,
								Enter:                true,
								AtBotShowChannelList: true,
							},
						},
						{
							ID: "5",
							RenderData: &keyboard.RenderData{
								Label:        "5",
								VisitedLabel: "5",
								Style:        0,
							},
							Action: &keyboard.Action{
								Type: 2,
								Permission: &keyboard.Permission{
									Type: keyboard.PermissionTypAll,
								},
								Data:                 "5",
								Reply:                true,
								Enter:                true,
								AtBotShowChannelList: true,
							},
						},
						{
							ID: "6",
							RenderData: &keyboard.RenderData{
								Label:        "6",
								VisitedLabel: "6",
								Style:        0,
							},
							Action: &keyboard.Action{
								Type: 2,
								Permission: &keyboard.Permission{
									Type: keyboard.PermissionTypAll,
								},
								Data:                 "6",
								Reply:                true,
								Enter:                true,
								AtBotShowChannelList: true,
							},
						},
						{
							ID: "7",
							RenderData: &keyboard.RenderData{
								Label:        "7",
								VisitedLabel: "7",
								Style:        0,
							},
							Action: &keyboard.Action{
								Type: 2,
								Permission: &keyboard.Permission{
									Type: keyboard.PermissionTypAll,
								},
								Data:                 "7",
								Reply:                true,
								Enter:                true,
								AtBotShowChannelList: true,
							},
						},
					},
				},
				{
					Buttons: []*keyboard.Button{
						{
							ID: "3",
							RenderData: &keyboard.RenderData{
								Label:        "赛季信息",
								VisitedLabel: "赛季信息",
								Style:        0,
							},
							Action: &keyboard.Action{
								Type: 2,
								Permission: &keyboard.Permission{
									Type: keyboard.PermissionTypAll,
								},
								Data:                 "赛季信息",
								Reply:                true,
								Enter:                true,
								AtBotShowChannelList: true,
							},
						},
						{
							ID: "4",
							RenderData: &keyboard.RenderData{
								Label:        "爱魔方吧",
								VisitedLabel: "孙一仝",
								Style:        1,
							},
							Action: &keyboard.Action{
								Type: 0,
								Permission: &keyboard.Permission{
									Type: keyboard.PermissionTypAll,
								},
								Data:                 "https://2mf8.cn/",
								Reply:                true,
								Enter:                true,
								AtBotShowChannelList: true,
							},
						},
					},
				},
			}
			if err == nil {
				md := fmt.Sprintf("# 标题 \\n## 二级标题\\n[直发指令](mqqapi://aio/inlinecmd?command=%s&reply=false&enter=true)\\n[手动指令](mqqapi://aio/inlinecmd?command=%s&reply=false&enter=false)\\n[🔗爱魔方吧](https://2mf8.cn/)", url.PathEscape("直发指令"), url.PathEscape("手动指令"))
				resp, err := bot.SendForwardMsg(gmi.Data.Nickname, md, kc)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(resp.Data, resp.Echo)
				lm := gonebot.NewMsg().LongMsg(resp.Data)
				rsp, err := bot.SendGroupMessage(groupId, lm, false)
				fmt.Println(rsp.Echo, err)
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
					log.Println(retStuct.ReplyMsg.Text)
					if retStuct.ReplyMsg != nil {
						newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
						if retStuct.ReplyMsg.Image != "" {
							newMsg = newMsg.Image(retStuct.ReplyMsg.Image)
						}
						resp, err := bot.SendGroupMessage(groupId, newMsg, false)
						fmt.Println(resp, resp.Data.MessageId, err)
					}
					break
				}
				if retStuct.ReqType == utils.GroupBan {
					if retStuct.BanId == 0 {
						if retStuct.ReplyMsg != nil {
							newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
							bot.SendGroupMessage(groupId, newMsg, false)
							break
						}
						break
					} else {
						bot.SetGroupBan(groupId, retStuct.BanId, retStuct.Duration)
						if retStuct.ReplyMsg != nil {
							newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
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
					if retStuct.BanId == 0 {
						break
					}
					bot.SetGroupKick(groupId, retStuct.BanId, retStuct.RejectAddAgain)
					if retStuct.ReplyMsg != nil {
						newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMessage(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.GroupLeave {
					bot.SetGroupLeave(groupId, false)
					if retStuct.ReplyMsg != nil {
						newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMessage(groupId, newMsg, false)
					}
					break
				}
				if retStuct.ReqType == utils.DeleteMsg {
					bot.DeleteMsg(retStuct.MsgId)
					if retStuct.ReplyMsg != nil {
						newMsg := gonebot.NewMsg().Text(retStuct.ReplyMsg.Text)
						bot.SendGroupMessage(groupId, newMsg, false)
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
			fmt.Println("[失败] 创建机器人失败")
		}
	})

	if err := router.Run(":8082"); err != nil {
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
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Warn("登录失败，请检查 appid 和 AccessToken 是否正确。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
		log.Printf("%+v, err:%v", ws, err)
	}
	var groupMessage event.GroupAtMessageEventHandler = func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		groupId := data.GroupId
		userId := data.Author.UserId
		content := strings.TrimSpace(data.Content)
		msgId := data.MsgId
		reg4 := regexp.MustCompile("/")
		content = strings.TrimSpace(reg4.ReplaceAllString(content, ""))
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)
		log.Printf("[INFO] GroupId(%v) UserId(%v) -> %v", groupId, userId, content)
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
		for v, i := range database.AllConfig.Plugins {
			fmt.Println(v,i)
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
					log.Println(retStuct.ReplyMsg.Text)
					if retStuct.ReplyMsg != nil {
						fmt.Println(retStuct.ReplyMsg.Text)
						if retStuct.ReplyMsg.Image != "" {
							resp, _ := api.PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Image, SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: retStuct.ReplyMsg.Text, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
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
								Content: "\n" + retStuct.ReplyMsg.Text, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
								MsgID:   data.MsgId,
								MsgType: 0,
							}
							api.PostGroupMessage(ctx, groupId, newMsg)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							resp, _ := api.PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: retStuct.ReplyMsg.Text, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
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
									Content: retStuct.ReplyMsg.Text, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
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
									Content: retStuct.ReplyMsg.Text, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
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
	intent := websocket.RegisterHandlers(groupMessage)
	botgo.NewSessionManager().Start(ws, token, &intent)
}
