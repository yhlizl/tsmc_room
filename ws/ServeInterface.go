package ws

import "github.com/gogf/gf/net/ghttp"

type ServeInterface interface {
	RunWs(r *ghttp.Request)
	GetOnlineUserCount() int
	GetOnlineRoomUserCount(roomId int) int
}
