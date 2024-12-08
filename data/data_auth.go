package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "gopkg.in/guregu/null.v3/zero"
)

var AIMap = make(map[string]string, 0)

type AuthSync struct {
	IsTrue  bool
	SubSync map[string]string
}

func AuthRead() (m map[string]string, err error) {
	jsonFile, err := os.Open("a2j.json")
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

func AuthCreate(k, v string) {
	m, err := AuthRead()
	if err != nil {
		IMap[k] = v
		b, _ := json.MarshalIndent(IMap, "", "\t")
		err := os.WriteFile("a2j.json", b, 0644)
		if err != nil {
			fmt.Println("订阅存储失败")
		}
	} else {
		m[k] = v
		b, _ := json.MarshalIndent(m, "", "\t")
		err := os.WriteFile("a2j.json", b, 0644)
		if err != nil {
			fmt.Println("订阅存储失败")
		}
	}
}

func AuthDelete(k string) {
	m, err := AuthRead()
	if err != nil {
		fmt.Println("删除失败")
	} else {
		delete(m, k)
		b, e := json.MarshalIndent(m, "", "\t")
		fmt.Println(string(b), e)
		err := os.WriteFile("a2j.json", b, 0644)
		if err != nil {
			fmt.Println("订阅存储失败")
		}
	}
}
