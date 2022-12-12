package data

import (
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/guregu/null.v3"
	_ "gopkg.in/guregu/null.v3/zero"
)

type CuberPrice struct {
	Id          int64       `json:"id"`
	GuildId     string      `json:"guild_id"`
	ChannelId   string      `json:"channel_id"`
	Brand       null.String `json:"brand"`
	Item        string      `json:"item"`
	Price       null.String `json:"price"`
	Shipping    null.String `json:"shipping"`
	Updater     string      `json:"updater"`
	GmtModified null.Time   `json:"gmt_modified"`
}

func GetItem(guildId, channelId string, item string) (cp CuberPrice, err error) {
	cp = CuberPrice{}
	err = Db.QueryRow("select ID, guild_id, channel_id, brand, item, price, shipping, updater, gmt_modified from [kequ5060].[dbo].[guild_price] where guild_id = $1 and channel_id = $3 and item = $2", guildId, item, channelId).Scan(&cp.Id, &cp.GuildId, &cp.ChannelId, &cp.Brand, &cp.Item, &cp.Price, &cp.Shipping, &cp.Updater, &cp.GmtModified)
	return
}

func GetItems(guildId, channelId string, key string) (cps []CuberPrice, err error) {
	statment := fmt.Sprintf("select ID, guild_id, channel_id, brand, item, price, shipping, updater, gmt_modified from [kequ5060].[dbo].[guild_price] where guild_id = %s and channel_id = %s and item like '%%%s%%'", guildId, channelId, key)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		cp := CuberPrice{}
		err = rows.Scan(&cp.Id, &cp.GuildId, &cp.ChannelId, &cp.Brand, &cp.Item, &cp.Price, &cp.Shipping, &cp.Updater, &cp.GmtModified)
		cps = append(cps, cp)
	}
	return
}

func (cp *CuberPrice) ItemCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[guild_price] values ($1, $2, $3, $4, $5, $6, $7, $8) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(cp.GuildId, cp.ChannelId, cp.ChannelId, cp.Brand, cp.Item, cp.Price, cp.Shipping, cp.Updater, cp.GmtModified).Scan(&cp.Id)
	return
}

func (cp *CuberPrice) ItemUpdate(price null.String, shipping null.String) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[guild_price] set guild_id = $2, channel_id = $9, brand = $3, item = $4, price = $5, shipping = $6, updater = $7, gmt_modified = $8 where ID = $1", cp.Id, cp.GuildId, cp.Brand, cp.Item, price.String, shipping.String, cp.Updater, cp.GmtModified, cp.ChannelId)
	return
}

func (cp *CuberPrice) ItemDeleteById() (err error) {
	_, err = Db.Exec("delete from [kequ5060].[dbo].[guild_price] where ID = $1", cp.Id)
	return
}

func ItemSave(guildId, channelId string, brand null.String, item string, price null.String, shipping null.String, updater string, gmtModified null.Time) (err error) {
	cp := CuberPrice{
		GuildId:     guildId,
		ChannelId:   channelId,
		Brand:       brand,
		Item:        item,
		Price:       price,
		Shipping:    shipping,
		Updater:     updater,
		GmtModified: gmtModified,
	}
	cp_get, err := GetItem(guildId, channelId, item)
	if err != nil {
		err = cp.ItemCreate()
		return
	}
	err = cp_get.ItemUpdate(price, shipping)
	return
}

// ItemDeleteByGuildIdAndName
func IDBGAN(guildId, channelId, item string) (err error) {
	cp_get_d, err := GetItem(guildId, channelId, item)
	if err != nil {
		return
	}
	err = cp_get_d.ItemDeleteById()
	return
}
