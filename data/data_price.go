package database

import (
	"fmt"
	"time"

	"github.com/2mf8/GoTBot/config"
	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/guregu/null.v3"
	_ "gopkg.in/guregu/null.v3/zero"
)

type CuberPrice struct {
	Id            int64       `json:"id"`
	GuildId       string      `json:"guild_id"`
	ChannelId     string      `json:"channel_id"`
	Brand         null.String `json:"brand"`
	Item          string      `json:"item"`
	Price         null.String `json:"price"`
	Shipping      null.String `json:"shipping"`
	Updater       null.String `json:"updater"`
	GmtModified   int64       `json:"gmt_modified"`
	IsMagnetism   bool        `json:"is_magnetism"`
	MagnetismType null.String `json:"magnetism_type"`
}

type TCuberPrice struct {
	Id            int64       `json:"id"`
	GuildId       string      `json:"guild_id"`
	ChannelId     string      `json:"channel_id"`
	Brand         null.String `json:"brand"`
	Item          string      `json:"item"`
	Price         null.String `json:"price"`
	Shipping      null.String `json:"shipping"`
	Updater       null.String `json:"updater"`
	GmtModified   null.Time   `json:"gmt_modified"`
	IsMagnetism   bool        `json:"is_magnetism"`
	MagnetismType null.String `json:"magnetism_type"`
}

// is_magnetism
func GetItem(guildId, channelId string, item string) (cp CuberPrice, err error) {
	statment := ""
	cp = CuberPrice{}
	statment = fmt.Sprintf("select ID, guild_id, channel_id, brand, item, price, shipping, updater, gmt_modified, is_magnetism, magnetism_type from [%s].[dbo].[guild_shop] where guild_id = $1 and channel_id = $3 and item = $2", config.Conf.DatabaseName)
	err = Db.QueryRow(statment, guildId, item, channelId).Scan(&cp.Id, &cp.GuildId, &cp.ChannelId, &cp.Brand, &cp.Item, &cp.Price, &cp.Shipping, &cp.Updater, &cp.GmtModified, &cp.IsMagnetism, &cp.MagnetismType)
	return
}

func GetItems(guildId, channelId string, key string) (cps []CuberPrice, err error) {
	statment := fmt.Sprintf("select ID, guild_id, channel_id, brand, item, price, shipping, updater, gmt_modified, is_magnetism, magnetism_type from [%s].[dbo].[guild_shop] where guild_id = '%s' and channel_id = '%s' and item like '%%%s%%'", config.Conf.DatabaseName, guildId, channelId, key)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		cp := CuberPrice{}
		err = rows.Scan(&cp.Id, &cp.GuildId, &cp.ChannelId, &cp.Brand, &cp.Item, &cp.Price, &cp.Shipping, &cp.Updater, &cp.GmtModified, &cp.IsMagnetism, &cp.MagnetismType)
		cps = append(cps, cp)
	}
	return
}

func GetItemsAll(guildId, channelId string) (cps []CuberPrice, err error) {
	statment := fmt.Sprintf("select ID, guild_id, channel_id, brand, item, price, shipping, updater, gmt_modified, is_magnetism, magnetism_type from [%s].[dbo].[guild_shop] where guild_id = '%s' and channel_id = '%s'", config.Conf.DatabaseName, guildId, channelId)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		cp := CuberPrice{}
		err = rows.Scan(&cp.Id, &cp.GuildId, &cp.ChannelId, &cp.Brand, &cp.Item, &cp.Price, &cp.Shipping, &cp.Updater, &cp.GmtModified, &cp.IsMagnetism, &cp.MagnetismType)
		cps = append(cps, cp)
	}
	return
}

func (cp *CuberPrice) ItemCreate() (id int64, err error) {
	statement := fmt.Sprintf("insert into [%s].[dbo].[guild_shop] values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) select @@identity", config.Conf.DatabaseName)
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(cp.GuildId, cp.ChannelId, cp.Brand, cp.Item, cp.Price, cp.Shipping, cp.Updater, cp.GmtModified, cp.IsMagnetism, cp.MagnetismType).Scan(&cp.Id)
	return cp.Id, err
}

func (cp *CuberPrice) ItemUpdate() (err error) {
	statment := fmt.Sprintf("update [%s].[dbo].[guild_shop] set guild_id = $2, channel_id = $9, brand = $3, item = $4, price = $5, shipping = $6, updater = $7, gmt_modified = $8, is_magnetism = $10, magnetism_type = $11 where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, cp.Id, cp.GuildId, cp.Brand, cp.Item, cp.Price.String, cp.Shipping.String, cp.Updater, cp.GmtModified, cp.ChannelId, cp.IsMagnetism, cp.MagnetismType)
	return
}

func (cp *CuberPrice) ItemDeleteById() (err error) {
	statment := fmt.Sprintf("delete from [%s].[dbo].[guild_shop] where ID = $1", config.Conf.DatabaseName)
	_, err = Db.Exec(statment, cp.Id)
	return
}

func ItemSave(guildId, channelId string, brand null.String, item string, price null.String, shipping null.String, updater null.String, gmtModified int64, is_magnetism bool, magnetism_type null.String) (id int64, err error) {
	cp := CuberPrice{
		GuildId:       guildId,
		ChannelId:     channelId,
		Brand:         brand,
		Item:          item,
		Price:         price,
		Shipping:      shipping,
		Updater:       updater,
		GmtModified:   gmtModified,
		IsMagnetism:   is_magnetism,
		MagnetismType: magnetism_type,
	}
	cp_get, err := GetItem(guildId, channelId, item)
	if err != nil {
		id, err = cp.ItemCreate()
		fmt.Println(id, err)
		return
	}
	cp_get.Price = price
	if shipping.String != "" {
		cp_get.Shipping = shipping
	}
	cp_get.IsMagnetism = is_magnetism
	err = cp_get.ItemUpdate()
	fmt.Println(err)
	return cp_get.Id, err
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

func GetAll() (err error) {
	ii := 0
	statment := fmt.Sprintf("select guild_id, channel_id, brand, item, price, shipping, updater, gmt_modified, is_magnetism, magnetism_type from [%s].[dbo].[guild_shop]", config.Conf.DatabaseName)
	rows, err := Db.Query(statment)
	fmt.Println(err)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println(ii)
		cp := TCuberPrice{}
		err = rows.Scan(&cp.GuildId, &cp.ChannelId, &cp.Brand, &cp.Item, &cp.Price, &cp.Shipping, &cp.Updater, &cp.GmtModified, &cp.IsMagnetism, &cp.MagnetismType)
		//if strings.Contains(cp.Item, "磁") {
		//cp.IsMagnetism = true
		ItemSave(cp.GuildId, cp.ChannelId, cp.Brand, cp.Item, cp.Price, cp.Shipping, cp.Updater, time.Now().Unix(), cp.IsMagnetism, cp.MagnetismType)
		//}
		ii++
	}
	return
}
