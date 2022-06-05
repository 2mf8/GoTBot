package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	. "github.com/2mf8/go-tbot-for-rq/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
)

type Switch struct {
	Id          int       `json:"id"`
	GroupId     int64     `json:"group_id"`
	PluginName  string    `json:"plugin_name"`
	GmtModified time.Time `json:"gmt_modified"`
	Stop        bool      `json:"stop"`
}

type SwitchStop struct {
	GroupId    int64 `json:"group_id"`
	PluginName []string
}

type Plu struct {
	PluginName string `json:"plugin_name"`
}

type SwitchSync struct {
	IsTrue       bool `json:"synchronization"`
	PluginSwitch *Switch
}

var ctx = context.Background()

func (bot_switch *Switch) SwitchCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[zbot_plugin_switch] values ($1, $2, $3, $4) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(bot_switch.GroupId, bot_switch.PluginName, bot_switch.GmtModified, bot_switch.Stop).Scan(&bot_switch.Id)

	bot_switch_sync := SwitchSync{
		IsTrue: true,
		PluginSwitch: &Switch{
			Id:          bot_switch.Id,
			GroupId:     bot_switch.GroupId,
			PluginName:  bot_switch.PluginName,
			GmtModified: bot_switch.GmtModified,
			Stop:        bot_switch.Stop,
		},
	}

	// byte write
	bw := strconv.Itoa(int(bot_switch.GroupId)) + "_" + bot_switch.PluginName
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

func (bot_switch *Switch) SwitchUpdate(stop bool) (err error) {

	_, err = Db.Exec("update [kequ5060].[dbo].[zbot_plugin_switch] set group_id = $2, plugin_name = $3, gmt_modified = $4, stop = $5 where ID = $1", bot_switch.Id, bot_switch.GroupId, bot_switch.PluginName, bot_switch.GmtModified, stop)
	if err != nil {
		return err
	}

	ubot_switch := Switch{
		Id:          bot_switch.Id,
		GroupId:     bot_switch.GroupId,
		PluginName:  bot_switch.PluginName,
		GmtModified: bot_switch.GmtModified,
		Stop:        bool(stop),
	}

	bot_switch_sync := SwitchSync{
		IsTrue:       true,
		PluginSwitch: &ubot_switch,
	}

	bw := strconv.Itoa(int(bot_switch.GroupId)) + "_" + bot_switch.PluginName
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

func SwitchSave(groupId int64, pluginName string, gmtModified time.Time, stop bool) (err error) {
	bot_switch := Switch{
		GroupId:     groupId,
		PluginName:  pluginName,
		GmtModified: time.Now(),
		Stop:        bool(stop),
	}
	bot_switch_sync := SwitchSync{
		IsTrue:       true,
		PluginSwitch: &bot_switch,
	}
	switch_get, err := SGBGIAPN(groupId, pluginName)
	if err != nil || switch_get.IsTrue == false {
		err = bot_switch_sync.PluginSwitch.SwitchCreate()
		return err
	}
	err = switch_get.PluginSwitch.SwitchUpdate(stop)
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

//SGBGIAS SwitchGetsByGroupIdAndStop
//bsss bot_switch_stops
//bss bot_switch_stop
func SGBGIAS(groupId int64) (bsss SwitchStop, err error) {
	bps := make([]string, 0)
	rows, err := Db.Query("select plugin_name from [kequ5060].[dbo].[zbot_plugin_switch] where group_id = $1 and stop = $2", groupId, true)
	if err != nil {
		return
	}
	defer rows.Close()
	var plu Plu
	for rows.Next() {
		rows.Scan(&plu.PluginName)
		bps = append(bps, plu.PluginName)
	}
	bsss = SwitchStop{
		groupId,
		bps,
	}
	return
}

//SGBGIAPN SwitchGetByGroupIdAndPluginName
func SGBGIAPN(groupId int64, pluginName string) (bot_switch_sync SwitchSync, err error) {
	bot_switch := Switch{}
	bot_switch_sync = SwitchSync{
		IsTrue:       true,
		PluginSwitch: &bot_switch,
	}
	bw := strconv.Itoa(int(groupId)) + "_" + pluginName
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
		err = Db.QueryRow("select ID, group_id, plugin_name, gmt_modified, stop from [kequ5060].[dbo].[zbot_plugin_switch] where group_id = $1 and plugin_name = $2", groupId, pluginName).Scan(&bot_switch_sync.PluginSwitch.Id, &bot_switch_sync.PluginSwitch.GroupId, &bot_switch_sync.PluginSwitch.PluginName, &bot_switch_sync.PluginSwitch.GmtModified, &bot_switch_sync.PluginSwitch.Stop)
		info := fmt.Sprintf("%s", err)
		if StartsWith(info, "sql") || StartsWith(info, "unable"){
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
