package main

import (
	"context"
	"fmt"
	"math/rand"
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
)

type Push struct {
	Bot     *gonebot.Bot
	GroupId int64
}

var push = Push{}

var pushes = make(map[int64]*Push)

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
		rand.Seed(time.Now().UnixNano())
		userRole := public.IsAdmin(event.Sender.Role)
		fmt.Println(messageId, card, super, event.Sender.Nickname)
		gi, _ := bot.GetGroupInfo(groupId, true)
		gmi, _ := bot.GetGroupMemberInfo(groupId, bot.BotId, true)
		botIsAdmin := public.IsAdmin(gmi.Data.Role)
		log.Infof("[INFO] BotId(%v) GroupId(%v) UserId(%v) <- %s", botId, groupId, userId, rawMsg)

		fmt.Println("权限测试", super, botIsAdmin, userRole, gi.Data.GroupName)
		reg := regexp.MustCompile(`\[CQ:at,qq=[0-9]+\]`)
		reg1 := regexp.MustCompile(`\[CQ:reply,id=[0-9]+\]`)
		
		ss := reg.FindAllString(rawMsg, -1)
		s1 := reg1.FindAllString(rawMsg, -1)
		ns := ""
		if len(ss) == 0 {
			ns = rawMsg
		} else {
			if len(s1) > 0{
				ns = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(rawMsg, ss[0], "."), " ", ""), s1[0], ""), "..", ".")
			}else{
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
				md := "# 标题 \\n## 二级标题"
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
			retStuct := utils.PluginSet[i].Do(&ctx, botId, groupId, userId, gi.Data.GroupName, messageId, ns, card, botIsAdmin, userRole, super)
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
