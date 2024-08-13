package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "gopkg.in/guregu/null.v3/zero"
)

var IMap = make(map[string]string, 0)

type SubscribeSync struct {
	IsTrue  bool
	SubSync map[string]string
}

func SubscribeRead() (m map[string]string, err error) {
	jsonFile, err := os.Open("m2j.json")
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
	json.Unmarshal(jsonData, &m)
	return
}

func SubscribeCreate(k, v string) {
	m, err := SubscribeRead()
	if err != nil {
		IMap[k] = v
		b, _ := json.MarshalIndent(IMap, "", "\t")
		err := os.WriteFile("m2j.json", b, 0644)
		if err != nil {
			fmt.Println("订阅存储失败")
		}
	} else {
		m[k] = v
		b, _ := json.MarshalIndent(m, "", "\t")
		err := os.WriteFile("m2j.json", b, 0644)
		if err != nil {
			fmt.Println("订阅存储失败")
		}
	}
}

func SubscribeDelete(k string) {
	m, err := SubscribeRead()
	if err != nil {
		fmt.Println("删除失败")
	} else {
		delete(m, k)
		b, e := json.MarshalIndent(m, "", "\t")
		fmt.Println(string(b), e)
		err := os.WriteFile("m2j.json", b, 0644)
		if err != nil {
			fmt.Println("订阅存储失败")
		}
	}
}
