package public

import (
	"io"
	"strconv"
	"strings"

	"github.com/2mf8/GoPbBot"
	. "github.com/2mf8/GoTBot/config"
	"github.com/BurntSushi/toml"
)

type DataBase struct {
	User     string
	Password string
	Url      string
	Port     int
}

type Redis struct {
	Url      string
	Port     int
	Password string
	Table    int
	PoolSize int
}

type PluginConfig struct{
	Conf []string
	ChannelConf []string
}

func StartsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func EndsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func Contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if StartsWith(s[i:], substr) {
			return true
		}
	}
	return false
}

func IsAdmin(bot *pbbot.Bot, groupId, userId int64) bool {
	memberInfo, _ := bot.GetGroupMemberInfo(groupId, userId, true)
	if strings.ToLower(memberInfo.Role) == "admin" || strings.ToLower(memberInfo.Role) == "owner" {
		return true
	}
	return false
}

func IsGuildAdmin(role []string) bool {
	if role == nil {
		return false
	} else {
		return true
	}
}

func TbotConf() (c PluginConfig, err error) {
	_, err = toml.DecodeFile("conf.toml", Conf)
	pc := PluginConfig{
		Conf: Conf.Plugins,
		ChannelConf: Conf.ChannelPlugins,
	}
	return pc, err
}

func IsBotAdmin(userId int64) bool {
	_, _ = toml.DecodeFile("conf.toml", Conf)
	for _, uId := range Conf.Admins {
		if userId == uId {
			return true
		}
	}
	return false
}

func DataBaseSet() (dbset DataBase, err error) {
	var user string
	var password string
	var url string
	var port int = 0
	_, err = toml.DecodeFile("conf.toml", Conf)
	if Conf.DatabaseUser == "" {
		user = "sa"
	} else {
		user = Conf.DatabaseUser
	}
	if Conf.DatabasePassword == "" {
		password = "@#$mima45"
	} else {
		password = Conf.DatabasePassword
	}
	if Conf.DatabaseServer == "" {
		url = "127.0.0.1"
	} else {
		url = Conf.ScrambleServer
	}
	port = Conf.DatabasePort
	dbset = DataBase{
		User:     user,
		Password: password,
		Url:      url,
		Port:     port,
	}
	return
}

func RedisSet() (dbset Redis, err error) {
	var url string
	var port int = 0
	var password string
	var table int
	var poolSize int
	_, err = toml.DecodeFile("conf.toml", Conf)
	if Conf.RedisServer == "" {
		url = "127.0.0.1"
	} else {
		url = Conf.RedisServer
	}
	password = Conf.RedisPassword
	port = Conf.RedisPort
	table = Conf.RedisTable
	poolSize = Conf.RedisPoolSize
	dbset = Redis{
		Url:      url,
		Port:     port,
		Password: password,
		Table:    table,
		PoolSize: poolSize,
	}
	return
}

/*func Port() int {
	_, err = toml.DecodeFile("conf.toml", Conf)
}*/

func IsConnErr(err error) bool {
	var needNewConn bool
	if err == nil {
		return false
	}
	if err == io.EOF {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connect: connection refused") {
		needNewConn = true
	}
	return needNewConn
}

func Prefix(s string, p string) (r string, b bool) {
	if StartsWith(s, p) {
		r = strings.TrimSpace(string([]byte(s)[len(p):]))
		return r, true
	}
	r = s
	return r, false
}

func ArrayStringToArrayInt64(s []string) (g []int64) {
	for _, str := range s {
		i, e := strconv.Atoi(str)
		if e != nil {
			continue
		}
		g = append(g, int64(i))
	}
	return g
}