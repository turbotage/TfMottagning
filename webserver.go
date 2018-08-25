package main

import (
	"database/sql"
	"log"
	"net/http"
	"flag"

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

	var ip = flag.String("database_ip", "127.0.0.1", "the ip to the database")
	var port = flag.String("database_port", "5555", "the port to the database")
	var user = flag.String("database_username", "turbotage", "the username to the database")
	var password = flag.String("database_password", "klassuger", "the password to the database")
	var dbname = flag.String("database_name", "tfnolla", "the database name")

	db, err := sql.Open("mysql", *user + ":" + *password + "@tcp(" + *ip + ":" + *port + ")/" + *dbname)

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
			err := json.Unmarshal([]byte(msg), winRequest)
			if len(winRequest.ChallengeID) > 20 {

			}
			if len(winRequest.NollaName) > 20 {

			}
			if len(winRequest.Password) > 20 {

			}
			rows, err := db.Query("SELECT name FROM users WHERE ChallengeID=?", winRequest.ChallengeID);
			log.Println(rows)

			db.Query("INSERT INTO Challenges () VALUES(")

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
