package plugins

import "github.com/2mf8/GoTBot/utils"

var PluginMap = map[string]utils.Plugin{
	"群管": &Admin{},
	"Bind": &Bind{},
	"屏蔽": &Block{},
	"赛季": &Competition{},
	"守卫": &Guard{},
	"查价": &PricePlugin{},
	"学习": &LearnPlugin{},
	"Log": &Log{},
	"Rank": &RankPlugin{},
	"复读": &Repeat{},
	"回复": &Reply{},
	"打乱": &ScramblePlugin{},
	"订阅": &Sub{},
	"开关": &BotSwitch{},
	"WCA": &WCA{},
}
