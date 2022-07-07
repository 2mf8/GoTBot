package data

import (
	"database/sql"
	"flag"
	"fmt"

	//"time"
	"github.com/fatih/color"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
)

var Db *sql.DB
var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "wr@#kequ5060", "the database password")
var iport *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "127.0.0.1", "the database server")
var user = flag.String("user", "sa", "the database user")
var Pool *redis.Pool
var redis_url = flag.String("redis_addr", "127.0.0.1", "the redis url")
var redis_port *int = flag.Int("redis_port", 6379, "the redis port")
var redis_password = flag.String("redis_password", "", "the redis password")
var redis_db *int = flag.Int("redis_db", 0, "the redis db")
var redis_pool_size *int = flag.Int("redis_pool_size", 1000, "the redis pool size")

func init() {
	var err error

	flag.Parse()

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *iport)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
		fmt.Printf(" redis_url:%s", *redis_url)
		fmt.Printf(" redis_port:%d", *redis_port)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;encrypt=disable", *server, *user, *password, *iport)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	Db, err = sql.Open("mssql", connString)

	if err != nil {
		color.Red("[自检] MSSQL连接失败, 请检查是否启动相关服务")
		panic(err)
	}
	color.Green("[自检] MSSQL连接成功")

	//defer Db.Close()

	Pool = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   0,
		IdleTimeout: 0, //300 * time.Second,
		Dial:        dial,
	}
}

func dial() (conn redis.Conn, err error) {
	addr := fmt.Sprintf("%s:%d", *redis_url, *redis_port)
	conn, err = redis.Dial("tcp", addr, redis.DialPassword(*redis_password), redis.DialDatabase(*redis_db))
	if err != nil {
		color.Red("[自检] Redis连接失败, 请检查是否启动相关服务")
		return
	}
	color.Green("[自检] Redis连接成功")
	//defer conn.Close()
	return
}
