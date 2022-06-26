package main

import (
	"fmt"
	"net/http"

	"github.com/leonardchinonso/chat_app_test/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Username: "",
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool1 := websocket.NewPool()
	pool2 := websocket.NewPool()

	go pool1.Start()
	go pool2.Start()

	http.HandleFunc("/room1", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool1, w, r)
	})

	http.HandleFunc("/room2", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool2, w, r)
	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
