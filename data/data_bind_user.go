package database

import (
	"encoding/json"
	"fmt"

	"github.com/2mf8/GoTBot/config"
	"github.com/2mf8/GoTBot/public"
	"github.com/tencent-connect/botgo/log"
	"gopkg.in/guregu/null.v3"
)

type BindUserInfo struct {
	Id       int64
	UserId   string
	UserName null.String
	Avatar   null.String
	Uin      int64
}

type BindUserInfoSync struct {
	IsTrue   bool
	UserInfo *BindUserInfo
}

func BindUserInfoGet(userId string, uin int64) (buis BindUserInfoSync, err error) {
	bs := BindUserInfoSync{}
	key := fmt.Sprintf("%s_%v", userId, uin)
	bi, err := RedisGet(key)
	if err != nil {
		b := BindUserInfo{}
		statment := fmt.Sprintf("select ID, user_id, user_name, avatar, uin from [%s].[dbo].[bind_userinfo] where user_id = $1 or uin = $2", config.Conf.DatabaseName)
		err = Db.QueryRow(statment, userId, uin).Scan(&b.Id, &b.UserId, &b.UserName, &b.Avatar, &b.Uin)
		info := fmt.Sprintf("%s", err)
		if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
			if public.StartsWith(info, "unable") {
				log.Warn(info)
			}
			return bs, err
		}
		bs = BindUserInfoSync{
			IsTrue:   false,
			UserInfo: &b,
		}
		return bs, nil
	}
	json.Unmarshal(bi, &bs)
	return bs, nil
}

func (b *BindUserInfo) BindUserInfoCreate() (err error) {
	statement := fmt.Sprintf("insert into [%s].[dbo].[bind_userinfo] values ($1, $2, $3, $4) select @@identity", config.Conf.DatabaseName)
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(b.UserId, b.UserName, b.Avatar, b.Uin).Scan(&b.Id)
	key := fmt.Sprintf("%s_%v", b.UserId, b.Uin)
	bs := BindUserInfoSync{
		IsTrue:   true,
		UserInfo: b,
	}
	bi, err := json.Marshal(bs)
	if err != nil {
		return
	}
	RedisSet(key, bi)
	return
}

func (b *BindUserInfo) BindUserInfoUpdate() (err error) {
	statment := fmt.Sprintf("update [%s].[dbo].[bind_userinfo] set user_id = $2, user_name = $3, avatar = $4, uin = $5 where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, b.Id, b.UserId, b.UserName, b.Avatar, b.Uin)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s_%v", b.UserId, b.Uin)
	bs := BindUserInfoSync{
		IsTrue:   true,
		UserInfo: b,
	}
	bi, err := json.Marshal(bs)
	if err != nil {
		return
	}
	RedisSet(key, bi)
	return
}

func BindUserInfoSave(userId string, userName, avatar null.String, uin int64) (err error) {
	b := BindUserInfo{
		UserId:   userId,
		UserName: userName,
		Avatar:   avatar,
		Uin:      uin,
	}
	bs := BindUserInfoSync{
		IsTrue:   true,
		UserInfo: &b,
	}
	bsg, err := BindUserInfoGet(userId, uin)
	if err != nil || !bsg.IsTrue  {
		err = b.BindUserInfoCreate()
		return
	}
	bs.UserInfo = &b
	err = bs.UserInfo.BindUserInfoUpdate()
	return
}
