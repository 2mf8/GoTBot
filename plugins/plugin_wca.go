package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	database "github.com/2mf8/GoTBot/data"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type WCA struct {
}

func (wca *WCA) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	
	s, b := public.Prefix(rawMsg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	if strings.HasPrefix(s, "wca") {
		w_m := strings.TrimSpace(strings.TrimSpace(string([]byte(s)[len("wca"):])))
		fmt.Println(w_m)
		url := "https://www.2mf8.cn:8300/wca/people/" + url.PathEscape(w_m)
		fmt.Println(url)
		resp, _ := http.Get(url)
		s := database.GetInfoResponse{}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal([]byte(body), &s)
		//fmt.Println(fmt.Sprintf("%+v", s))
		//fmt.Println(len(s.Data))
		if s.Data.TotalElements == 1 {
			fmt.Println(s.Data.Data[0].Id)
			s_r := database.WcaPersonHandler(s.Data.Data[0].Id)
			// fmt.Println(s_r)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: s_r,
				},
				MsgId: messageId.Common,
				ReqType: utils.GroupMsg,
			}
		} else if s.Data.TotalElements > 99 {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "搜索范围太大！",
				},
				MsgId: messageId.Common,
				ReqType: utils.GroupMsg,
			}
		} else {
			rankList := ""
			s_r := ""
			count := 0
			for _, l := range s.Data.Data {
				if count < 4 {
					rankList += "\n" + l.Id + " | " + l.Name
				} else {
					rankList += "\n" + l.Id + " | " + l.Name
					if count == 19 {
						rankList += "\n..."
						break
					}
				}
				count++
			}
			if s.Data.TotalElements == 0 {
				s_r = "暂无相关记录，请换个名字或输入对应的WCAID进行搜索。"
			} else {
				s_r = strconv.Itoa(s.Data.TotalElements) + "条记录" + rankList + "\n请换个名字或输入对应的WCAID进行搜索。"
			}
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: s_r,
				},
				MsgId: messageId.Common,
				ReqType: utils.GroupMsg,
			}
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}
