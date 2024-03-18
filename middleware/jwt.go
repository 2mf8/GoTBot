package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	database "github.com/2mf8/GoTBot/data"
	"github.com/2mf8/GoTBot/status"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(database.AllConfig.JwtKey)
var jwtRefreshKey = []byte(database.AllConfig.RefreshKey)
var setUsername = ""
var setUserRole = 0
var setServerNumber = ""
var setUserId = ""

type JwtClaims struct {
	UserId       string `json:"user_id"`
	Username     string `json:"user_name"`
	UserRole     int    `json:"user_role"` // 1<<1 黄小姐 1<<2 奇乐 1<<30 系统
	UserAvatar   string `json:"user_avatar"`
	ServerNumber string `json:"server_number"`
	Email        string `josn:"email"`
	jwt.RegisteredClaims
}

// var LoginMap map[string]int = map[string]int{}

func GetKeys() [2][]uint8 {
	var bs [2][]uint8
	bs[0] = jwtKey
	bs[1] = jwtRefreshKey
	return bs
}

func GenJwtClaims(user_id, user_name, user_avatar, server_number, email string, user_role, timeout int) JwtClaims {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(timeout)))
	claims := JwtClaims{
		UserId:       user_id,
		Username:     user_name,
		UserRole:     user_role,
		UserAvatar:   user_avatar,
		ServerNumber: server_number,
		Email:        email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
		},
		//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(timeout))),
	}
	// LoginMap[user_id] = claims.UserRole
	return claims
}

func (j *JwtClaims) GenToken(key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, j)
	return token.SignedString(key)
}

func GenTokens(user_id, user_name, user_avatar, server_number, email string, user_role, timeout int) ([2]string, int64, [2]error) {
	var tokens [2]string
	var errs [2]error
	_timeout := time.Now().Add(time.Hour* time.Duration(timeout)).Unix()
	tg := GenJwtClaims(user_id, user_name, user_avatar, server_number, email, user_role, timeout)
	t, e1 := tg.GenToken(jwtKey)
	tokens[0] = t
	errs[0] = e1
	rtg := GenJwtClaims(user_id, user_name, user_avatar, server_number, email, user_role, timeout*36)
	rt, e2 := rtg.GenToken(jwtRefreshKey)
	tokens[1] = rt
	errs[1] = e2
	return tokens, _timeout, errs
}

func RefreshTokens(refreshTokens [2]string, timeout int) ([2]string, string, int64, [2]error) {
	tokens := [2]string{}
	errs := [2]error{}
	_timeout := time.Now().Add(time.Hour* time.Duration(timeout)).Unix()
	bs := GetKeys()
	if refreshTokens[0] == "" {
		status := "token为空，请正常登录"
		return tokens, status, -1, errs
	}
	_r, err := ParseToken(refreshTokens[0], bs[0])
	fmt.Println(err)
	if err != nil {
		rj, e := ParseToken(refreshTokens[1], bs[1])
		fmt.Println(rj, e)
		if e != nil {
			status := "token已过期，请正常登录"
			return tokens, status, -1, errs
		}
		tg := GenJwtClaims(rj.UserId, rj.Username, rj.UserAvatar, rj.ServerNumber, rj.Email, rj.UserRole, timeout)
		t, e1 := tg.GenToken(bs[0])
		tokens[0] = t
		errs[0] = e1
		setUsername = rj.Username
		setUserRole = rj.UserRole
		setServerNumber = rj.ServerNumber
		setUserId = rj.UserId
		rtg := GenJwtClaims(rj.UserId, rj.Username, rj.UserAvatar, rj.ServerNumber, rj.Email, rj.UserRole, timeout*36)
		rt, e2 := rtg.GenToken(bs[1])
		tokens[1] = rt
		errs[1] = e2
		status := "token刷新成功"
		fmt.Printf(`curl -H "Authorization: Bearer %s" -H "Refresh: Bearer %s" http://localhost:8080/price/四`, tokens[0], tokens[1])
		return tokens, status, _timeout, errs
	}
	setUsername = _r.Username
	setUserRole = _r.UserRole
	setServerNumber = _r.ServerNumber
	setUserId = _r.UserId
	status := "登录状态正常"
	tokens = refreshTokens
	return tokens, status, _timeout, errs
}

func ParseToken(tokenString string, key []byte) (*JwtClaims, error) {
	var j = new(JwtClaims)
	token, err := jwt.ParseWithClaims(tokenString, j, func(t *jwt.Token) (interface{}, error) { return key, nil })
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return j, nil
	}
	return nil, errors.New("invaild token")
}

func JWTAuthMiddlewareOrRefreshToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		tokens := [2]string{}
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Println(authHeader)
		//if c.Request.Method != "GET" {
		if authHeader == "" {
			refreshHeader := c.Request.Header.Get("Refresh")
			fmt.Println(refreshHeader)
			if refreshHeader != "" {
				rparts := strings.SplitN(refreshHeader, " ", 2)
				if len(rparts) == 2 && rparts[0] == "Bearer" {
					tokens[1] = rparts[1]
				}
			}
			_, _status, _, _ := RefreshTokens(tokens, 2)
			if !strings.Contains(_status, "状态正常") {
				c.JSON(int(status.ExpectationFailed), gin.H{
					"code": 200,
					"msg":  _status,
				})
				c.Abort()
				return
			}
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(int(status.Unauthorized), gin.H{
				"code": status.TokenInvalid,
				"msg":  "访问失败，无效的token，请登录。",
			})
			c.Abort()
			return
		}
		c.Set("user_name", setUsername)
		c.Set("user_role", setUserRole)
		c.Set("server_number", setServerNumber)
		c.Set("user_id", setUserId)
		c.Next()
		/*} else {
			c.JSON(int(status.MethodNotAllowed), gin.H{
				"code": status.MethodNotAllowed,
				"msg":  "请使用 POST 方法提交",
			})
			c.Abort()
			return
		}*/
	}
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			c.JSON(int(status.Unauthorized), gin.H{
				"code": status.TokenNull,
				"msg":  "访问失败，token为空，请登录",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(int(status.Unauthorized), gin.H{
				"code": status.TokenInvalid,
				"msg":  "访问失败，无效的token，请登录。",
			})
			c.Abort()
			return
		}
		ms, err := ParseToken(parts[1], jwtKey)
		if err != nil {
			refreshHeader := c.Request.Header.Get("Refresh")
			fmt.Println(refreshHeader)
			if refreshHeader == "" {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.RefreshTokenNull,
					"msg":  "访问失败，refreshToken为空，请登录",
				})
				c.Abort()
				return
			}
			rparts := strings.SplitN(refreshHeader, " ", 2)
			if !(len(rparts) == 2 && rparts[0] == "Bearer") {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.RefreshTokenInvalid,
					"msg":  "访问失败，无效的refreshToken，请登录。",
				})
				c.Abort()
				return
			}
			rms, err := ParseToken(rparts[1], jwtRefreshKey)
			if err != nil {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.RefreshTokenInvalid,
					"msg":  "访问失败，无效的refreshToken，请登录。",
				})
				c.Abort()
				return
			}
			c.Set("Username", rms.Username)
			c.Next()
		} else {
			c.Set("Username", ms.Username)
			c.Next()
		}
	}
}
