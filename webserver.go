package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

/*
type AddWinnerRequest struct {
	ChallengeID string 'json:"challenge_id"'
	NollaName string 'json:"nolla"'
	Password string 'json:"password"'
}
*/

var xlChallenges *excelize.File
var challengeSheet string

func getChallenge(challengeID int) string {
	return xlChallenges.GetCellValue(challengeSheet, "A"+strconv.Itoa(challengeID+1))
}

func main() {

	challengeSheet = "Blad1"
	//var ip = flag.String("database_ip", "127.0.0.1", "the ip to the database")
	//var port = flag.String("database_port", "5555", "the port to the database")
	//var user = flag.String("database_username", "turbotage", "the username to the database")
	//var password = flag.String("database_password", "klassuger", "the password to the database")
	//var dbname = flag.String("database_name", "tfnolla", "the database name")

	//db, err := sql.Open("mysql", *user + ":" + *password + "@tcp(" + *ip + ":" + *port + ")/" + *dbname)

	var err error

	xlChallenges, err = excelize.OpenFile("./Utmaningar.xlsx")

	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println(getChallenge(1))

	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		//join them to room
		c.Join("chat")
	})

	type Message struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	//handle custom event
	server.On("send", func(c *gosocketio.Channel, msg Message) string {
		//send event to all in room
		c.BroadcastTo("chat", "message", msg)
		log.Println(msg.Message)
		return "OK"
	})

	//setup http server
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	serveMux.Handle("/", http.FileServer(http.Dir("assets")))
	log.Panic(http.ListenAndServe(":5000", serveMux))

}
