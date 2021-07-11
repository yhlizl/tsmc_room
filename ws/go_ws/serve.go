package go_ws

import (
	"encoding/json"
	"go-gofram-chat/app/models"
	"go-gofram-chat/app/service/helper"
	"go-gofram-chat/ws"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gorilla/websocket"
)

// client status
type wsClients struct {
	Conn *websocket.Conn `json:"conn"`

	RemoteAddr string `json:"remote_addr"`

	Uid float64 `json:"uid"`

	Username string `json:"username"`

	RoomId string `json:"room_id"`

	AvatarId string `json:"avatar_id"`
}

// client & serve  message tx
type msg struct {
	Status int             `json:"status"`
	Data   interface{}     `json:"data"`
	Conn   *websocket.Conn `json:"conn"`
}

// var init
var (
	wsUpgrader = websocket.Upgrader{}

	clientMsg = msg{}

	mutex = sync.Mutex{}

	//rooms = [roomCount + 1][]wsClients{}
	rooms = make(map[int][]wsClients)

	enterRooms = make(chan wsClients)

	sMsg = make(chan msg)

	offline = make(chan *websocket.Conn)

	chNotify = make(chan int, 1)
)

// define type of msg status
const msgTypeOnline = 1        // online
const msgTypeOffline = 2       // offline
const msgTypeSend = 3          // sending msg
const msgTypeGetOnlineUser = 4 // getting online user list
const msgTypePrivateChat = 5   // pivate msg

const roomCount = 6 // room count

type GoServe struct {
	ws.ServeInterface
}

func (goServe *GoServe) RunWs(r *ghttp.Request) {
	// 使用 channel goroutine
	Run(r)
}

func (goServe *GoServe) GetOnlineUserCount() int {
	return GetOnlineUserCount()
}

func (goServe *GoServe) GetOnlineRoomUserCount(roomId int) int {
	return GetOnlineRoomUserCount(roomId)
}

func Run(r *ghttp.Request) {

	// @see https://github.com/gorilla/websocket/issues/523
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, _ := wsUpgrader.Upgrade(r.Response.RawWriter(), r.Request, nil)

	defer c.Close()

	go read(c)
	go write()

	select {}

}

func read(c *websocket.Conn) {

	defer func() {
		//catch read panic
		if err := recover(); err != nil {
			log.Println("read error : ", err)
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		// log.Println("client message", string(message),c.RemoteAddr())
		if err != nil { // 離線通知
			offline <- c
			log.Println("ReadMessage error1", err)
			return
		}
		serveMsgStr := message
		// 處理心跳響應 , heartbeat為與客戶端約定的值
		if string(serveMsgStr) == `heartbeat` {
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":0,"data":"heartbeat ok"}`))
			continue
		}
		json.Unmarshal(message, &clientMsg)

		// log.Println("來自客戶端的消息", clientMsg,c.RemoteAddr())

		if clientMsg.Data != nil {
			if clientMsg.Status == msgTypeOnline { // 進入房間，建立連接
				roomId, _ := getRoomId()

				enterRooms <- wsClients{
					Conn:       c,
					RemoteAddr: c.RemoteAddr().String(),
					Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
					Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
					RoomId:     roomId,
					AvatarId:   clientMsg.Data.(map[string]interface{})["avatar_id"].(string),
				}
			}

			_, serveMsg := formatServeMsgStr(clientMsg.Status, c)
			sMsg <- serveMsg
		}
	}
}

func write() {

	defer func() {
		//捕獲write拋出的panic
		if err := recover(); err != nil {
			log.Println("write發生錯誤", err)
		}
	}()

	for {
		select {
		case r := <-enterRooms:
			handleConnClients(r.Conn)
		case cl := <-sMsg:
			serveMsgStr, _ := json.Marshal(cl)
			switch cl.Status {
			case msgTypeOnline, msgTypeSend:
				notify(cl.Conn, string(serveMsgStr))
			case msgTypeGetOnlineUser:
				chNotify <- 1
				cl.Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				<-chNotify
			case msgTypePrivateChat:
				chNotify <- 1
				toC := findToUserCoonClient()
				if toC != nil {
					toC.(wsClients).Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				}
				<-chNotify
			}
		case o := <-offline:
			disconnect(o)
		}
	}
}

func handleConnClients(c *websocket.Conn) {
	roomId, roomIdInt := getRoomId()

	for cKey, wcl := range rooms[roomIdInt] {
		if wcl.Uid == clientMsg.Data.(map[string]interface{})["uid"].(float64) {
			// 通知當前用戶下線
			wcl.Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
			mutex.Lock()
			rooms[roomIdInt] = append(rooms[roomIdInt][:cKey], rooms[roomIdInt][cKey+1:]...)
			mutex.Unlock()
			wcl.Conn.Close()
		}
	}

	mutex.Lock()
	rooms[roomIdInt] = append(rooms[roomIdInt], wsClients{
		Conn:       c,
		RemoteAddr: c.RemoteAddr().String(),
		Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
		Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
		RoomId:     roomId,
		AvatarId:   clientMsg.Data.(map[string]interface{})["avatar_id"].(string),
	})
	mutex.Unlock()
}

// 獲取私聊的用戶連接
func findToUserCoonClient() interface{} {
	_, roomIdInt := getRoomId()
	toUserUid := clientMsg.Data.(map[string]interface{})["to_uid"].(string)

	for _, c := range rooms[roomIdInt] {
		stringUid := strconv.FormatFloat(c.Uid, 'f', -1, 64)
		if stringUid == toUserUid {
			return c
		}
	}

	return nil
}

// 統一消息發放
func notify(conn *websocket.Conn, msg string) {
	chNotify <- 1 // 利用channel阻塞 避免併發去對同一個連接發送消息出現panic: concurrent write to websocket connection這樣的異常
	_, roomIdInt := getRoomId()
	for _, con := range rooms[roomIdInt] {
		if con.RemoteAddr != conn.RemoteAddr().String() {
			con.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
	<-chNotify
}

// 離線通知
func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomId()
	for index, con := range rooms[roomIdInt] {
		if con.RemoteAddr == conn.RemoteAddr().String() {
			data := map[string]interface{}{
				"username": con.Username,
				"uid":      con.Uid,
				"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
			}

			jsonStrServeMsg := msg{
				Status: msgTypeOffline,
				Data:   data,
			}
			serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

			disMsg := string(serveMsgStr)

			mutex.Lock()
			rooms[roomIdInt] = append(rooms[roomIdInt][:index], rooms[roomIdInt][index+1:]...)
			mutex.Unlock()
			con.Conn.Close()
			notify(conn, disMsg)
		}
	}
}

// 格式化傳送給客戶端的消息數據
func formatServeMsgStr(status int, conn *websocket.Conn) ([]byte, msg) {

	roomId, roomIdInt := getRoomId()

	data := map[string]interface{}{
		"username": clientMsg.Data.(map[string]interface{})["username"].(string),
		"uid":      clientMsg.Data.(map[string]interface{})["uid"].(float64),
		"room_id":  roomId,
		"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
	}

	if status == msgTypeSend || status == msgTypePrivateChat {
		data["avatar_id"] = clientMsg.Data.(map[string]interface{})["avatar_id"].(string)
		content := clientMsg.Data.(map[string]interface{})["content"].(string)

		data["content"] = content
		if helper.MbStrLen(content) > 800 {
			// 直接截斷
			data["content"] = string([]rune(content)[:800])
		}

		toUidStr := clientMsg.Data.(map[string]interface{})["to_uid"].(string)
		toUid, _ := strconv.Atoi(toUidStr)

		// 保存消息
		stringUid := strconv.FormatFloat(data["uid"].(float64), 'f', -1, 64)
		intUid, _ := strconv.Atoi(stringUid)

		if _, ok := clientMsg.Data.(map[string]interface{})["image_url"]; ok {
			// 存在圖片
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data["content"],
				"room_id":    data["room_id"],
				"image_url":  clientMsg.Data.(map[string]interface{})["image_url"].(string),
			})
		} else {
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"room_id":    data["room_id"],
				"content":    data["content"],
			})
		}

	}

	if status == msgTypeGetOnlineUser {
		ro := rooms[roomIdInt]
		data["count"] = len(ro)
		data["list"] = ro
	}

	jsonStrServeMsg := msg{
		Status: status,
		Data:   data,
		Conn:   conn,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

	return serveMsgStr, jsonStrServeMsg
}

func getRoomId() (string, int) {

	roomId := clientMsg.Data.(map[string]interface{})["room_id"].(string)

	roomIdInt, _ := strconv.Atoi(roomId)

	// room := clientMsg.Data.(map[string]interface{})["room_id"].(float64)

	// roomIdInt := int(room)
	// roomId := strconv.Itoa(roomIdInt)
	return roomId, roomIdInt
}

// =======================對外方法=====================================

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i <= roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}

func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
