package database

import (
	"fmt"
	"time"

	"github.com/2mf8/GoTBot/config"
	"github.com/2mf8/GoTBot/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/tencent-connect/botgo/log"
	_ "gopkg.in/guregu/null.v3/zero"
)

type PBlock struct {
	Id          int       `json:"id"`
	GuildId     string    `json:"guild_id"`
	UserId      string    `json:"user_id"`
	IsPBlock    bool      `json:"ispblock"`
	AdminId     string    `json:"admin_id"`
	GmtModified time.Time `json:"gmt_modified"`
}

func PBlockGet(guildId, userId string) (p PBlock, err error) {
	pblock := PBlock{}
	statment := fmt.Sprintf("select ID, guild_id, user_id, admin_id, gmt_modified, ispblock from [%s].[dbo].[guild_pblock] where user_id = $1 and ispblock = $2 and guild_id = $3", config.Conf.DatabaseName)
	err = Db.QueryRow(statment, userId, true, guildId).Scan(&pblock.Id, &pblock.GuildId, &pblock.UserId, &pblock.AdminId, &pblock.GmtModified, &pblock.IsPBlock)
	info := fmt.Sprintf("%s", err)
	if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
		if public.StartsWith(info, "unable") {
			log.Warn(info)
		}
		return
	}
	return
}

func (pBlock *PBlock) PBlockCreate() (err error) {
	statement := fmt.Sprintf("insert into [%s].[dbo].[guild_pblock] values ($1, $2, $3, $4, $5) select @@identity", config.Conf.DatabaseName)
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(pBlock.GuildId, pBlock.UserId, pBlock.AdminId, pBlock.GmtModified, pBlock.IsPBlock).Scan(&pBlock.Id)
	return
}

func (pBlock *PBlock) PBlockUpdate() (err error) {
	statment := fmt.Sprintf("update [%s].[dbo].[guild_pblock] set guild_id = $2, user_id = $3, ispblock = $4, admin_id = $5, gmt_modified = $6 where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, pBlock.Id, pBlock.GuildId, pBlock.UserId, pBlock.IsPBlock, pBlock.AdminId, pBlock.GmtModified)
	return
}

func PBlockSave(guildId, userId, adminId string, ispblock bool, gmtModified time.Time) (err error) {
	pblock := PBlock{
		GuildId:     guildId,
		UserId:      userId,
		IsPBlock:    ispblock,
		AdminId:     adminId,
		GmtModified: gmtModified,
	}
	pblock_get, err := PBlockGet(guildId, userId)
	if err != nil {
		pblock.PBlockCreate()
		return
	}
	pblock_get.IsPBlock = ispblock
	pblock_get.PBlockUpdate()
	return
}
