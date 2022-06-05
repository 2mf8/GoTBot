package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	. "github.com/2mf8/tbotGo/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/guregu/null.v3"
	_ "gopkg.in/guregu/null.v3/zero"
)

type Learn struct {
	Id          int64
	Ask         string
	GroupId     int64
	AdminId     int64
	Answer      null.String
	GmtModified time.Time
	//Pass        bool
}

type LearnSync struct{
	IsTrue bool
	LearnSync *Learn
}

func LearnGet(groupId int64, ask string) (learnSync LearnSync, err error) {
	learn := Learn{}
	learnSync = LearnSync{
		IsTrue: true,
		LearnSync: &learn,
	}

	var vb []byte
	var bw_set []byte

	bw := strconv.Itoa(int(groupId)) + "_" + ask
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-学习", bw)
		err = Db.QueryRow("select * from [kequ5060].[dbo].[zbot_learn] where group_id = $1 and ask = $2", groupId, ask).Scan(&learn.Id, &learn.Ask, &learn.GroupId, &learn.AdminId, &learn.Answer, &learn.GmtModified)
		info := fmt.Sprintf("%s", err)
		if StartsWith(info, "sql") || StartsWith(info, "unable"){
			if StartsWith(info, "unable") {
				fmt.Println(info)
			}
			learnSync = LearnSync{
				IsTrue: false,
				LearnSync: &learn,
			}
		}
		bw_set, _ = json.Marshal(&learnSync)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return
	}
	err = json.Unmarshal(vb, &learnSync)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	fmt.Println("[Redis] Key(", bw, ") Value(", learnSync.IsTrue, *learnSync.LearnSync, ")")  //测试用
	return
}

func (learn *Learn) LearnCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[zbot_learn] values ($1, $2, $3, $4, $5) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(learn.Ask, learn.GroupId, learn.AdminId, learn.Answer, learn.GmtModified).Scan(&learn.Id)

	learnSync := LearnSync{
		IsTrue: true,
		LearnSync: &Learn{
			Id: learn.Id,
			Ask: learn.Ask,
			GroupId: learn.GroupId,
			AdminId: learn.AdminId,
			Answer: learn.Answer,
			GmtModified: learn.GmtModified,
		},
	}

	bw := strconv.Itoa(int(learn.GroupId)) + "_" + learn.Ask
	var bw_set []byte
	bw_set, _ = json.Marshal(&learnSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)
	return
}

func (learn *Learn) LearnUpdate(answer null.String) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[zbot_learn] set ask = $2, group_id = $3, admin_id = $4, answer = $5, gmt_modified = $6 where ID = $1", learn.Id, learn.Ask, learn.GroupId, learn.AdminId, answer.String, learn.GmtModified)
	
	learnSync := LearnSync{
		IsTrue: true,
		LearnSync: &Learn{
			Id: learn.Id,
			Ask: learn.Ask,
			GroupId: learn.GroupId,
			AdminId: learn.AdminId,
			Answer: answer,
			GmtModified: learn.GmtModified,
		},
	}

	bw := strconv.Itoa(int(learn.GroupId)) + "_" + learn.Ask
	var bw_set []byte
	bw_set, _ = json.Marshal(&learnSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)

	return
}

func (learn *Learn) LearnDeleteByGroupIdAndAsk() (err error) {
	_, err = Db.Exec("delete from [kequ5060].[dbo].[zbot_learn] where group_id = $1 and ask = $2", learn.GroupId, learn.Ask)
	return
}

func LearnSave(ask string, groupId int64, adminId int64, answer null.String, gmtModified time.Time) (err error) {
	learn := Learn{
		Ask:         ask,
		GroupId:     groupId,
		AdminId:     adminId,
		Answer:      answer,
		GmtModified: gmtModified,
	}
	learn_get, err := LearnGet(groupId, ask)
	if err != nil || learn_get.IsTrue == false {
		err = learn.LearnCreate()
		return
	}
	err = learn_get.LearnSync.LearnUpdate(answer)
	return
}

func LDBGAA(groupId int64, ask string) (err error) {
	learn_get, err := LearnGet(groupId, ask)
	if err != nil {
		return
	}
	learn_get.LearnSync.LearnDeleteByGroupIdAndAsk()

	learnSync := LearnSync{
		IsTrue: true,
		LearnSync: &Learn{},
	}

	bw := strconv.Itoa(int(groupId)) + "_" + ask
	var bw_set []byte
	bw_set, _ = json.Marshal(&learnSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)
	return
}
