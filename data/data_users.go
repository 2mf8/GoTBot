package database

import (
	"fmt"

	"github.com/2mf8/GoTBot/config"
	"github.com/2mf8/GoTBot/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/2mf8/Better-Bot-Go/log"
	"gopkg.in/guregu/null.v3"
	_ "gopkg.in/guregu/null.v3/zero"
)

type UserInfo struct {
	Id            int64       `json:"id"`
	UserId        null.String `json:"user_id"`
	Username      null.String `json:"user_name"`
	UserRole      int         `json:"user_role"` // 1<<1 黄小姐 1<<2 奇乐 1<<30 系统
	UserAvatar    null.String `json:"user_avatar"`
	ServerNumber  null.String `json:"server_number"`
	Password      null.String `json:"password"`
	Email         null.String `josn:"email"`
	QQUnionId     null.String `json:"qq_union_id"`
	WeixinUnionId null.String `json:"weixin_union_id"`
}

func UserInfoGet(userid, email, qq_union_id, weixin_union_id string) (ui UserInfo, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, user_role, user_avatar, server_number, password, email, qq_union_id, weixin_union_id from [%s].[dbo].[user_info] where user_id = $1 or email = $2 or qq_union_id = $3 or weixin_union_id = $4", config.Conf.DatabaseName)
	err = Db.QueryRow(statment, userid, email, qq_union_id, weixin_union_id).Scan(&ui.Id, &ui.UserId, &ui.Username, &ui.UserRole, &ui.UserAvatar, &ui.ServerNumber, &ui.Password, &ui.Email, &ui.QQUnionId, &ui.WeixinUnionId)
	info := fmt.Sprintf("%s", err)
	if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
		if public.StartsWith(info, "unable") {
			log.Warn(info)
		}
	}
	return
}

func UserInfoGetByBot(userid string) (ui UserInfo, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, user_role, user_avatar, server_number, password, email, qq_union_id, weixin_union_id from [%s].[dbo].[user_info] where user_id = $1", config.Conf.DatabaseName)
	err = Db.QueryRow(statment, userid).Scan(&ui.Id, &ui.UserId, &ui.Username, &ui.UserRole, &ui.UserAvatar, &ui.ServerNumber, &ui.Password, &ui.Email, &ui.QQUnionId, &ui.WeixinUnionId)
	info := fmt.Sprintf("%s", err)
	if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
		if public.StartsWith(info, "unable") {
			log.Warn(info)
		}
	}
	return
}

func (ui *UserInfo) UserInfoUpdate() error {
	statment := fmt.Sprintf("update [%s].[dbo].[user_info] set user_id = $2, user_name = $3, user_role = $4, user_avatar = $5, server_number = $6, password = $7, email = $8, qq_union_id = $9, weixin_union_id = $10 where Id = $1", config.Conf.DatabaseName)
	_, err := Db.Exec(statment, ui.Id, ui.UserId.String, ui.Username.String, ui.UserRole, ui.UserAvatar.String, ui.ServerNumber.String, ui.Password.String, ui.Email.String, ui.QQUnionId.String, ui.WeixinUnionId.String)
	return err
}

// user_info
func (ui *UserInfo) UserInfoCreate() error {
	statment := fmt.Sprintf("insert into [%s].[dbo].[user_info] values ($1 ,$2 ,$3 ,$4 ,$5 ,$6 ,$7 ,$8 ,$9) select @@identity", config.Conf.DatabaseName)
	stmt, err := Db.Prepare(statment)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ui.UserId.String, ui.Username.String, ui.UserRole, ui.UserAvatar.String, ui.ServerNumber.String, ui.Password.String, ui.Email.String, ui.QQUnionId.String, ui.WeixinUnionId.String).Scan(&ui.Id)
	return err
}

func UserInfoSave(user_id, user_name, user_avatar, server_number, password, email, qq_union_id, weixin_union_id null.String, user_role int) (err error) {
	ui := UserInfo{
		UserId:        user_id,
		Username:      user_name,
		UserRole:      user_role,
		UserAvatar:    user_avatar,
		ServerNumber:  server_number,
		Password:      password,
		Email:         email,
		QQUnionId:     qq_union_id,
		WeixinUnionId: weixin_union_id,
	}
	ui_get, err := UserInfoGetByBot(user_id.String)
	if err != nil {
		err = ui.UserInfoCreate()
		return
	}
	if user_id.String != "" {
		ui_get.UserId = user_id
	}
	if user_name.String != "" {
		ui_get.Username = user_name
	}
	if user_avatar.String != "" {
		ui_get.UserAvatar = user_avatar
	}
	if server_number.String != "" {
		ui_get.ServerNumber = server_number
	}
	if password.String != "" {
		ui_get.Password = password
	}
	if email.String != "" {
		ui_get.Email = email
	}
	if qq_union_id.String != "" {
		ui_get.QQUnionId = qq_union_id
	}
	if weixin_union_id.String != "" {
		ui_get.WeixinUnionId = weixin_union_id
	}
	if user_role > 0 {
		ui_get.UserRole = user_role
	}
	err = ui_get.UserInfoUpdate()
	return
}

func (ui *UserInfo) UserInfoDelete() error {
	statment := fmt.Sprintf("delete from [%s].[dbo].[user_info] where ID = $1", config.Conf.DatabaseName)
	_, err := Db.Exec(statment, ui.Id)
	return err
}

func UID(userid, email, qq_union_id, weixin_union_id string) error {
	uid, err := UserInfoGet(userid, email, qq_union_id, weixin_union_id)
	if err != nil {
		return err
	}
	err = uid.UserInfoDelete()
	return err
}
