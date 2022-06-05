package data

import (
	"strconv"
	"fmt"
	"encoding/json"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
	_ "gopkg.in/guregu/null.v3/zero"
	. "github.com/2mf8/go-tbot-for-rq/public"
)

type Subscribe struct {
	Id             int64
	OriginGroupId  int64
	ReplaceGroupId int64
	AdminId        int64
}

type SubscribeSync struct {
	IsTrue        bool
	SubSync *Subscribe
}

func GetSubscribe(groupId int64) (sub_sync SubscribeSync, err error) {
	sub := Subscribe{}
	sub_sync = SubscribeSync{
		IsTrue: true,
		SubSync: &sub,
	}

	// byte write
	bw := strconv.Itoa(int(groupId)) + "_sub"
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()
	// value byte
	var vb []byte
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-魔友价", bw)
		err = Db.QueryRow("select * from [kequ5060].[dbo].[zbot_replace] where orgin_group_id = $1", groupId).Scan(&sub_sync.SubSync.Id, &sub_sync.SubSync.OriginGroupId, &sub_sync.SubSync.ReplaceGroupId, &sub_sync.SubSync.AdminId)
		info := fmt.Sprintf("%s", err)
		if StartsWith(info, "sql") || StartsWith(info, "unable"){
			if StartsWith(info, "unable") {
				fmt.Println(info)
			}
			sub_sync = SubscribeSync{
				IsTrue:       false,
				SubSync: &sub,
			}
		}
		var bw_set []byte
		bw_set, _ = json.Marshal(&sub_sync)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return
	}
	err = json.Unmarshal(vb, &sub_sync)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	//fmt.Println("[Redis] Key(", bw, ") Value(", sub_sync.IsTrue, *sub_sync.SubSync, ")")  //测试用
	return
}

func (sub *Subscribe) SubCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[zbot_replace] values ($1, $2, $3) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(sub.OriginGroupId, sub.ReplaceGroupId, sub.AdminId).Scan(&sub.Id)

	sub_sync := SubscribeSync{
		IsTrue: true,
		SubSync: &Subscribe{
			Id: sub.Id,
			OriginGroupId: sub.OriginGroupId,
			ReplaceGroupId: sub.ReplaceGroupId,
			AdminId: sub.AdminId,
		},
	}

	bw := strconv.Itoa(int(sub.OriginGroupId)) + "_sub"
	var sub_redis []byte
	sub_redis, err = json.Marshal(&sub_sync)
	if err != nil {
		fmt.Println("[错误] Marshal序列化出错")
	}
	c := Pool.Get()
	defer c.Close()
	c.Send("SET", bw, sub_redis)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Printf("[收到] %#v\n", v)

	return
}

func (sub *Subscribe) SubUpdate(replace_group_id int64) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[zbot_replace] set orgin_group_id = $2, replace_group_id = $3, admin_id = $4 where ID = $1", sub.Id, sub.OriginGroupId, replace_group_id, sub.AdminId)
	
	u_sub := Subscribe{
		Id: sub.Id,
		OriginGroupId: sub.OriginGroupId,
		ReplaceGroupId: replace_group_id,
		AdminId: sub.AdminId,
	}

	sub_sync := SubscribeSync{
		IsTrue: true,
		SubSync: &u_sub,
	}

	bw := strconv.Itoa(int(sub.OriginGroupId)) + "_sub"
	var sub_redis []byte
	sub_redis, err = json.Marshal(&sub_sync)
	if err != nil {
		fmt.Println("[错误] Marshal序列化出错")
	}
	c := Pool.Get()
	defer c.Close()
	c.Send("SET", bw, sub_redis)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Printf("[收到] %#v\n", v)

	return
}

func (sub *Subscribe) SubDelete() (err error) {
	_, err = Db.Exec("delete from [kequ5060].[dbo].[zbot_replace] where ID = $1", sub.Id)
	return
}

func SubSave(orgin_group_id int64, replace_group_id int64, admin_id int64) (err error) {
	sub := Subscribe{
		OriginGroupId:  orgin_group_id,
		ReplaceGroupId: replace_group_id,
		AdminId:        admin_id,
	}
	sub_get, err := GetSubscribe(orgin_group_id)
	if err != nil || sub_get.IsTrue == false {
		err = sub.SubCreate()
		return
	}
	err = sub_get.SubSync.SubUpdate(replace_group_id)
	return
}

func SubDeleteByGroupId(groupId int64) (err error) {
	sub_get, _ := GetSubscribe(groupId)
	sub_get.SubSync.SubDelete()
	
	sub_sync := SubscribeSync{
		IsTrue: false,
		SubSync: &Subscribe{},
	}

	bw := strconv.Itoa(int(groupId)) + "_sub"
	var sub_redis []byte
	sub_redis, err = json.Marshal(&sub_sync)
	if err != nil {
		fmt.Println("[错误] Marshal序列化出错")
	}
	c := Pool.Get()
	defer c.Close()
	c.Send("SET", bw, sub_redis)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Printf("[收到] %#v\n", v)

	return
}
