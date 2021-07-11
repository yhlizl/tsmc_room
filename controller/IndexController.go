package controller

import (
	"go-gofram-chat/app/service/helper"
	"go-gofram-chat/app/service/message_service"
	"go-gofram-chat/app/service/user_service"
	"go-gofram-chat/ws/primary"
	"strconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func Index(r *ghttp.Request) {
	// 已登錄跳轉room界面，多頁面應該考慮放在中間件實現
	userInfo := user_service.GetUserInfo(r)
	if len(userInfo) > 0 {
		r.Response.RedirectTo("/home")
		return
	}

	OnlineUserCount := primary.OnlineUserCount()
	r.Response.WriteTpl("login.html", g.Map{
		"OnlineUserCount": OnlineUserCount,
	})

}

func Login(r *ghttp.Request) {
	user_service.Login(r)
}

func Logout(r *ghttp.Request) {
	user_service.Logout(r)
}

func Home(r *ghttp.Request) {
	userInfo := user_service.GetUserInfo(r)
	rooms := []map[string]interface{}{
		{"id": 1, "num": primary.OnlineRoomUserCount(1), "roomname": "台積大家庭"},
		{"id": 2, "num": primary.OnlineRoomUserCount(2), "roomname": "台積 竹科"},
		{"id": 3, "num": primary.OnlineRoomUserCount(3), "roomname": "台積 中科"},
		{"id": 4, "num": primary.OnlineRoomUserCount(4), "roomname": "台積 南科"},
		{"id": 5, "num": primary.OnlineRoomUserCount(5), "roomname": "台積 龍潭"},
		{"id": 6, "num": primary.OnlineRoomUserCount(6), "roomname": "台積 Arizona"},
	}

	r.Response.WriteTpl("index.html", g.Map{
		"rooms":     rooms,
		"user_info": userInfo,
	})

	// c.HTML(http.StatusOK, "index.html", gin.H{
	// 	"rooms":     rooms,
	// 	"user_info": userInfo,
	// })
}

func Room(r *ghttp.Request) {
	roomId := r.GetString("room_id")

	rooms := []string{"1", "2", "3", "4", "5", "6"}

	if !helper.InArray(roomId, rooms) {
		r.Response.RedirectTo("/room/1")
		return
	}

	userInfo := user_service.GetUserInfo(r)
	msgList := message_service.GetLimitMsg(roomId, 0)

	roomname := ""
	switch roomId {
	case "1":
		roomname = "台積大家庭"
	case "2":
		roomname = "台積 竹科"
	case "3":
		roomname = "台積 中科"
	case "4":
		roomname = "台積 南科"
	case "5":
		roomname = "台積 龍潭"
	case "6":
		roomname = "台積 Arizona"
	}
	//fmt.Println(roomId)
	r.Response.WriteTpl("room.html", g.Map{
		"user_info":      userInfo,
		"msg_list":       msgList,
		"msg_list_count": len(msgList),
		"room_id":        roomId,
		"roomename":      roomname,
	})

}

func PrivateChat(r *ghttp.Request) {

	roomId := r.GetString("room_id")
	toUid := r.GetString("uid")

	userInfo := user_service.GetUserInfo(r)

	uid := strconv.Itoa(int(userInfo["uid"].(uint)))

	msgList := message_service.GetLimitPrivateMsg(uid, toUid, 0)
	r.Response.WriteTpl("private_chat.html", g.Map{
		"user_info": userInfo,
		"msg_list":  msgList,
		"room_id":   roomId,
	})
}

func Pagination(r *ghttp.Request) {
	roomId := r.GetString("room_id")
	toUid := r.GetString("uid")
	offset := r.GetString("offset")
	offsetInt, e := strconv.Atoi(offset)
	if e != nil || offsetInt <= 0 {
		offsetInt = 0
	}

	rooms := []string{"1", "2", "3", "4", "5", "6"}

	if !helper.InArray(roomId, rooms) {
		r.Response.WriteJson(g.Map{
			"code": 0,
			"data": map[string]interface{}{
				"list": nil,
			},
		})
		return
	}

	msgList := []map[string]interface{}{}
	if toUid != "" {
		userInfo := user_service.GetUserInfo(r)

		uid := strconv.Itoa(int(userInfo["uid"].(uint)))

		msgList = message_service.GetLimitPrivateMsg(uid, toUid, offsetInt)
	} else {
		msgList = message_service.GetLimitMsg(roomId, offsetInt)
	}

	r.Response.WriteJson(g.Map{
		"code": 0,
		"data": map[string]interface{}{
			"list": msgList,
		},
	})
}
