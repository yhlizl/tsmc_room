package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func start() {
	var addr = flag.String("addr", "localhost:8322", "http service address")

	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	p := os.Args

	log.Println("Args", p)

	d := make(map[string]interface{})
	d["status"] = 1

	// string轉成int64：
	uid, _ := strconv.ParseFloat(p[1], 64)

	d["data"] = map[string]interface{}{
		"uid":       uid,
		"room_id":   "1",
		"avatar_id": "4",
		"username":  "suiji" + p[1],
	}

	c.WriteJSON(d)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
	}

}

func main() {
	start()
}
