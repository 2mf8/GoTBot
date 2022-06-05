package plugins

import(
	"context"
	. "github.com/2mf8/go-tbot-for-rq/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
)

type Reply struct{
}

func (rep *Reply) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	return MESSAGE_IGNORE
}


func init(){
	Register("回复", &Reply{})
}