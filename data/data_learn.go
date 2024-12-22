package database

import (
	"fmt"
	"time"

	"github.com/2mf8/GoTBot/config"
	"github.com/2mf8/GoTBot/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/2mf8/Better-Bot-Go/log"
	"gopkg.in/guregu/null.v3"
	_ "gopkg.in/guregu/null.v3/zero"
)

type Learn struct {
	Id          int64
	Ask         string
	GuildId     string
	ChannelId   string
	AdminId     string
	Answer      null.String
	GmtModified time.Time
	Pass        bool
}
type TLearn struct {
	Id          int64
	Ask         string
	GuildId     string
	ChannelId   null.String
	AdminId     string
	Answer      null.String
	GmtModified time.Time
	Pass        bool
}

func LearnUpdateZbot() {
	statment := fmt.Sprintf("select ID, ask, guild_id, channel_id, admin_id, answer, gmt_modified, pass from [%s].[dbo].[zbot_learn]", config.Conf.DatabaseName)
	rows, err := Db.Query(statment)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		l := TLearn{}
		_ = rows.Scan(&l.Id, &l.Ask, &l.GuildId, &l.ChannelId, &l.AdminId, &l.Answer, &l.GmtModified, &l.Pass)
		l.ChannelId = null.NewString(l.GuildId, true)
		l.LearnUpdateZbot()
	}
}

func (learn *TLearn) LearnUpdateZbot() (err error) {
	statment := fmt.Sprintf("update [%s].[dbo].[zbot_learn] set ask = $2, guild_id = $3, channel_id = $4, admin_id = $5, answer = $6, gmt_modified = $7, pass = $8 where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, learn.Id, learn.Ask, learn.GuildId, learn.ChannelId, learn.AdminId, learn.Answer.String, learn.GmtModified, learn.Pass)
	return
}

func LearnGet(guildId, channelId, ask string) (l Learn, err error) {
	l = Learn{}
	statment := fmt.Sprintf("select ID, ask, guild_id, channel_id, admin_id, answer, gmt_modified, pass from [%s].[dbo].[guild_learn] where guild_id = $1 and ask = $2 and channel_id = $3", config.Conf.DatabaseName)
	err = Db.QueryRow(statment, guildId, ask, channelId).Scan(&l.Id, &l.Ask, &l.GuildId, &l.ChannelId, &l.AdminId, &l.Answer, &l.GmtModified, &l.Pass)
	info := fmt.Sprintf("%s", err)
	if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
		if public.StartsWith(info, "unable") {
			log.Warn(info)
		}
	}
	return
}

func (learn *Learn) LearnCreate() (err error) {
	statement := fmt.Sprintf("insert into [%s].[dbo].[guild_learn] values ($1, $2, $3, $4, $5, $6, $7) select @@identity", config.Conf.DatabaseName)
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(learn.Ask, learn.GuildId, learn.ChannelId, learn.AdminId, learn.Answer, learn.GmtModified, learn.Pass).Scan(&learn.Id)
	return
}

func (learn *Learn) LearnUpdate() (err error) {
	statment := fmt.Sprintf("update [%s].[dbo].[guild_learn] set ask = $2, guild_id = $3, channel_id = $4, admin_id = $5, answer = $6, gmt_modified = $7, pass = $8 where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, learn.Id, learn.Ask, learn.GuildId, learn.ChannelId, learn.AdminId, learn.Answer.String, learn.GmtModified, learn.Pass)
	return
}

// LearnDeleteByGuildIdAndAskAndChannelId
func (learn *Learn) LDBGIAAACI() (err error) {
	statment := fmt.Sprintf("delete from [%s].[dbo].[guild_learn] where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, learn.Id)
	return
}

func LearnSave(ask, guildId, channelId, adminId string, answer null.String, gmtModified time.Time, pass bool) (err error) {
	learn := Learn{
		Ask:         ask,
		GuildId:     guildId,
		ChannelId:   channelId,
		AdminId:     adminId,
		Answer:      answer,
		GmtModified: gmtModified,
		Pass:        pass,
	}
	learn_get, err := LearnGet(guildId, channelId, ask)
	if err != nil {
		err = learn.LearnCreate()
		return
	}
	learn_get.Answer = answer
	err = learn_get.LearnUpdate()
	return
}

func LDBGAA(guildId, channelId, ask string) (err error) {
	learn_get, err := LearnGet(guildId, channelId, ask)
	if err != nil {
		return
	}
	learn_get.LDBGIAAACI()
	return
}
