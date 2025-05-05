package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	bot1 "github.com/2mf8/Better-Bot-Go"
	bytesimage "github.com/2mf8/Better-Bot-Go/bytes_image"
	"github.com/2mf8/Better-Bot-Go/dto"
	"github.com/2mf8/Better-Bot-Go/openapi"
	v1 "github.com/2mf8/Better-Bot-Go/openapi/v1"
	"github.com/2mf8/Better-Bot-Go/token"
	"github.com/2mf8/Better-Bot-Go/webhook"
	database "github.com/2mf8/GoTBot/data"
	"github.com/2mf8/GoTBot/plugins"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
	log "github.com/sirupsen/logrus"
)

type contextKey string

const (
	GloKey contextKey = "2mf8" // 可以使用字符串或其他值标识不同的 key
)

func main() {
	webhook.InitLog()
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

	ctx := context.WithValue(context.Background(), GloKey, "cn2mf8")

	as := webhook.ReadSetting()
	for _, v := range as.Apps {
		atr := v1.GetAccessToken(fmt.Sprintf("%v", v.AppId), v.AppSecret)
		iat, err := strconv.Atoi(atr.ExpiresIn)
		if err == nil && atr.AccessToken != "" {
			aei := time.Now().Unix() + int64(iat)
			token := token.BotToken(v.AppId, atr.AccessToken, string(token.TypeQQBot))
			if v.IsSandBox {
				api := bot1.NewSandboxOpenAPI(token).WithTimeout(3 * time.Second)
				go bot1.AuthAcessAdd(fmt.Sprintf("%v", v.AppId), &bot1.AccessToken{AccessToken: atr.AccessToken, ExpiresIn: aei, Api: api, AppSecret: v.AppSecret, IsSandBox: v.IsSandBox, Appid: v.AppId})
			} else {
				api := bot1.NewOpenAPI(token).WithTimeout(3 * time.Second)
				go bot1.AuthAcessAdd(fmt.Sprintf("%v", v.AppId), &bot1.AccessToken{AccessToken: atr.AccessToken, ExpiresIn: aei, Api: api, AppSecret: v.AppSecret, IsSandBox: v.IsSandBox, Appid: v.AppId})
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	b, _ := json.Marshal(as)
	fmt.Println("配置", string(b))
	webhook.GroupAtMessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		groupId := data.GroupId
		userId := data.Author.UserId
		content := strings.TrimSpace(data.Content)
		msgId := data.MsgId
		content, _ = public.Prefix(content, "/")
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)
		content = fmt.Sprintf(".%s", content)
		sg, _ := database.SGBGIACI(groupId, groupId)
		if content == ".GetID" {
			newMsg := &dto.GroupMessageToCreate{
				Content: userId,
				MsgID:   data.MsgId,
				MsgType: 0,
			}
			bot1.SendApi(bot.XBotAppid[0]).PostGroupMessage(ctx, groupId, newMsg)
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
			mi, err := bot1.SendApi(bot.XBotAppid[0]).PostGroupMessage(ctx, data.GroupId, &dto.C2CMessageToCreate{
				Content: "测试撤回",
				MsgType: dto.C2CMsgTypeText,
				MsgID:   data.MsgId,
			})
			if err == nil {
				fmt.Println(mi.Id, mi.Timestamp)
				go func() {
					time.Sleep(time.Second * 10)
					bot1.SendApi(bot.XBotAppid[0]).DelGroupBotMessage(ctx, data.GroupId, mi.Id, openapi.RetractMessageOptionHidetip)
				}()
			} else {
				fmt.Println(err)
			}
		}

		for _, i := range database.AllConfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := plugins.PluginMap[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &msgIdType, content, "", true, false, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						msg := fmt.Sprintf("\n%s", strings.TrimSpace(retStuct.ReplyMsg.Text))
						if retStuct.ReplyMsg.Image != "" {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Image)
							if err == nil {
								resp, _ := bot1.SendApi(bot.XBotAppid[0]).PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
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
									bot1.SendApi(bot.XBotAppid[0]).PostGroupMessage(ctx, groupId, newMsg)
								}
							}
						} else {
							newMsg := &dto.GroupMessageToCreate{
								Content: msg, //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
								MsgID:   data.MsgId,
								MsgType: 0,
							}
							bot1.SendApi(bot.XBotAppid[0]).PostGroupMessage(ctx, groupId, newMsg)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								resp, _ := bot1.SendApi(bot.XBotAppid[0]).PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
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
									bot1.SendApi(bot.XBotAppid[0]).PostGroupMessage(ctx, groupId, newMsg)
								}
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								resp, _ := bot1.SendApi(bot.XBotAppid[0]).PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
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
									bot1.SendApi(bot.XBotAppid[0]).PostGroupMessage(ctx, groupId, newMsg)
								}
							}
							s1, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[2])
							if err == nil {
								resp1, _ := bot1.SendApi(bot.XBotAppid[0]).PostGroupRichMediaMessage(ctx, groupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, FileData: s1, SrvSendMsg: false})
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
									bot1.SendApi(bot.XBotAppid[0]).PostGroupMessage(ctx, groupId, newMsg)
								}
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
		content = fmt.Sprintf(".%s", content)
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
			mi, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CMessage(ctx, data.Author.UserOpenId, &dto.C2CMessageToCreate{
				Content: "测试撤回",
				MsgType: dto.C2CMsgTypeText,
				MsgID:   data.Id,
			})
			if err == nil {
				fmt.Println(mi.Id, mi.Timestamp)
				go func() {
					time.Sleep(time.Second * 10)
					bot1.SendApi(bot.XBotAppid[0]).DelC2CMessage(ctx, data.Author.UserOpenId, mi.Id, openapi.RetractMessageOptionHidetip)
				}()
			} else {
				fmt.Println(err)
			}
		}

		for _, i := range database.AllConfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := plugins.PluginMap[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &msgIdType, content, "", true, false, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						if retStuct.ReplyMsg.Image != "" {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Image)
							if err == nil {
								resp, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
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
									_, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
									log.Info(err)
								}
							}
						} else {
							newMsg := &dto.C2CMessageToCreate{
								Content: strings.TrimSpace(retStuct.ReplyMsg.Text),
								MsgID:   data.Id,
								MsgType: 0,
								MsgReq:  1,
							}
							_, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
							log.Info(err)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								resp, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
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
									_, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
									log.Info(err)
								}
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								resp, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
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
									_, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
									log.Info(err)
								}
							}
							s1, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[2])
							if err == nil {
								resp1, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.C2CRichMediaMessageToCreate{FileType: 1, FileData: s1, SrvSendMsg: false})
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
									_, err := bot1.SendApi(bot.XBotAppid[0]).PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
									log.Info(err)
								}
							}
						}
					}
					break
				}
			}
		}
		return nil
	}
	webhook.ATMessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSATMessageData) error {
		guildId := data.GuildID
		channelId := data.ChannelID
		userId := data.Author.ID
		content := strings.TrimSpace(data.Content)
		msgId := data.ID
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)
		me, _ := bot1.SendApi(bot.XBotAppid[0]).Me(ctx)
		reg7 := regexp.MustCompile(fmt.Sprintf("<@!%s> ", me.ID))
		reg4 := regexp.MustCompile(fmt.Sprintf("<@!%s> /", me.ID))
		content = strings.TrimSpace(reg4.ReplaceAllString(content, "."))
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
			retStuct := plugins.PluginMap[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &msgIdType, content, "", true, false, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						msg := strings.TrimSpace(retStuct.ReplyMsg.Text)
						if retStuct.ReplyMsg.Image != "" {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Image)
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
									"content": msg,
								}, "333.png", bytes.NewReader(s))
							}
						} else {
							bot1.SendApi(bot.XBotAppid[0]).PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								MsgID:   data.ID,
							})
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
								}, "333.png", bytes.NewReader(s))
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
								}, "333.png", bytes.NewReader(s))
							}
							s1, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[2])
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
								}, "333.png", bytes.NewReader(s1))
							}
						}
					}
					break
				}
			}
		}
		return nil
	}
	webhook.MessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageData) error {
		guildId := data.GuildID
		channelId := data.ChannelID
		userId := data.Author.ID
		content := strings.TrimSpace(data.Content)
		msgId := data.ID
		super := public.IsBotAdmin(userId, database.AllConfig.Admins)

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
			retStuct := plugins.PluginMap[i].Do(&ctx, &botType, &groupIdType, &userIdType, "", &msgIdType, content, "", true, false, super)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GroupMsg {
					if retStuct.ReplyMsg != nil {
						msg := strings.TrimSpace(retStuct.ReplyMsg.Text)
						if retStuct.ReplyMsg.Image != "" {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Image)
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
									"content": msg,
								}, "333.png", bytes.NewReader(s))
							}
						} else {
							bot1.SendApi(bot.XBotAppid[0]).PostMessage(ctx, channelId, &dto.MessageToCreate{
								Content: msg,
								MsgID:   data.ID,
							})
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
								}, "333.png", bytes.NewReader(s))
							}
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							s, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[1])
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
								}, "333.png", bytes.NewReader(s))
							}
							s1, err := bytesimage.GetImageBytes(retStuct.ReplyMsg.Images[2])
							if err == nil {
								bot1.SendApi(bot.XBotAppid[0]).PostFormFileReaderImage(ctx, channelId, map[string]string{
									"msg_id":  data.ID,
								}, "333.png", bytes.NewReader(s1))
							}
						}
					}
					break
				}
			}
		}
		return nil
	}
	webhook.InitGin(as.IsOpen)
	select {}
}
