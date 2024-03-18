package database

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Info struct {
	Msg                  string        `json:"msg"`
	PageLast             bool          `json:"pageLast"`
	PageEmpty            bool          `json:"pageEmpty"`
	Data                 []*PeopleInfo `json:"data"`
	TotalPages           int           `json:"totalPages"`
	PageFirst            bool          `json:"pageFirst"`
	PageSize             int           `json:"pageSize"`
	PageNumberOfElements int           `json:"pageNumberOfElements"`
	PageNum              int           `json:"pageNum"`
	Retcode              int           `json:"retcode"`
	TotalElements        int           `json:"totalElements"`
}

type PeopleInfo struct {
	CountryId string `json:"countryId"`
	Gender    string `json:"gender"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	SubId     int    `json:"subId"`
}

type WcaResult struct {
	PeopleInfo     PeopleInfo       `json:"info"`
	WcaItemsResult []*WcaItemResult `json:"data"`
}

type WcaItemResult struct {
	Event   any    `json:"event"`
	Best    string `json:"best"`
	Average string `json:"average"`
}

type RankInfo struct {
	Msg                  string `json:"msg"`
	PageLast             bool   `json:"pageLast"`
	PageEmpty            bool   `json:"pageEmpty"`
	Data                 []PeopleRankInfo
	TotalPages           int  `json:"totalPages"`
	PageFirst            bool `json:"pageFirst"`
	PageSize             int  `json:"pageSize"`
	PageNumberOfElements int  `json:"pageNumberOfElements"`
	PageNum              int  `json:"pageNum"`
	Retcode              int  `json:"retcode"`
	TotalElements        int  `json:"totalElements"`
}

type PeopleRankInfo struct {
	Best          int    `json:"best"`
	ContinentRank int    `json:"continentRank"`
	CountryRank   int    `json:"countryRank"`
	EventId       string `json:"eventId"`
	PersonId      string `json:"personId"`
	WorldRank     int    `json:"worldRank"`
}

type RankResult struct {
	PeopleInfo      PeopleInfo      `json:"info"`
	RankItemsResult []*RankItemInfo `json:"data"`
}

type RankItemInfo struct {
	EventId       string         `json:"eventId"`
	Best          *RankBest      `json:"best"`
	ContinentRank *RankContinent `json:"continentRank"`
	CountryRank   *RankCountry   `json:"countryRank"`
	WorldRank     *RankWorld     `json:"worldRank"`
}

type RankBest struct {
	Best    string `json:"best"`
	Average string `json:"average"`
}

type RankContinent struct {
	Best    int `json:"best"`
	Average int `json:"average"`
}

type RankCountry struct {
	Best    int `json:"best"`
	Average int `json:"average"`
}

type RankWorld struct {
	Best    int `json:"best"`
	Average int `json:"average"`
}

type GetWcaResultResponse struct {
	Code    int32      `json:"code"`
	Data    *WcaResult `json:"data"`
	Message string     `json:"message"`
}

type GetInfoResponse struct {
	Code    int32  `json:"code"`
	Data    *Info  `json:"data"`
	Message string `json:"message"`
}

func WcaPersonHandler(s string) string {
	singleUrl := "https://www.2mf8.cn:8300/wca/result/" + url.PathEscape(s)
	resp, _ := http.Get(singleUrl)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	gr := GetWcaResultResponse{}
	json.Unmarshal([]byte(body), &gr)
	result := fmt.Sprintf("%s\n%s,%s,%s\n事件 最佳 | 平均", gr.Data.PeopleInfo.Name, gr.Data.PeopleInfo.Id, gr.Data.PeopleInfo.CountryId, gr.Data.PeopleInfo.Gender)
	for _, i := range gr.Data.WcaItemsResult {
		result += fmt.Sprintf("\n%v %s | %s", i.Event, i.Best, i.Average)
	}
	return result
}

func SearchPeople(nameOrId string) GetInfoResponse {
	url := fmt.Sprintf("https://www.2mf8.cn:8300/wca/people/%s", url.PathEscape(nameOrId))
	resp, _ := http.Get(url)
	s := GetInfoResponse{}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &s)
	return s
}
