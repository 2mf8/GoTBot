package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gomodule/redigo/redis"
)

type ServerAuthSet struct {
	Groups []string
}

type ServerAuthSetSync struct {
	IsTrue            bool
	ServerAuthSetSync *ServerAuthSet
}

func (s *ServerAuthSet) ServerAuthsSet() error {
	output, err := json.MarshalIndent(&s, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return err
	}
	err = os.WriteFile("server_Auths.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return err
	}
	return nil
}

func ServerAuthsRead() (s ServerAuthSet, err error) {
	jsonFile, err := os.Open("server_Auths.json")
	if err != nil {
		fmt.Println("Error reading JSON File:", err)
		return
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	json.Unmarshal(jsonData, &s)
	//fmt.Println(k)
	return
}

func ServerAuthsGet() (s ServerAuthSetSync, err error) {
	sns := ServerAuthSet{
		Groups: []string{},
	}
	snss := ServerAuthSetSync{
		IsTrue:            true,
		ServerAuthSetSync: &sns,
	}
	var vb []byte
	var bw_set []byte

	bw := "server_Auths"
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-守卫", bw)
		snr, err := ServerAuthsRead()
		if err != nil {
			snss = ServerAuthSetSync{
				IsTrue:            false,
				ServerAuthSetSync: &sns,
			}
			snss.ServerAuthSetSync.ServerAuthsSet()
		}
		snss.ServerAuthSetSync = &snr
		bw_set, _ = json.Marshal(&snss)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return snss, err
	}
	err = json.Unmarshal(vb, &snss)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	fmt.Println("[Redis] Key(", bw, ") Value(", snss.IsTrue, *snss.ServerAuthSetSync, ")") //测试用
	return snss, err
}

func (s *ServerAuthSetSync) ServerAuthUpdate(groupId string) error {
	if groupId != "" {
		s.ServerAuthSetSync.Groups = append(s.ServerAuthSetSync.Groups, groupId)
	}
	bw := "server_Auths"
	var bw_set []byte
	ServerAuthSetSync := ServerAuthSetSync{
		IsTrue:            true,
		ServerAuthSetSync: s.ServerAuthSetSync,
	}
	bw_set, _ = json.Marshal(&ServerAuthSetSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
	err = s.ServerAuthSetSync.ServerAuthsSet()
	return err
}

func (s *ServerAuthSetSync) ServerAuthsDelete(groupId string) {
	groupsTemp := []string{}
	if groupId != "" {
		for _, v := range s.ServerAuthSetSync.Groups {
			if v != groupId {
				groupsTemp = append(groupsTemp, v)
			}
			continue
		}
		s.ServerAuthSetSync.Groups = groupsTemp
		bw := "server_Auths"
		var bw_set []byte
		ServerAuthSetSync := ServerAuthSetSync{
			IsTrue:            true,
			ServerAuthSetSync: s.ServerAuthSetSync,
		}
		bw_set, _ = json.Marshal(&ServerAuthSetSync)
		c := Pool.Get()
		defer c.Close()
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, err := c.Receive()
		if err != nil {
			fmt.Println("[错误] Receive出错")
		}
		_ = fmt.Sprintf("%#v", v)
		s.ServerAuthSetSync.ServerAuthsSet()
	}
}

func IsExist(groupId string) bool {
	sgg, _ := ServerAuthsGet()
	for _, v := range sgg.ServerAuthSetSync.Groups {
		if v == groupId {
			return true
		}
	}
	return false
}
