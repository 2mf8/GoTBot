package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	database "github.com/2mf8/GoTBot/data"
	"github.com/2mf8/GoTBot/public"
	"github.com/2mf8/GoTBot/utils"
)

type RankPlugin struct{}

func (rp *RankPlugin) Do(ctx *context.Context, botId *utils.BotIdType, groupId *utils.GroupIdType, userId *utils.UserIdType, groupName string, messageId *utils.MsgIdType, rawMsg, card string, botRole, userRole, super bool) (retStuct utils.RetStuct) {
	s, b := public.Prefix(rawMsg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	fmt.Println("rank？")
	if strings.HasPrefix(s, "rank") {
		w_m := strings.TrimSpace(strings.TrimSpace(string([]byte(s)[len("rank"):])))
		ss := strings.Split(w_m, "-")
		if len(ss) > 0 {
			url := "https://www.2mf8.cn:8300/wca/people/" + url.PathEscape(ss[0])
			resp, _ := http.Get(url)
			s := database.GetInfoResponse{}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal([]byte(body), &s)
			//fmt.Println(fmt.Sprintf("%+v", s))
			//fmt.Println(len(s.Data))
			s_r := "获取Rank失败"
			if s.Data.TotalElements == 1 {
				fmt.Println(s.Data.Data[0].Id)
				if len(ss) > 1 {
					s_r = database.RankPersonHandler(s.Data.Data[0].Id, ss[1])
				} else {
					s_r = database.RankPersonHandler(s.Data.Data[0].Id, "sin")
				}
				// fmt.Println(s_r)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: s_r,
					},
					MsgId: messageId.Common,
					OfficalMsgId: messageId.Offical,
					ReqType: utils.GroupMsg,
				}
			} else if s.Data.TotalElements > 99 {
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: "搜索范围太大！",
					},
					MsgId: messageId.Common,
					OfficalMsgId: messageId.Offical,
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
					OfficalMsgId: messageId.Offical,
					ReqType: utils.GroupMsg,
				}
			}
		}
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	} else {
		e, r, t, g := database.ToRank(s)
		if e != "" || r != "" || t != "" || g != "" {
			ns := strings.TrimSpace(strings.TrimPrefix(s, "rank"))
			rs := database.GetTop10Rank(database.ToGetRankInfo(ns))
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, ns)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: rs,
				},
				ReqType: utils.GroupMsg,
				OfficalMsgId: messageId.Offical,
			}
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("Rank", &RankPlugin{})
}
