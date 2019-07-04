package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var SQL *sql.DB
var clients = make(map[int]map[*websocket.Conn]bool)
var broadcast = make(map[int](chan Push))

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var err error
	SQL, err = sql.Open("mysql", "------------/hackathon?parseTime=true")
	fs := http.FileServer(http.Dir("../public"))

	http.Handle("/", fs)
	http.HandleFunc("/ws", sendmesssage)

	go handleMessages()
	err = http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}

func sendmesssage(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		var user_id, group_id int
		if msg.Session == "" {
			log.Printf("no session id")
			return
		}
		user_id, group_id, err = GetUserGroup(msg.Session)
		_, ok := clients[group_id]
		if !ok {
			clients[group_id] = make(map[*websocket.Conn]bool)
		}
		clients[group_id][ws] = true
		log.Println(len(clients))
		log.Println(len(clients[group_id]))
		if err != nil {
			log.Printf("error in sendmessage: %v", err)
			delete(clients[group_id], ws)
			break
		}
		var msgs Push
		msgs.Username, err = GetUserName(user_id)
		msgs.Message = msg.Message
		msgs.User_id = user_id
		log.Println(msgs.Username + " said " + msgs.Message)
		broadcast[group_id] = make(chan Push)
		broadcast[group_id] <- msgs
	}
}

func handleMessages() {

	for {
		for group_id := range broadcast {
			log.Println(group_id)
			msg := <-broadcast[group_id]
			_, ok := broadcast[group_id]
			if ok {
				delete(broadcast, group_id)
			}
			log.Println(msg)
			for client, value := range clients[group_id] {
				log.Println(value)
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error in handler: %v", err)
					client.Close()
					delete(clients[group_id], client)
				}
			}
		}
	}
}
