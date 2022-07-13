package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/2mf8/go-tbot-for-rq/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
)

type Switch struct {
	Id             int   `json:"id"`
	GroupId        int64 `json:"group_id"`
	IsCloseOrGuard int64 `json:"is_close_or_guard"`
}

type SwitchSync struct {
	IsTrue       bool `json:"synchronization"`
	PluginSwitch *Switch
}

type intent int

const (
	PluginGuard intent = 1 << iota // 守卫
	PluginBlock // 个人屏蔽
	PluginSwitch // 开关
	PluginRepeat // 复读
	PluginReply // 回复
	PluginAdmin // 群管
	PluginSubscribe //查价订阅
	PluginPrice // 查价
	PluginScramble // 打乱
	PluginLearn // 群学习
)

var IntentMap = map[intent]string{
	PluginGuard:     "守卫",
	PluginBlock:     "屏蔽",
	PluginSwitch:    "开关",
	PluginRepeat:    "复读",
	PluginReply:     "回复",
	PluginAdmin:     "群管",
	PluginSubscribe: "订阅",
	PluginPrice:     "查价",
	PluginScramble:  "打乱",
	PluginLearn:     "学习",
}

var SwitchMap = map[string]intent{
	"守卫": PluginGuard,
	"屏蔽": PluginBlock,
	"开关": PluginSwitch,
	"复读": PluginRepeat,
	"回复": PluginReply,
	"群管": PluginAdmin,
	"订阅": PluginSubscribe,
	"查价": PluginPrice,
	"打乱": PluginScramble,
	"学习": PluginLearn,
}

var ctx = context.Background()

func (bot_switch *Switch) SwitchCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[zbot_plugin_switch] values ($1, $2) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(bot_switch.GroupId, bot_switch.IsCloseOrGuard).Scan(&bot_switch.Id)

	bot_switch_sync := SwitchSync{
		IsTrue: true,
		PluginSwitch: &Switch{
			Id:             bot_switch.Id,
			GroupId:        bot_switch.GroupId,
			IsCloseOrGuard: bot_switch.IsCloseOrGuard,
		},
	}

	// byte write
	bw := strconv.Itoa(int(bot_switch.GroupId)) + "_switchorguard"
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

func (bot_switch *Switch) SwitchUpdate(isCloseOrGuard int64) (err error) {

	_, err = Db.Exec("update [kequ5060].[dbo].[zbot_plugin_switch] set group_id = $2, is_close_or_guard = $3 where ID = $1", bot_switch.Id, bot_switch.GroupId, isCloseOrGuard)
	if err != nil {
		return err
	}

	ubot_switch := Switch{
		Id:             bot_switch.Id,
		GroupId:        bot_switch.GroupId,
		IsCloseOrGuard: isCloseOrGuard,
	}

	bot_switch_sync := SwitchSync{
		IsTrue:       true,
		PluginSwitch: &ubot_switch,
	}

	bw := strconv.Itoa(int(bot_switch.GroupId)) + "_switchorguard"
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

func SwitchSave(groupId int64, isCloseOrGuard int64, isClose bool) (err error) {
	var icog int64 = 0
	bot_switch := Switch{
		GroupId:        groupId,
		IsCloseOrGuard: isCloseOrGuard,
	}
	bot_switch_sync := SwitchSync{
		IsTrue:       true,
		PluginSwitch: &bot_switch,
	}
	switch_get, err := SGBGI(groupId)
	if err != nil || !switch_get.IsTrue {
		err = bot_switch_sync.PluginSwitch.SwitchCreate()
		return err
	}
	if isClose {
		icog = switch_get.PluginSwitch.IsCloseOrGuard | isCloseOrGuard
	} else {
		icog = ^isCloseOrGuard & switch_get.PluginSwitch.IsCloseOrGuard
	}
	err = switch_get.PluginSwitch.SwitchUpdate(icog)
	return err
}

// SDBGI SwitchDeleteByGroupId
func SDBGI(groupId int64) (err error) {
	_, err = Db.Exec("delete [kequ5060].[dbo].[zbot_plugin_switch] where group_id = $1", groupId)
	if err != nil {
		return err
	}
	return
}

//SGBGI SwitchGetByGroupId
func SGBGI(groupId int64) (bot_switch_sync SwitchSync, err error) {
	bot_switch := Switch{}
	bot_switch_sync = SwitchSync{
		IsTrue:       true,
		PluginSwitch: &bot_switch,
	}
	bw := strconv.Itoa(int(groupId)) + "_switchorguard"
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
		fmt.Println("[查询] 首次查询-开关", bw)
		err = Db.QueryRow("select ID, group_id, is_close_or_guard from [kequ5060].[dbo].[zbot_plugin_switch] where group_id = $1", groupId).Scan(&bot_switch_sync.PluginSwitch.Id, &bot_switch_sync.PluginSwitch.GroupId, &bot_switch_sync.PluginSwitch.IsCloseOrGuard)
		info := fmt.Sprintf("%s", err)
		if StartsWith(info, "sql") || StartsWith(info, "unable") {
			if StartsWith(info, "unable") {
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
	//fmt.Println("[Redis] Key(", bw, ") Value(", bot_switch_sync.IsTrue, *bot_switch_sync.PluginSwitch, ")")  //测试用
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
