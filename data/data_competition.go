package database

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type CompContent struct {
	Two      string
	Three    string
	Four     string
	Five     string
	Six      string
	Seven    string
	Skewb    string
	Pyraminx string
	Square   string
	Megaminx string
	Clock    string
}

type CompOptions struct {
	Sessions     int
	StartTime    int64
	EndTime      int64
	Items        []string
	CompContents *CompContent
}

var ScrambleMap = map[string]string{
	"2":     "222",
	"3":     "333",
	"4":     "444",
	"5":     "555",
	"6":     "666",
	"7":     "777",
	"sk":    "skewb",
	"py":    "pyram",
	"sq":    "sq1",
	"cl":    "clock",
	"mx":    "minx",
	"fm":    "333fm",
	"222":   "222",
	"333":   "333",
	"444":   "444",
	"555":   "555",
	"666":   "666",
	"777":   "777",
	"skewb": "skewb",
	"pyram": "pyram",
	"sq1":   "sq1",
	"clock": "clock",
	"minx":  "minx",
	"333fm": "333fm",
}

var ScrambleIndexMap = map[string]int{
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
}

// var scrambleItems = []string{"222", "444", "555", "skewb", "pyram", "sq1", "clock", "minx"}

//var scrambleItems = []string{"222", "333", "444"}

func GetScrambles(s []string, n int) (string, CompContent, error) {
	var c CompContent
	for _, v := range s {
		resp, err := http.Get("http://2mf8.cn:2014/scramble/.txt?=" + v + "*" + strconv.Itoa(n))
		if err != nil {
			return "获取失败", CompContent{}, err
		}
		vs, _ := io.ReadAll(resp.Body)
		switch v {
		case "222":
			c.Two = string(vs)
		case "333":
			c.Three = string(vs)
		case "444":
			c.Four = string(vs)
		case "555":
			c.Five = string(vs)
		case "666":
			c.Six = string(vs)
		case "777":
			c.Seven = string(vs)
		case "skewb":
			c.Skewb = string(vs)
		case "pyram":
			c.Pyraminx = string(vs)
		case "sq1":
			c.Square = string(vs)
		case "clock":
			c.Clock = string(vs)
		case "minx":
			c.Megaminx = string(vs)
		}
	}
	return "", c, nil
}

func (c *CompOptions) CompetitionCreate(day int, si []string) (err error) {
	i := c.Sessions
	_, compContents, _ := GetScrambles(si, 5)
	filepath := "private/competition/"
	filename := fmt.Sprintf("%scompetition%d.json", filepath, i)
	co := CompOptions{
		Sessions:     i,
		StartTime:    time.Now().Unix(),
		EndTime:      time.Now().AddDate(0, 0, day).Unix(),
		Items:        si,
		CompContents: &compContents,
	}
	if !PathExists(filepath) {
		if err := os.MkdirAll(filepath, 0777); err != nil {
			fmt.Println("failed to mkdir")
			return err
		}
	}
	output, err := json.MarshalIndent(co, "", "\t")
	fmt.Println(len(output))
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	err = os.WriteFile(filename, output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return
	}
	err = os.WriteFile("private/competition/latest.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return
	}
	return
}

func (c *CompOptions) CompetitionUpdate(sa []string) (tip string, err error) {
	if time.Now().Unix() > c.EndTime || time.Now().Unix() < c.StartTime {
		return "赛季项目更新错误，赛季不存在或已过期", nil
	}
	filename := fmt.Sprintf("private/competition/competition%d.json", c.Sessions)
	var itemTemp []string
	for _, v := range sa {
		fmt.Println(v)
		itemContent := JudgeItem(v, c.Items)
		if itemContent == "" {
			itemTemp = append(itemTemp, v)
		}
	}
	c.Items = append(c.Items, itemTemp...)
	if len(itemTemp) != 0 {
		_, compContents, _ := GetScrambles(itemTemp, 5)
		for _, vi := range itemTemp {
			switch vi {
			case "222":
				c.CompContents.Two = compContents.Two
			case "333":
				c.CompContents.Three = compContents.Three
			case "444":
				c.CompContents.Four = compContents.Four
			case "555":
				c.CompContents.Five = compContents.Five
			case "666":
				c.CompContents.Six = compContents.Six
			case "777":
				c.CompContents.Seven = compContents.Seven
			case "skewb":
				c.CompContents.Skewb = compContents.Skewb
			case "pyram":
				c.CompContents.Pyraminx = compContents.Pyraminx
			case "sq1":
				c.CompContents.Square = compContents.Square
			case "clock":
				c.CompContents.Clock = compContents.Clock
			case "minx":
				c.CompContents.Megaminx = compContents.Megaminx
			}
		}
	} else {
		return "无新增内容", nil
	}
	output, err := json.MarshalIndent(c, "", "\t")
	fmt.Println(len(output))
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	err = os.WriteFile(filename, output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return
	}
	err = os.WriteFile("private/competition/latest.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return
	}
	return
}

func CompetitionRead() (c CompOptions, err error) {
	jsonFile, err := os.Open("private/competition/latest.json")
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
	json.Unmarshal(jsonData, &c)
	return
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func JudgeItem(str string, items []string) string {
	for _, item := range items {
		//if strings.Index(str, k) != -1 {
		if strings.Contains(str, item) {
			return item
		}
	}
	return ""
}

func ToGetScramble(s string) string {
	tgc, ok := ScrambleMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetScrambleIndex(s string) int {
	tgi, ok := ScrambleIndexMap[s]
	if !ok {
		tgi = 0
	}
	return tgi
}
