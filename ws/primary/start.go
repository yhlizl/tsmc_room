package primary

import (
	"go-gofram-chat/ws"
	"go-gofram-chat/ws/go_ws"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 定義 serve 的映射關係
var serveMap = map[string]ws.ServeInterface{
	"Serve":   &ws.Serve{},
	"GoServe": &go_ws.GoServe{},
}

func Create() ws.ServeInterface {
	// GoServe or Serve
	_type := g.Cfg().GetString("conf.serve_type")
	//log.Println("type of :", _type)
	//_type := viper.GetString("app.serve_type")
	return serveMap[_type]
}

func Start(r *ghttp.Request) {
	Create().RunWs(r)
}

func OnlineUserCount() int {
	return Create().GetOnlineUserCount()
}

func OnlineRoomUserCount(roomId int) int {
	return Create().GetOnlineRoomUserCount(roomId)
}
