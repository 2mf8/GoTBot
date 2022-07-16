package config

type Config struct {
	Plugins          []string
	ChannelPlugins   []string
	Admins           []int64
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     int
	DatabaseServer   string
	ServerPort       int
	ScrambleServer   string
	RedisServer      string
	RedisPort        int
	RedisPassword    string
	RedisTable       int
	RedisPoolSize    int
}

var Conf *Config

func init() {
	Conf = &Config{}
}
