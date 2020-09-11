package main

import (
	"fmt"
	"github.com/chehsunliu/poker"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
  "bytes"
)

// User Variables
var logfile = "/log.txt"
var dir = "/"

// System Variables
var gamestarted = false
var availablecards = []string{}
var players = []string{}
var playerhands= [][]poker.Card{}
var gamestatus = "Waiting for players..."
var deck = poker.NewDeck()

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return 0 - 1
}

func buildBuffer(nodes []poker.Card) string {
    buf := &bytes.Buffer{}
    buf.WriteString("{")
    buf.WriteByte('"')
    buf.WriteString("cards")
    buf.WriteByte('"')
    buf.WriteString(":[")
    for i, v := range nodes {
        if i > 0 {
            buf.WriteByte(',')
        }
        buf.WriteByte('"')
        buf.WriteString(fmt.Sprintf("%v",v))
        buf.WriteByte('"')
    }
    buf.WriteString("]}")
    return buf.String()
}

func logtofile(text string) {
	f, _ := os.OpenFile(dir+logfile, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	logtime := time.Now().Format("[01-02-2006 15:04:05.0000]: ")
	f.WriteString(logtime + text + "<br>")
	fmt.Printf(logtime + text + "\n")
}

func truncatelogfile() {
	f, _ := os.OpenFile(dir+logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	f.Truncate(0)
	fmt.Fprintf(f, "")
	f.Close()
}

func page(w http.ResponseWriter, req *http.Request) {
	dat, _ := ioutil.ReadFile(dir + "/index.html")
	fmt.Fprintf(w, string(dat))
	logtofile("Served main page to X.X.X." + strings.Split(string(req.RemoteAddr), ".")[3])
}

func join(w http.ResponseWriter, req *http.Request) {
	logtofile("Player associated to X.X.X." + strings.Split(string(req.RemoteAddr), ".")[3] + " has joined!")
	fmt.Fprintf(w, "Joined")
	if Find(players, string(strings.Split(string(req.RemoteAddr), ":")[0])) == 0-1 {
		players = append(players, string(strings.Split(string(req.RemoteAddr), ":")[0]))
		logtofile("New player")
	} else {
		logtofile("Player rejoining...")
	}
}

func gethand(w http.ResponseWriter, req *http.Request) {
  if(gamestarted){
  playerindex := Find(players, string(strings.Split(string(req.RemoteAddr), ":")[0]))
	if (playerindex != 0-1 ){
    fmt.Fprintf(w, buildBuffer(playerhands[playerindex]))
	}
}
}

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, gamestatus)
}

func round() {
	gamestatus = "Round 1"
	// Make new deck;
	deck = poker.NewDeck()
  deck.Shuffle();
	for index, player := range players {
		playerhands = append(playerhands, deck.Draw(2))
    logtofile("Generating hand for player "+string(index)+" at "+player)
	}
}

func start(w http.ResponseWriter, req *http.Request) {
	if !gamestarted {
		logtofile("Game starting!")
		gamestarted = true
		fmt.Fprintf(w, "Game started!")
		logtofile("Game started!")
		round()
	} else {
		fmt.Fprintf(w, "Game already started.")
	}
}

func main() {
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	truncatelogfile()
	logtofile("Server starting...")
	assets := http.FileServer(http.Dir(dir + "/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
	http.HandleFunc("/", page)
	http.HandleFunc("/start/", start)
	http.HandleFunc("/join/", join)
	http.HandleFunc("/status/", status)
  http.HandleFunc("/hand/",gethand)
	logtofile("Server started...")
	http.ListenAndServe(":11000", nil)
}
