package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/2mf8/GoTBot/config"
	"github.com/2mf8/GoTBot/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
)

type Switch struct {
	Id             int       `json:"id"`
	GuildId        string    `json:"guild_id"`
	ChannelId      string    `json:"channel_id"`
	IsCloseOrGuard int64     `json:"is_close_or_guard"`
	AdminId        string    `json:"admin_id"`
	GmtModified    time.Time `json:"gmt_modified"`
}

type SwitchSync struct {
	IsTrue       bool `json:"synchronization"`
	PluginSwitch *Switch
}

type intent int

const (
	PluginGuard intent = 1 << iota // 守卫
	PluginLog
	PluginBlock    // 个人屏蔽
	PluginSwitch   // 开关
	PluginRepeat   // 复读
	PluginWCA      // WCA
	PluginReply    // 回复
	PluginAdmin    // 群管
	PluginPrice    // 查价
	PluginScramble // 打乱
	PluginRank     // Rank
	PluginLearn    // 频道学习
)

var IntentMap = map[intent]string{
	PluginLog:      "Log",
	PluginGuard:    "守卫",
	PluginBlock:    "屏蔽",
	PluginSwitch:   "开关",
	PluginRepeat:   "复读",
	PluginWCA:      "WCA",
	PluginReply:    "回复",
	PluginAdmin:    "群管",
	PluginPrice:    "查价",
	PluginScramble: "打乱",
	PluginRank:     "Rank",
	PluginLearn:    "学习",
}

var SwitchMap = map[string]intent{
	"Log":  PluginLog,
	"守卫":   PluginGuard,
	"屏蔽":   PluginBlock,
	"开关":   PluginSwitch,
	"复读":   PluginRepeat,
	"WCA":  PluginWCA,
	"回复":   PluginReply,
	"群管":   PluginAdmin,
	"查价":   PluginPrice,
	"打乱":   PluginScramble,
	"Rank": PluginRank,
	"学习":   PluginLearn,
}

func (bot_switch *Switch) SwitchCreate() (err error) {
	statement := fmt.Sprintf("insert into [%s].[dbo].[guild_switch] values ($1, $2, $3, $4, $5) select @@identity", config.Conf.DatabaseName)
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(bot_switch.GuildId, bot_switch.ChannelId, bot_switch.IsCloseOrGuard, bot_switch.AdminId, bot_switch.GmtModified).Scan(&bot_switch.Id)

	bot_switch_sync := SwitchSync{
		IsTrue: true,
		PluginSwitch: &Switch{
			Id:             bot_switch.Id,
			GuildId:        bot_switch.GuildId,
			ChannelId:      bot_switch.ChannelId,
			IsCloseOrGuard: bot_switch.IsCloseOrGuard,
			AdminId:        bot_switch.AdminId,
			GmtModified:    bot_switch.GmtModified,
		},
	}

	// byte write
	bw := bot_switch.GuildId + "_" + bot_switch.ChannelId + "_switchorguard"
	var bot_switch_redis []byte
	bot_switch_redis, err = json.Marshal(&bot_switch_sync)
	if err != nil {
		fmt.Println("[错误] Marshal序列化出错")
	}
	c := Pool.Get()
	defer c.Close()
	c.Send("SET", bw, bot_switch_redis)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Printf("[收到] %#v\n", v)

	return
}

func (bot_switch *Switch) SwitchUpdate() (err error) {
	statment := fmt.Sprintf("update [%s].[dbo].[guild_switch] set guild_id = $2, channel_id = $3, is_close_or_guard = $4, admin_id = $5, gmt_modified = $6 where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, bot_switch.Id, bot_switch.GuildId, bot_switch.ChannelId, bot_switch.IsCloseOrGuard, bot_switch.AdminId, bot_switch.GmtModified)
	if err != nil {
		return err
	}

	ubot_switch := Switch{
		Id:             bot_switch.Id,
		GuildId:        bot_switch.GuildId,
		ChannelId:      bot_switch.ChannelId,
		IsCloseOrGuard: bot_switch.IsCloseOrGuard,
		AdminId:        bot_switch.AdminId,
		GmtModified:    bot_switch.GmtModified,
	}

	bot_switch_sync := SwitchSync{
		IsTrue:       true,
		PluginSwitch: &ubot_switch,
	}

	bw := bot_switch.GuildId + "_" + bot_switch.ChannelId + "_switchorguard"
	var bot_switch_redis []byte
	bot_switch_redis, err = json.Marshal(&bot_switch_sync)
	if err != nil {
		fmt.Println("[错误] Marshal序列化出错")
	}
	c := Pool.Get()
	defer c.Close()
	c.Send("SET", bw, bot_switch_redis)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Printf("[收到] %#v\n", v)

	return
}

func SwitchSave(guildId, channelId, adminId string, isCloseOrGuard int64, gmtModified time.Time, isClose bool) (err error) {
	var icog int64 = 0
	bot_switch := Switch{
		GuildId:        guildId,
		ChannelId:      channelId,
		IsCloseOrGuard: isCloseOrGuard,
		AdminId:        adminId,
		GmtModified:    gmtModified,
	}
	bot_switch_sync := SwitchSync{
		IsTrue:       true,
		PluginSwitch: &bot_switch,
	}
	switch_get, err := SGBGIACI(guildId, channelId)
	fmt.Println(switch_get.IsTrue)
	if err != nil || !switch_get.IsTrue {
		err = bot_switch_sync.PluginSwitch.SwitchCreate()
		return err
	}
	if isClose {
		icog = switch_get.PluginSwitch.IsCloseOrGuard | isCloseOrGuard
	} else {
		icog = ^isCloseOrGuard & switch_get.PluginSwitch.IsCloseOrGuard
	}
	switch_get.PluginSwitch.IsCloseOrGuard = icog
	switch_get.PluginSwitch.AdminId = adminId
	switch_get.PluginSwitch.GmtModified = gmtModified
	err = switch_get.PluginSwitch.SwitchUpdate()
	return err
}

// SDBGI SwitchDeleteByGuildIdAndChannelId
func SDBGIACI(guildId, channelId string) (err error) {
	statment := fmt.Sprintf("delete [%s].[dbo].[guild_switch] where guild_id = $1 and channel_id = $2", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, guildId, channelId)
	if err != nil {
		return err
	}
	return
}

// SGBGI SwitchGetByGuildIdAndChannelId
func SGBGIACI(guildId, channelId string) (bot_switch_sync SwitchSync, err error) {
	bot_switch := Switch{}
	bot_switch_sync = SwitchSync{
		IsTrue:       true,
		PluginSwitch: &bot_switch,
	}
	bw := guildId + "_" + channelId + "_switchorguard"
	c := Pool.Get()
	defer c.Close()
	/*exists, err := redis.Bool(c.Do("exists", bw))
	if err != nil {
		fmt.Println("不存在")
	}
	fmt.Println(exists)*/
	c.Send("Get", bw)
	c.Flush()
	// value byte
	var vb []byte
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		log.Println(err)
		fmt.Println("[查询] 首次查询-开关", bw)
		statment := fmt.Sprintf("select ID, guild_id, channel_id, is_close_or_guard, admin_id, gmt_modified from [%s].[dbo].[guild_switch] where guild_id = $1 and channel_id = $2", config.Conf.DatabaseName)
		err = Db.QueryRow(statment, guildId, channelId).Scan(&bot_switch_sync.PluginSwitch.Id, &bot_switch_sync.PluginSwitch.GuildId, &bot_switch_sync.PluginSwitch.ChannelId, &bot_switch_sync.PluginSwitch.IsCloseOrGuard, &bot_switch_sync.PluginSwitch.AdminId, &bot_switch_sync.PluginSwitch.GmtModified)
		info := fmt.Sprintf("%s", err)
		if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
			if public.StartsWith(info, "unable") {
				fmt.Println(info)
			}
			bot_switch_sync = SwitchSync{
				IsTrue:       false,
				PluginSwitch: &bot_switch,
			}
		}
		var bw_set []byte
		bw_set, _ = json.Marshal(&bot_switch_sync)
		c.Send("SET", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return
	}
	err = json.Unmarshal(vb, &bot_switch_sync)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	fmt.Println("[Redis] Key(", bw, ") Value(", bot_switch_sync.IsTrue, *bot_switch_sync.PluginSwitch, ")") //测试用
	return
}

func IntentMean(intent intent) string {
	mean, ok := IntentMap[intent]
	if !ok {
		mean = "unknown"
	}
	return mean
}

func PluginNameToIntent(s string) intent {
	intent, ok := SwitchMap[s]
	if !ok {
		intent = 0
	}
	return intent
}
