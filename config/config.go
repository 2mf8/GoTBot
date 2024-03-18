package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Plugins          []string
	AppId            uint64
	AccessToken      string
	ClientSecret     string
	Admins           []string
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     int
	DatabaseServer   string
	DatabaseName     string
	ServerPort       int
	ScrambleServer   string
	RedisServer      string
	RedisPort        int
	RedisPassword    string
	RedisTable       int
	RedisPoolSize    int
	JwtKey           string
	RefreshKey       string
}

var Conf *Config = &Config{}

func AllConfig() Config {
	_, err := toml.DecodeFile("conf.toml", Conf)
	if err != nil {
		return *Conf
	}
	return *Conf
}
