package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "github.com/2mf8/GoTBot/data"
	"github.com/2mf8/GoTBot/middleware"
	"github.com/2mf8/GoTBot/status"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Account string `json:"account"`
	Code    string `json:"code"`
}

// curl -X POST -H "Content-Type:application/json" -d "{\"username\":\"admin\", \"code\": \"123456\"}" http://localhost:8080/login
func TokenGetApi(c *gin.Context) {
	l := Login{}
	li := Login{}
	len := c.Request.ContentLength
	body := make([]byte, len)
	c.Request.Body.Read(body)
	fmt.Println(len, body, string(body))
	e := json.Unmarshal(body, &l)
	fmt.Println(e)
	if e != nil {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code": status.BadRequest,
			"msg":  "提交的数据格式有误，请参看接口文档说明。",
		})
		c.Abort()
		return
	}
	u, e := database.UserInfoGet(l.Account, l.Account, l.Account, l.Account)
	fmt.Println(e)
	if e != nil {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code": status.LoginError,
			"msg":  "查询失败，可能是账号未在系统中注册。请先注册或联系系统管理员。",
		})
		c.Abort()
		return
	}
	bv, err := database.RedisGet(l.Account)
	if err != nil {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code": status.LoginError,
			"msg":  "登录失败，验证码已失效",
		})
		c.Abort()
		return
	}
	json.Unmarshal(bv, &li)
	if l.Account == li.Account && l.Code == li.Code {
		ts, ti, es := middleware.GenTokens(u.UserId.String, u.Username.String, u.UserAvatar.String, u.ServerNumber.String, u.Email.String, u.UserRole, 2)
		if es[0] != nil {
			c.JSON(int(status.ExpectationFailed), gin.H{
				"code": status.GetTokenError,
				"msg":  "获取Token失败，请重试",
			})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":      200,
				"msg":       "登录成功",
				"token":     ts[0],
				"refresh":   ts[1],
				"expires_at":   ti - 600,
				"user_info": u,
			})
			//c.String(http.StatusOK, `curl -H "Authorization: Bearer %s" -H "Refresh: Bearer %s" http://localhost:8080/prices/四`, ts[0], ts[1])
		}
	} else {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code": status.LoginError,
			"msg":  "获取Token失败，用户名错误或验证码错误, 可能验证码已失效, 请前往机器人所在频道重新获取",
		})
		c.Abort()
		return
	}
}
