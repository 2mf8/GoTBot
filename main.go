package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	bot "github.com/2mf8/Better-Bot-Go"
	"github.com/2mf8/Better-Bot-Go/dto"
	"github.com/2mf8/Better-Bot-Go/dto/keyboard"
	"github.com/2mf8/Better-Bot-Go/openapi"
	"github.com/2mf8/Better-Bot-Go/token"
	"github.com/2mf8/Better-Bot-Go/webhook"
	database "github.com/2mf8/GoTBot/data"
	_ "github.com/2mf8/GoTBot/plugins"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
	gonebot "github.com/2mf8/GoneBot"
	gkb "github.com/2mf8/GoneBot/keyboard"
	"github.com/2mf8/GoneBot/onebot"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/* type SetBotCallBack struct {
	BotId      uint   `json:"bot_id,omitempty"`
	GroupId    int64  `json:"group_id,omitempty"`
	ButtonId   string `json:"button_id,omitempty"`   // button_id
	ButtonData string `json:"button_data,omitempty"` // button_data
} */

var Apis = make(map[string]openapi.OpenAPI, 0)

type Resolved struct {
	ButtonId   string `json:"button_id,omitempty"`   // button_id
	ButtonData string `json:"button_data,omitempty"` // button_data
}

/*
	 {
		"bot_id":101981675,
		"group_id":121196301,
		"button_id":"id",
		"ButtonData":"data",
	}
*/
func main() {
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
	gonebot.HandleGroupMessage = func(bot *gonebot.Bot, ievent *onebot.GroupMsgEvent) {
		groupId := ievent.GroupId
		rawMsg := ievent.RawMessage
		rand.New(rand.NewSource(time.Now().UnixNano()))

		if strings.HasPrefix(rawMsg, ".") || strings.HasPrefix(rawMsg, "%") {
			bot.SendGroupBotCallback(102070767, groupId, "1", rawMsg)
		}
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/onebot/v11/ws", func(c *gin.Context) {
		if err := gonebot.UpgradeWebsocket(c.Writer, c.Request); err != nil {
			log.Info("[失败] 创建机器人失败")
		}
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
	select {}
}

func StartOffical() {
	webhook.InitLog()
	as := webhook.ReadSetting()
	var ctx context.Context
	for i, v := range as.Apps {
		token := token.BotToken(v.AppId, v.Token, string(token.TypeBot))
		api := bot.NewOpenAPI(token).WithTimeout(3 * time.Second)
		Apis[i] = api
	}
	b, _ := json.Marshal(as)
	fmt.Println("配置", string(b))
	webhook.GroupAtMessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
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
		if content == ".GetID" {
			newMsg := &dto.GroupMessageToCreate{
				Content: userId,
				MsgID:   data.MsgId,
				MsgType: 0,
			}
			Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, groupId, newMsg)
		}
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
		if content == ".del" {
			mi, err := Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupId, &dto.C2CMessageToCreate{
				Content: "测试撤回",
				MsgType: dto.C2CMsgTypeText,
				MsgID: data.MsgId,
			})
			if err == nil {
				fmt.Println(mi.Id, mi.Timestamp)
				go func ()  {
					time.Sleep(time.Second * 10)
					Apis[bot.XBotAppid[0]].DelGroupBotMessage(ctx, data.GroupId, mi.Id, openapi.RetractMessageOptionHidetip)
				}()
			}else{
				fmt.Println(err)
			}
		}
		if content == ".get" {
			gm, err := Apis[bot.XBotAppid[0]].GetGroupMembers(ctx, data.GroupId, 0, 0)
			b,_:=json.Marshal(gm)
			fmt.Println(string(b), err)
		}
		if content == ".at" {
			Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupId, &dto.C2CMessageToCreate{
				Content: "测试<qqbot-at-user id="+data.Author.UserId+" />",
				MsgType: dto.C2CMsgTypeText,
			})
		}
		if content == ".k" {
			/* rows := keyboard.CustomKeyboard{} */
			/* kb := gkb.Builder().
			TextButton("测试", "已测试", "成功", false, true).
			UrlButton("爱魔方吧", "一仝", "https://2mf8.cn", false, true).
			SetRow().
			TextButton("测试", "已测试", "成功", false, true).
			SetRow()
			b, _:= json.Marshal(kb)
			json.Unmarshal(b, &rows) */
			fmt.Println("测试")
			Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupId, &dto.C2CMessageToCreate{
				Keyboard: &keyboard.MessageKeyboard{
					ID: "101981675_1734764173",
				},
				MsgType: dto.C2CMsgTypeMarkdown,
				MsgID:   data.MsgId,
			})
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
							resp, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Image, SrvSendMsg: false})
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
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, groupId, newMsg)
							}
						} else {
							newMsg := &dto.GroupMessageToCreate{
								Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
								MsgID:   data.MsgId,
								MsgType: 0,
							}
							Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, groupId, newMsg)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							resp, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
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
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, groupId, newMsg)
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							resp, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
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
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, groupId, newMsg)
							}
							resp1, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[2], SrvSendMsg: false})
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
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, groupId, newMsg)
							}
						}
					}
					break
				}
			}
		}
		return nil
	}
	webhook.C2CMessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSC2CMessageData) error {
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
		if content == ".del" {
			mi, err := Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, &dto.C2CMessageToCreate{
				Content: "测试撤回",
				MsgType: dto.C2CMsgTypeText,
				MsgID: data.Id,
			})
			if err == nil {
				fmt.Println(mi.Id, mi.Timestamp)
				go func ()  {
					time.Sleep(time.Second * 10)
					Apis[bot.XBotAppid[0]].DelC2CMessage(ctx, data.Author.UserOpenId, mi.Id, openapi.RetractMessageOptionHidetip)
				}()
			}else{
				fmt.Println(err)
			}
		}
		if content == "k" {
			kb := gkb.Builder().
				TextButton("测试", "已测试", "成功", false, true).
				UrlButton("爱魔方吧", "一仝", "https://2mf8.cn", false, true).
				SetRow().
				TextButton("测试", "已测试", "成功", false, true).
				SetRow()
			b, _ := json.Marshal(kb)
			fmt.Println(string(b))
			Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, &dto.C2CMessageToCreate{
				Keyboard: &keyboard.MessageKeyboard{
					ID: "101981675_1734764173",
				},
				MsgType: dto.C2CMsgTypeMarkdown,
				MsgID:   data.Id,
			})
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
							resp, err := Apis[bot.XBotAppid[0]].PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Image, SrvSendMsg: false})
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
								_, err := Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
						} else {
							newMsg := &dto.C2CMessageToCreate{
								Content: strings.TrimSpace(retStuct.ReplyMsg.Text),
								MsgID:   data.Id,
								MsgType: 0,
								MsgReq:  1,
							}
							_, err := Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
							log.Info(err)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							resp, err := Apis[bot.XBotAppid[0]].PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
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
								_, err := Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							resp, err := Apis[bot.XBotAppid[0]].PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
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
								_, err := Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
							resp1, err := Apis[bot.XBotAppid[0]].PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[2], SrvSendMsg: false})
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
								_, err := Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
								log.Info(err)
							}
						}
					}
					break
				}
			}
		}

		if data.Content == "测试" {
			resp, err := Apis[bot.XBotAppid[0]].PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, Url: "https://www.2mf8.cn/static/image/cube3/b1.png", SrvSendMsg: false})
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
				_, err := Apis[bot.XBotAppid[0]].PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
				log.Info(err)
			}
			return nil
		}
		return nil
	}
	webhook.MessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageData) error {
		guildId := data.GuildID
		channelId := data.ChannelID
		userId := data.Author.ID
		content := strings.TrimSpace(data.Content)
		msgId := data.ID
		reg4 := regexp.MustCompile("/")
		content = strings.TrimSpace(reg4.ReplaceAllString(content, ""))
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)
		ctx := context.WithValue(context.Background(), "key", "value")
		me, _ := Apis[bot.XBotAppid[0]].Me(ctx)
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
		if content == ".k" {
			/* rows := keyboard.CustomKeyboard{} */
			/* kb := gkb.Builder().
			TextButton("测试", "已测试", "成功", false, true).
			UrlButton("爱魔方吧", "一仝", "https://2mf8.cn", false, true).
			SetRow().
			TextButton("测试", "已测试", "成功", false, true).
			SetRow()
			b, _:= json.Marshal(kb)
			json.Unmarshal(b, &rows) */
			fmt.Println("测试")
			Apis[bot.XBotAppid[0]].PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{
				Keyboard: &keyboard.MessageKeyboard{
					ID: "101981675_1734764173",
				},
				MsgID:   data.ID,
			})
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
							Apis[bot.XBotAppid[0]].PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								Image:   retStuct.ReplyMsg.Image,
								MsgID:   data.ID,
							})
						} else {
							Apis[bot.XBotAppid[0]].PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								MsgID:   data.ID,
							})
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							Apis[bot.XBotAppid[0]].PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								Image:   retStuct.ReplyMsg.Images[1],
								MsgID:   data.ID,
							})
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							Apis[bot.XBotAppid[0]].PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								Image:   retStuct.ReplyMsg.Images[1],
								MsgID:   data.ID,
							})
							Apis[bot.XBotAppid[0]].PostMessage(ctx, channelId, &dto.MessageToCreate{
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
	webhook.InteractionEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSInteractionData) error {
		dr := &Resolved{}
		b, err := json.Marshal(data.Data.Resolved)
		if err != nil {
			return err
		}
		json.Unmarshal(b, dr)
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
		for _, i := range database.AllConfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			fmt.Println("eventId", data.ID)
			retStuct := utils.PluginSet[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &utils.MsgIdType{}, dr.ButtonData, "", true, false, false)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						msg := strings.TrimSpace(retStuct.ReplyMsg.Text)
						if retStuct.ReplyMsg.Image != "" {
							resp, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Image, SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									EventID: dto.EventType(data.ID),
									MsgType: 7,
									MsgReq:  1,
								}
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupOpenID, newMsg)
							}
						} else {
							newMsg := &dto.GroupMessageToCreate{
								Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
								EventID: dto.EventType(data.ID),
								MsgType: 0,
							}
							Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupOpenID, newMsg)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							resp, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									EventID: dto.EventType(data.ID),
									MsgType: 7,
									MsgReq:  2,
								}
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupOpenID, newMsg)
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							resp, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[1], SrvSendMsg: false})
							if resp != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp.FileInfo,
									},
									EventID: dto.EventType(data.ID),
									MsgType: 7,
									MsgReq:  2,
								}
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupOpenID, newMsg)
							}
							resp1, _ := Apis[bot.XBotAppid[0]].PostGroupRichMediaMessage(ctx, data.GroupOpenID, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: retStuct.ReplyMsg.Images[2], SrvSendMsg: false})
							if resp1 != nil {
								newMsg := &dto.GroupMessageToCreate{
									Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
									Media: &dto.FileInfo{
										FileInfo: resp1.FileInfo,
									},
									EventID: dto.EventType(data.ID),
									MsgType: 7,
									MsgReq:  3,
								}
								Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupOpenID, newMsg)
							}
						}
					}
					break
				}
			}
		}
		return nil
	}
	webhook.InitGin()
	select {}
}
