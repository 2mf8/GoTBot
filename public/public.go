package public

import (
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MessageType int

const (
	GuildMessage MessageType = iota
	GroupMessage
	PrivateMessage
	DirectMessage
	Undefined
)

type PluginConfig struct {
	Conf []string
}

type Node struct {
	XMLName xml.Name
	Attr    []xml.Attr `xml:",any,attr"`
}

type BotLogin struct {
	AppId       uint64
	AccessToken string
}

type GetAppAccessToken struct {
	AppId        string `json:"appId"`
	ClientSecret string `json:"clientSecret"`
}

type QQGetAppAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

func StartsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func EndsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func Contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if StartsWith(s[i:], substr) {
			return true
		}
	}
	return false
}

func IsAdmin(role string) bool {
	if strings.ToLower(role) == "admin" || strings.ToLower(role) == "owner" {
		return true
	}
	return false
}

func IsCompAdmin(roles []string) bool {
	for _, role := range roles {
		if role == "13818102" {
			return true
		}
	}
	return false
}

func IsBotAdmin(userId string, admins []string) bool {
	for _, uId := range admins {
		if userId == uId {
			return true
		}
	}
	return false
}

func IsConnErr(err error) bool {
	var needNewConn bool
	if err == nil {
		return false
	}
	if err == io.EOF {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connect: connection refused") {
		needNewConn = true
	}
	return needNewConn
}

func Prefix(s string, p string) (r string, b bool) {
	if StartsWith(s, p) {
		r = strings.TrimSpace(string([]byte(s)[len(p):]))
		return r, true
	}
	r = s
	return r, false
}

func ArrayStringToArrayInt64(s []string) (g []int64) {
	for _, str := range s {
		i, e := strconv.Atoi(str)
		if e != nil {
			continue
		}
		g = append(g, int64(i))
	}
	return g
}

func GuildAtConvert(str string) (string, []string) {

	var re = regexp.MustCompile(`<[\s\S\\s\\S]+?/>`)

	var node Node

	var users []string

	textList := re.Split(str, -1)
	text := strings.Join(textList, " ")

	codeList := re.FindAllString(str, -1)

	for _, c := range codeList {
		err := xml.Unmarshal([]byte(c), &node)
		if err != nil {
			continue
		}
	}

	attrMap := make(map[string]string)
	for _, attr := range node.Attr {
		attrMap[attr.Name.Local] = html.UnescapeString(attr.Value)
		users = append(users, attr.Value)
	}

	return text, users
}

func ConvertTime(str string) int32 {
	var duration int = -1
	reg4 := regexp.MustCompile("天")
	reg5 := regexp.MustCompile("小时")
	reg6 := regexp.MustCompile("时")
	reg7 := regexp.MustCompile("分")
	reg8 := regexp.MustCompile("秒")
	str4 := strings.TrimSpace(reg4.ReplaceAllString(str, "d"))
	str4 = strings.TrimSpace(reg5.ReplaceAllString(str4, "h"))
	str4 = strings.TrimSpace(reg6.ReplaceAllString(str4, "h"))
	str4 = strings.TrimSpace(reg7.ReplaceAllString(str4, "m"))
	str4 = strings.TrimSpace(reg8.ReplaceAllString(str4, "s"))
	str4 = str4 + "s"
	reg9 := regexp.MustCompile(`([0-9]+)(d|h|m|s)`)
	m := reg9.FindAllString(str4, -1)
	for _, v := range m {
		if EndsWith(v, "d") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("d")]))
			duration += num * 60 * 60 * 24
		}
		if EndsWith(v, "h") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("h")]))
			duration += num * 60 * 60
		}
		if EndsWith(v, "m") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("m")]))
			duration += num * 60
		}
		if EndsWith(v, "s") {
			num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("s")]))
			duration += num
		}
	}
	return int32(duration)
}

func ConvertJinTime(i int) string {
	var timeString string
	day := i / 86400
	hour := i % 86400 / 3600
	min := i % 3600 / 60
	sec := i % 60
	if i >= 86400 {
		timeString = fmt.Sprintf("%v 天 %v 小时 %v 分钟", day, hour, min)
		return timeString
	}
	if i < 60 {
		timeString = fmt.Sprintf("%v 秒钟", sec)
		return timeString
	}
	if i <= 3600 {
		timeString = fmt.Sprintf("%v 分钟 %v 秒钟", min, sec)
		return timeString
	}
	timeString = fmt.Sprintf("%v 小时 %v 分钟 %v 秒钟", hour, min, sec)
	return timeString
}

func ConvertGradeToInt(s string) (grade []int) {
	reg1 := regexp.MustCompile(`([0-9]*)(:*)([0-9]+)(\.*)([0-9]*)`)
	ss := reg1.FindAllString(s, -1)
	for _, i := range ss {
		if strings.Contains(i, ".") && !strings.Contains(i, ":") {
			f, _ := strconv.ParseFloat(i, 32)
			sf := fmt.Sprintf("%0.3f", f)
			n, _ := strconv.Atoi(strings.Replace(sf, ".", "", -1))
			grade = append(grade, n)
		} else if strings.Contains(i, ".") && strings.Contains(i, ":") {
			as := strings.Split(i, ":")
			f, _ := strconv.ParseFloat(as[1], 32)
			sf := fmt.Sprintf("%0.3f", f)
			n, _ := strconv.Atoi(strings.Replace(sf, ".", "", -1))
			n1, _ := strconv.Atoi(as[0])
			nt := n1*1000*60 + n
			grade = append(grade, nt)
		} else {
			n1, _ := strconv.Atoi(i)
			if n1 == 222 || n1 == 333 || n1 == 444 || n1 == 555 || n1 == 666 || n1 == 777 || n1 == 1 {
				continue
			}
			grade = append(grade, n1*1000)
		}
	}
	return grade
}

func RandomString(n int) string {
	var letters = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]byte, n)
	rand.NewSource(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
