package plugins

import(
	"context"
	. "github.com/2mf8/GoTBot/utils"
)

type CReply struct{
}

func (rep *CReply) ChannelDo(ctx *context.Context, botId, botChannelId int64, guildId, channelId, userId uint64, rawMsg, card string, super, userRole bool, rs, rd, rf int) (retStuct RetChannelStuct) {
	return RetChannelStuct{
		RetVal: MESSAGE_IGNORE,
	}
}


func init(){
	ChannelRegister("回复", &CReply{})
}