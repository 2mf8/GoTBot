package database

import (
	"fmt"

	"github.com/2mf8/GoTBot/config"
	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/guregu/null.v3"
	_ "gopkg.in/guregu/null.v3/zero"
)

type Achievement struct {
	Id       int64
	UserId   string
	UserName string
	Avatar   null.String
	Item     string
	Best     int
	Average  int
	Session  int
}

func (a *Achievement) AchievementCreate() (err error) {
	statement := fmt.Sprintf("insert into [%s].[dbo].[guild_achievement] values ($1, $2, $3, $4, $5, $6, $7) select @@identity", config.Conf.DatabaseName)
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(a.UserId, a.UserName, a.Avatar, a.Item, a.Best, a.Average, a.Session).Scan(&a.Id)
	return
}

func (a *Achievement) AchievementUpdate() (err error) {
	statement := fmt.Sprintf("update [%s].[dbo].[guild_achievement] set user_id = $2, user_name = $3, avatar = $8, item = $4, best = $5, average = $6, session = $7 where Id = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statement, a.Id, a.UserId, a.UserName, a.Item, a.Best, a.Average, a.Session, a.Avatar)
	return
}

func AchievementSave(userId, userName string, avatar null.String, item string, best, average, session int) (err error) {
	a := Achievement{
		UserId:   userId,
		UserName: userName,
		Avatar:   avatar,
		Item:     item,
		Best:     best,
		Average:  average,
		Session:  session,
	}
	a_get, err := AchievementGet(userId, item, session)
	if err != nil {
		err = a.AchievementCreate()
		return
	}
	if a_get.Best == -1 && a_get.Average == -1 {
		a_get.UserName = userName
		a_get.Avatar = avatar
		a_get.Best = best
		a_get.Average = average
		err = a_get.AchievementUpdate()
		return
	}
	if a_get.Best == -1 {
		a_get.UserName = userName
		a_get.Avatar = avatar
		a_get.Best = best
		err = a_get.AchievementUpdate()
		return
	}
	if a_get.Average == -1 {
		a_get.UserName = userName
		a_get.Avatar = avatar
		a_get.Average = average
		err = a_get.AchievementUpdate()
		return
	}
	if best < a_get.Best && average < a_get.Average {
		a_get.UserName = userName
		a_get.Avatar = avatar
		a_get.Best = best
		a_get.Average = average
		err = a_get.AchievementUpdate()
		return
	}
	if best < a_get.Best {
		a_get.UserName = userName
		a_get.Avatar = avatar
		a_get.Best = best
		err = a_get.AchievementUpdate()
		return
	}
	if average < a_get.Average {
		a_get.UserName = userName
		a_get.Avatar = avatar
		a_get.Average = average
		err = a_get.AchievementUpdate()
		return
	}
	a_get.UserName = userName
	a_get.Avatar = avatar
	err = a_get.AchievementUpdate()
	return
}

func AchievementGet(userId, item string, session int) (a Achievement, err error) {
	a = Achievement{}
	statment := fmt.Sprintf("select Id, user_id, user_name, avatar, item, best, average, session from [%s].[dbo].[guild_achievement] where user_id = $1 and session = $3 and item = $2", config.Conf.DatabaseName)
	err = Db.QueryRow(statment, userId, item, session).Scan(&a.Id, &a.UserId, &a.UserName, &a.Avatar, &a.Item, &a.Best, &a.Average, &a.Session)
	return
}

// AchievementDeleteByUserIdAndSession
func ADBUAS(userId string, session int) (err error) {
	statment := fmt.Sprintf("delete from [%s].[dbo].[guild_achievement] where user_id = $1 and session = $2", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, userId, session)
	return
}

// AchievementDeleteByUserIdAndItemAndSession
func ADBUAIAS(userId, item string, session int) (err error) {
	statment := fmt.Sprintf("delete from [%s].[dbo].[guild_achievement] where user_id = $1 and item = $2 and session = $3", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, userId, item, session)
	return
}

// AchievementGetByUserIdAndSession
func AGBUAS(userId string, session int) (as []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, avatar, item, best, average, session from [%s].[dbo].[guild_achievement] where user_id = '%s' and session = %d and best > -1", config.Conf.DatabaseName, userId, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Avatar, &a.Item, &a.Best, &a.Average, &a.Session)
		as = append(as, a)
	}
	return
}

// AchievementGetByItemAndSessionOrderByBestAsc
// desc 大 → 小
// asc 小 → 大
func AGBIASOBBA(item string, session int) (bs []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, avatar, item, best, average, session from [%s].[dbo].[guild_achievement] where item = '%s' and session = %d and best > -1 order by best asc", config.Conf.DatabaseName, item, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := Achievement{}
		err = rows.Scan(&b.Id, &b.UserId, &b.UserName, &b.Avatar, &b.Item, &b.Best, &b.Average, &b.Session)
		bs = append(bs, b)
	}
	return
}

// AchievementGetByItemAndSessionOrderByAverageAsc
// desc 大 → 小
// asc 小 → 大
func AGBIASOBAA(item string, session int) (as []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, avatar, item, best, average, session from [%s].[dbo].[guild_achievement] where item = '%s' and session = %d and average > -1 order by average asc", config.Conf.DatabaseName, item, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Avatar, &a.Item, &a.Best, &a.Average, &a.Session)
		as = append(as, a)
	}
	return
}

// AchievementGetBySessionOrderByItemAscAndBestAsc
// desc 大 → 小
// asc 小 → 大
func AGBSOBIAABA(session int) (as []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, avatar, item, best, average, session from [%s].[dbo].[guild_achievement] where session = %d and best > -1 order by item asc, best asc", config.Conf.DatabaseName, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Avatar, &a.Item, &a.Best, &a.Average, &a.Session)
		as = append(as, a)
	}
	return
}

// AchievementGetBySessionOrderByItemAscAndAverageAsc
// desc 大 → 小
// asc 小 → 大
func AGBSOBIAAAA(session int) (as []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, avatar, item, best, average, session from [%s].[dbo].[guild_achievement] where session = %d and average > -1 order by item asc, average asc", config.Conf.DatabaseName, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Avatar, &a.Item, &a.Best, &a.Average, &a.Session)
		as = append(as, a)
	}
	return
}

func AchievementGetCount(item string, best, average, session int) (i, j int, err error) {
	i = 0
	j = 0
	statment := fmt.Sprintf("select Id, user_id, user_name, avatar, item, best, average, session from [%s].[dbo].[guild_achievement] where item = '%s' and session = %d", config.Conf.DatabaseName, item, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Avatar, &a.Item, &a.Best, &a.Average, &a.Session)
		if a.Best > -1 && a.Best < best {
			i++
		}
		if a.Average > -1 && a.Average < average {
			j++
		}
	}
	return
}

func BestAndAverageTimeConvert(b, a int) (bc, ac string) {
	bt := "DNF"
	at := "DNF"
	bm := b / 60000
	bs := b % 60000 / 1000
	bms := b % 60000 % 1000
	am := a / 60000
	as := a % 60000 / 1000
	ams := a % 60000 % 1000
	if b > -1 && bm == 0 {
		if bms < 10 {
			bt = fmt.Sprintf("%d.00%d", bs, bms)
		} else if bms < 100 {
			bt = fmt.Sprintf("%d.0%d", bs, bms)
		} else {
			bt = fmt.Sprintf("%d.%d", bs, bms)
		}
	}
	if bm > 0 {
		if bs < 10 {
			if bms < 10 {
				bt = fmt.Sprintf("%d:0%d.00%d", bm, bs, bms)
			} else if bms < 100 {
				bt = fmt.Sprintf("%d:0%d.0%d", bm, bs, bms)
			} else {
				bt = fmt.Sprintf("%d:0%d.%d", bm, bs, bms)
			}
		} else {
			if bms < 10 {
				bt = fmt.Sprintf("%d:%d.00%d", bm, bs, bms)
			} else if bms < 100 {
				bt = fmt.Sprintf("%d:%d.0%d", bm, bs, bms)
			} else {
				bt = fmt.Sprintf("%d:%d.%d", bm, bs, bms)
			}
		}
	}

	if a > -1 && am == 0 {
		if ams < 10 {
			at = fmt.Sprintf("%d.00%d", as, ams)
		} else if ams < 100 {
			at = fmt.Sprintf("%d.0%d", as, ams)
		} else {
			at = fmt.Sprintf("%d.%d", as, ams)
		}
	}
	if am > 0 {
		if as < 10 {
			if ams < 10 {
				at = fmt.Sprintf("%d:0%d.00%d", am, as, ams)
			} else if ams < 100 {
				at = fmt.Sprintf("%d:0%d.0%d", am, as, ams)
			} else {
				at = fmt.Sprintf("%d:0%d.%d", am, as, ams)
			}
		} else {
			if ams < 10 {
				at = fmt.Sprintf("%d:%d.00%d", am, as, ams)
			} else if ams < 100 {
				at = fmt.Sprintf("%d:%d.0%d", am, as, ams)
			} else {
				at = fmt.Sprintf("%d:%d.%d", am, as, ams)
			}
		}
	}
	return bt, at
}
