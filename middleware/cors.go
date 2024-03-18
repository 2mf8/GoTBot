package middleware

import (
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept-Encoding, X-CSRF-Token, Cache-Control, Token, X-Access-Token, Content-Length, Accept, Authorization, Refresh")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func Bind(c *gin.Context, req any) error {
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if r, ok := req.(proto.Message); ok {
		if err := proto.Unmarshal(buf, r); err != nil {
			return err
		}
	} else {
		return errors.New("obj is not ProtoMessage")
	}
	return nil
}