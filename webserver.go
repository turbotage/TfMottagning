package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
)

type AddWinnerRequest struct {
	ChallengeID string 'json:"challenge_id"'
	NollaName string 'json:"nolla"'
	Password string 'json:password'
}

func main() {

	db, err := sql.Open("mysql", "turbotage:klassuger@tcp(127.0.0.1:3306)/test")

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")

		so.On("query:challenges", func(msg string) {
			log.Println(msg)
		})

		so.On("request:add-winner", func(msg string) {
			var winRequest AddWinnerRequest
		})

		so.On("disconnection", func() {
			log.Println("client disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))

}
