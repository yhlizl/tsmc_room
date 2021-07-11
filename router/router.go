package router

import (
	"go-gofram-chat/app/service/session"
	"go-gofram-chat/controller"
	"go-gofram-chat/ws/primary"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	s.AddStaticPath("", "public/static")
	
	session.EnableCookieSession(s)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/", controller.Index)
		group.POST("/login", controller.Login)
		group.GET("/logout", controller.Logout)
		group.GET("/ws", primary.Start)
		group.Middleware(session.AuthSessionMiddle())
		group.GET("/home", controller.Home)
		group.GET("/room/:room_id", controller.Room)
		group.GET("/private-chat", controller.PrivateChat)
		group.POST("/img-kr-upload", controller.ImgKrUpload)
		group.GET("/pagination", controller.Pagination)
	})
}
