package main

import (
	"encoding/json"
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
var nRowsC int
var nColsC int

var xlPasswords *excelize.File
var passwordSheet string
var nRowsP int
var nColsP int

type AddWinnerRequest struct {
	ChallengeID int    `json:"challenge_id"`
	NollaName   string `json:"nolla"`
	Password    string `json:"password"`
}

type TableRowResponseData struct {
	CID     int    `json:"cid"`
	CName   string `json:"cname"`
	Points  int    `json:"points"`
	Nolla   string `json:"nolla"`
	Phadder string `json:"phadder"`
}

func getChallenge(challengeID int) string {
	return xlChallenges.GetCellValue(challengeSheet, "A"+strconv.Itoa(challengeID+1))
}

func getNolla(challengeID int) string {
	return xlChallenges.GetCellValue(challengeSheet, "C"+strconv.Itoa(challengeID+1))
}

func getTable() ([]TableRowResponseData, int) {
	tableRowResponse := make([]TableRowResponseData, nRowsC-1)
	length := 0
	rows := xlChallenges.GetRows(challengeSheet)
	for i := 1; i < nRowsC; i++ {
		tableRowResponse[i-1].CID = i
		tableRowResponse[i-1].CName = rows[i][0]
		tableRowResponse[i-1].Points, _ = strconv.Atoi(rows[i][1])
		tableRowResponse[i-1].Nolla = rows[i][2]
		tableRowResponse[i-1].Phadder = rows[i][3]
		length = i
	}
	return tableRowResponse, length
}

func findPassword(password string) int {
	for i := 1; i <= nRowsP; i++ {
		if password == xlPasswords.GetCellValue(passwordSheet, "B"+strconv.Itoa(i)) {
			return i
		}
	}
	return -1
}

func getPhadderFromPass(num int) string {
	return xlPasswords.GetCellValue(passwordSheet, "A"+strconv.Itoa(num))
}

func setNolla(num int, str string) {
	xlChallenges.SetCellValue(challengeSheet, "C"+strconv.Itoa(num+1), str)
}

func setPhadder(num int, str string) {
	xlChallenges.SetCellValue(challengeSheet, "D"+strconv.Itoa(num+1), str)
}

func main() {

	challengeSheet = "Blad1"
	nRowsC = 52
	nColsC = 5

	passwordSheet = "Blad1"
	nRowsP = 25
	nColsP = 2

	//var ip = flag.String("database_ip", "127.0.0.1", "the ip to the database")
	//var port = flag.String("database_port", "5555", "the port to the database")
	//var user = flag.String("database_username", "turbotage", "the username to the database")
	//var password = flag.String("database_password", "klassuger", "the password to the database")
	//var dbname = flag.String("database_name", "tfnolla", "the database name")

	//db, err := sql.Open("mysql", *user + ":" + *password + "@tcp(" + *ip + ":" + *port + ")/" + *dbname)

	var err error

	xlChallenges, err = excelize.OpenFile("./Utmaningar.xlsx")
	xlPasswords, err = excelize.OpenFile("./Passwords.xlsx")

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
		c.Emit("update", "jas")
	})

	server.On("table-req", func(c *gosocketio.Channel, num int) string {
		arr, _ := getTable()
		b, err := json.Marshal(arr)
		if err != nil {
			log.Println(err)
		}
		c.Emit("table-response", string(b))
		return "OK"
	})

	server.On("message", func(c *gosocketio.Channel, msg string) string {
		log.Println(msg)
		return "OK"
	})

	//handle custom event
	server.On("add-winner", func(c *gosocketio.Channel, winReq AddWinnerRequest) string {
		//send event to all in room

		challengeID := winReq.ChallengeID

		if (challengeID < 1) || (challengeID > (nRowsC - 1)) {
			c.Emit("alert", "Det specifierade Utmanings ID har ingen matchande utmaning")
			return "ERROR"
		}

		nolla := getNolla(challengeID)
		if len(nolla) != 0 {
			c.Emit("alert", nolla+" har redan klarat denna utmaning")
			return "ERROR"
		}
		if len(winReq.NollaName) > 50 {
			c.Emit("alert", "nollegruppsnamn är för långt")
			return "ERROR"
		}

		num := findPassword(winReq.Password)
		if num == -1 {
			c.Emit("alert", "inget matchande lösenord")
			return "ERROR"
		}
		phadder := getPhadderFromPass(num)
		setNolla(challengeID, winReq.NollaName)
		setPhadder(challengeID, phadder)

		xlChallenges.Save()

		c.Emit("update", "added winner")

		return "OK"
	})

	//setup http server
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	serveMux.Handle("/", http.FileServer(http.Dir("assets")))
	log.Panic(http.ListenAndServeTLS(":443", "server.crt", "server.key", serveMux))

}
