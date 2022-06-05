package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	. "github.com/2mf8/tbotGo/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
	_ "gopkg.in/guregu/null.v3/zero"
)

type PBlock struct {
	Id          int
	UserId      int64
	IsPBlock    bool
	AdminId     int64
	GmtModified time.Time
}

type PBlockSync struct {
	IsTrue     bool
	PBlockSync *PBlock
}

func PBlockGet(userId int64) (pblockSync PBlockSync, err error) {
	pblock := PBlock{}
	pblockSync = PBlockSync{
		IsTrue: true,
		PBlockSync: &pblock,
	}

	bw := "pblock_" + strconv.Itoa(int(userId))
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()

	var vb []byte
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-个人屏蔽", bw)
		err = Db.QueryRow("select * from [kequ5060].[dbo].[zbot_pblock] where user_id = $1 and ispblock = $2", userId, true).Scan(&pblock.Id, &pblock.UserId, &pblock.AdminId, &pblock.GmtModified, &pblock.IsPBlock)
		info := fmt.Sprintf("%s", err)
		if StartsWith(info, "sql") || StartsWith(info, "unable") {
			if StartsWith(info, "unable") {
				fmt.Println(info)
			}
			pblockSync = PBlockSync{
				IsTrue:     false,
				PBlockSync: &pblock,
			}
		}
		var bw_set []byte
		bw_set, _ = json.Marshal(&pblockSync)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return
	}
	err = json.Unmarshal(vb, &pblockSync)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	//fmt.Println("[Redis] Key(", bw, ") Value(", pblockSync.IsTrue, *pblockSync.PBlockSync, ")") //测试用
	return
}

func (pBlock *PBlock) PBlockCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[zbot_pblock] (user_id, admin_id, gmt_modified, ispblock) values ($1, $2, $3, $4) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(pBlock.UserId, pBlock.AdminId, pBlock.GmtModified, pBlock.IsPBlock).Scan(&pBlock.Id)

	pblockSync := PBlockSync{
		IsTrue: true,
		PBlockSync: &PBlock{
			Id:          pBlock.Id,
			UserId:      pBlock.UserId,
			IsPBlock:    pBlock.IsPBlock,
			AdminId:     pBlock.AdminId,
			GmtModified: pBlock.GmtModified,
		},
	}

	bw := "pblock_" + strconv.Itoa(int(pBlock.UserId))
	var bw_set []byte
	bw_set, _ = json.Marshal(&pblockSync)
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

func (pBlock *PBlock) PBlockUpdate(ispblock bool) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[zbot_pblock] set user_id = $2, ispblock = $3, admin_id = $4, gmt_modified = $5 where ID = $1", pBlock.Id, pBlock.UserId, pBlock.IsPBlock, pBlock.AdminId, pBlock.GmtModified)

	pblockSync := PBlockSync{
		IsTrue: true,
		PBlockSync: &PBlock{
			Id:          pBlock.Id,
			UserId:      pBlock.UserId,
			IsPBlock:    ispblock,
			AdminId:     pBlock.AdminId,
			GmtModified: pBlock.GmtModified,
		},
	}

	bw := "pblock_" + strconv.Itoa(int(pBlock.UserId))
	var bw_set []byte
	bw_set, _ = json.Marshal(&pblockSync)
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

func PBlockSave(userId int64, ispblock bool, adminId int64, gmtModified time.Time) (err error) {
	pblock := PBlock{
		UserId:      userId,
		IsPBlock:    ispblock,
		AdminId:     adminId,
		GmtModified: gmtModified,
	}
	pblock_get, err := PBlockGet(userId)
	if err != nil || pblock_get.IsTrue == false {
		pblock.PBlockCreate()
		return
	}
	pblock_get.PBlockSync.PBlockUpdate(ispblock)
	return
}
