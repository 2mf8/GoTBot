package database

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func RedisGet(key string) ([]byte, error) {
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", key)
	c.Flush()
	vb, err := redis.Bytes(c.Receive())
	if err != nil {
		return []byte{}, err
	}
	return vb, nil
}

func RedisSet(key string, bw_set []byte) {
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", key, bw_set, "EX", "300")
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
}

func RedisRemove(key string) {
	c := Pool.Get()
	defer c.Close()
	c.Send("Del", key)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
}
