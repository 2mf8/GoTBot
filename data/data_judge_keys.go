package data

import(
	"strings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/gomodule/redigo/redis"
)

type JudgeKeys struct {
	Keys []string
}

type JudgekeysSync struct {
	IsTrue        bool
	JudgekeysSync *JudgeKeys
}

func Judge(str string, key JudgeKeys) string {
	for _, k := range key.Keys {
		//if strings.Index(str, k) != -1 {
		if strings.Contains(str, k) {
			return k
		}
	}
	return ""
}

func JudgeIndex(str string, key JudgeKeys) int {
	for i, k := range key.Keys {
		//if strings.Index(str, k) != -1 {
		if strings.Contains(str, k) {
			return i
		}
	}
	return -1
}

func GetJudgeKeys() (key JudgekeysSync, err error) {
	judgekeys := JudgeKeys{}
	key = JudgekeysSync{
		IsTrue: true,
		JudgekeysSync: &judgekeys,
	}
	var vb []byte
	var bw_set []byte

	bw := "judgekeys"
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-守卫", bw)
		jk, err := JudgeKeysRead()
		key.JudgekeysSync = &jk
		if err != nil {
			key = JudgekeysSync{
				IsTrue: false,
				JudgekeysSync: &judgekeys,
			}
			key.JudgekeysSync.JudgeKeysCreate()
		}
		bw_set, _ = json.Marshal(&key)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return key, err
	}
	err = json.Unmarshal(vb, &key)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	//fmt.Println("[Redis] Key(", bw, ") Value(", key.IsTrue, *key.JudgekeysSync, ")")  //测试用
	return
}

func (k *JudgeKeys) JudgeKeysCreate() error {
	output, err := json.MarshalIndent(&k, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return err
	}
	err = ioutil.WriteFile("judgekeys.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return err
	}
	return nil
}

func JudgeKeysRead() (k JudgeKeys, err error) {
	jsonFile, err := os.Open("judgekeys.json")
	if err != nil {
		fmt.Println("Error reading JSON File:", err)
		return
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	json.Unmarshal(jsonData, &k)
	//fmt.Println(k)
	return
}

func (k *JudgekeysSync) JudgeKeysUpdate(uk ...string) error {
	for _, v := range uk {
		if Judge(v, *k.JudgekeysSync) == "" && v != "" {
			k.JudgekeysSync.Keys = append(k.JudgekeysSync.Keys, v)
		}
	}
	bw := "judgekeys"
	var bw_set []byte
	judgekeysSync := JudgekeysSync{
		IsTrue: true,
		JudgekeysSync: k.JudgekeysSync,
	}
	bw_set, _ = json.Marshal(&judgekeysSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)
	err = k.JudgekeysSync.JudgeKeysCreate()
	return err
}

func (k *JudgekeysSync) JudgeKeysDelete(dk ...string) {
	for _, v := range dk {
		if v == "" {
			continue
		}
		i := JudgeIndex(v, *k.JudgekeysSync)
		if i != -1 {
			if k.JudgekeysSync.Keys[i+1:] != nil {
				k.JudgekeysSync.Keys = append(k.JudgekeysSync.Keys[:i], k.JudgekeysSync.Keys[i+1:]...)
				i--
			}
		}
		bw := "judgekeys"
		var bw_set []byte
		judgekeysSync := JudgekeysSync{
			IsTrue: true,
			JudgekeysSync: k.JudgekeysSync,
		}
		bw_set, _ = json.Marshal(&judgekeysSync)
		c := Pool.Get()
		defer c.Close()
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, err := c.Receive()
		if err != nil {
			fmt.Println("[错误] Receive出错")
		}
		fmt.Sprintf("%#v", v)
		k.JudgekeysSync.JudgeKeysCreate()
	}
}
