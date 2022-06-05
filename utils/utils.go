package utils

import (
	"context"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
)

type Plugin interface {
	Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint)
}

var PluginSet map[string]Plugin

const (
	MESSAGE_BLOCK  uint = 0
	MESSAGE_IGNORE uint = 1
)

func init() {
	PluginSet = make(map[string]Plugin)
}

func Register(k string, v Plugin) {
	PluginSet[k] = v
}
