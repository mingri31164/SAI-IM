package test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func serverWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		message2 := append(message, " from server"...)
		err = c.WriteMessage(mt, message2)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	//http.HandleFunc("/ws", serverWs)
	//fmt.Println("Server started on port 1234")
	//log.Fatal(http.ListenAndServe(":1234", nil))

	m := 9 >> 1
	fmt.Println(m)
}
