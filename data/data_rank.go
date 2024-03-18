package database

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Top10RankInfo struct {
	Event    string      `json:"event,omitempty"`
	Region   string      `json:"region,omitempty"`
	Type     string      `json:"type,omitempty"`
	Gender   string      `json:"gender,omitempty"`
	ItemList []*ItemList `json:"itemList,omitempty"`
}

type ItemList struct {
	PlayerName string `json:"playerName,omitempty"`
	BestResult string `json:"bestResult,omitempty"`
}

type GetTop10RankInfo struct {
	Code    int32          `json:"code"`
	Data    *Top10RankInfo `json:"data,omitempty"`
	Message string         `json:"message"`
}

type GetRankResultInfo struct {
	Code    int32       `json:"code"`
	Data    *RankResult `json:"data,omitempty"`
	Message string      `json:"message"`
}

var eventMap = map[string]string{
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

var regionMap = map[string]string{
	"wr":  "wr",  // world record
	"nr":  "nr",  // record
	"asr": "asr", // asia record
	"afr": "afr", // africa record
	"er":  "er",  // europe record
	"nar": "nar", // north america record
	"sar": "sar", // south america record
	"ocr": "ocr", // oceania record
}

var typeMap = map[string]string{
	"sin":     "sin", // single result
	"avg":     "avg", // average result
	"single":  "sin",
	"average": "avg",
	"最佳":      "sin",
	"平均":      "avg",
}

var genderMap = map[string]string{
	"m":      "m",
	"f":      "f",
	"all":    "all",
	"male":   "m",
	"female": "f",
	"男":      "m",
	"女":      "f",
}

func ToGetEvent(s string) string {
	tgc, ok := eventMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetRegion(s string) string {
	tgc, ok := regionMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetType(s string) string {
	tgc, ok := typeMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetGender(s string) string {
	tgc, ok := genderMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetRankInfo(s string) (event, region, rtype, gender string) {
	e := ""
	r := ""
	t := ""
	g := ""
	strs := strings.Split(s, " ")
	for _, v := range strs {
		if e == "" {
			e = ToGetEvent(strings.TrimSpace(v))
		}
		if r == "" {
			r = ToGetRegion(strings.TrimSpace(v))
		}
		if t == "" {
			t = ToGetType(strings.TrimSpace(v))
		}
		if g == "" {
			g = ToGetGender(strings.TrimSpace(v))
		}
	}
	if e == "" {
		e = "333"
	}
	if r == "" {
		r = "nr"
	}
	if t == "" {
		t = "sin"
	}
	if g == "" {
		g = "all"
	}
	fmt.Println("ToGetRankInfo", e, r, t, g)
	return e, r, t, g
}

func GetTop10Rank(event, region, rtype, gender string) string {
	url := fmt.Sprintf("https://2mf8.cn:8300/top10rank/%s/%s/%s/%s", event, region, rtype, gender)
	resp, _ := http.Get(url)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	gr := GetTop10RankInfo{}
	json.Unmarshal([]byte(body), &gr)
	r := ""
	t := ""
	g := ""
	if gr.Data.Region == "wr" {
		r = "世界纪录"
	}
	if gr.Data.Region == "nr" {
		r = "中国记录"
	}
	if gr.Data.Region == "asr" {
		r = "亚洲记录"
	}
	if gr.Data.Region == "afr" {
		r = "非洲记录"
	}
	if gr.Data.Region == "er" {
		r = "欧洲记录"
	}
	if gr.Data.Region == "nar" {
		r = "北美洲记录"
	}
	if gr.Data.Region == "sar" {
		r = "南美洲记录"
	}
	if gr.Data.Region == "ocr" {
		r = "大洋洲记录"
	}
	if gr.Data.Type == "sin" {
		t = "最佳成绩"
	}
	if gr.Data.Type == "avg" {
		t = "平均成绩"
	}
	if gr.Data.Gender == "m" {
		g = "男"
	}
	if gr.Data.Gender == "f" {
		g = "女"
	}
	if gr.Data.Gender == "all" {
		g = "不限"
	}
	result := fmt.Sprintf("事件：%s\n区域：%s\n类型：%s\n性别：%s\n姓名 | 成绩", gr.Data.Event, r, t, g)
	for _, i := range gr.Data.ItemList {
		result += fmt.Sprintf("\n%s | %s", i.PlayerName, i.BestResult)
	}
	return result
}

func ToRank(s string) (event, region, rtype, gender string) {
	e := ""
	r := ""
	t := ""
	g := ""
	strs := strings.Split(s, " ")
	for _, v := range strs {
		if e == "" {
			e = ToGetEvent(strings.TrimSpace(v))
		}
		if r == "" {
			r = ToGetRegion(strings.TrimSpace(v))
		}
		if t == "" {
			t = ToGetType(strings.TrimSpace(v))
		}
		if g == "" {
			g = ToGetGender(strings.TrimSpace(v))
		}
	}
	return e, r, t, g
}

func RankPersonHandler(s, rtype string) string {
	singleUrl := "https://www.2mf8.cn:8300/rank/result/" + url.PathEscape(s)
	resp, _ := http.Get(singleUrl)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	gr := GetRankResultInfo{}
	json.Unmarshal([]byte(body), &gr)
	t := ToGetType(rtype)
	if t == "avg" {
		result := fmt.Sprintf("%s\n%s,%s,%s\n事件 平均成绩 | 国家 | 大陆 | 世界", gr.Data.PeopleInfo.Name, gr.Data.PeopleInfo.Id, gr.Data.PeopleInfo.CountryId, gr.Data.PeopleInfo.Gender)
		for _, i := range gr.Data.RankItemsResult {
			result += fmt.Sprintf("\n%v %v | %v | %v | %v", i.EventId, i.Best.Average, i.CountryRank.Average, i.ContinentRank.Average, i.WorldRank.Average)
		}
		return result
	} else {
		result := fmt.Sprintf("%s\n%s,%s,%s\n事件 最佳成绩 | 国家 | 大陆 | 世界", gr.Data.PeopleInfo.Name, gr.Data.PeopleInfo.Id, gr.Data.PeopleInfo.CountryId, gr.Data.PeopleInfo.Gender)
		for _, i := range gr.Data.RankItemsResult {
			result += fmt.Sprintf("\n%v %v | %v | %v | %v", i.EventId, i.Best.Best, i.CountryRank.Best, i.ContinentRank.Best, i.WorldRank.Best)
		}
		return result
	}
}
