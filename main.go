package main

import (
	"flag"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	token := flag.String("token", "none", "Token by bot")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	NewSocket(updates)
}

func NewSocket(updates tgbotapi.UpdatesChannel) {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		so.Join("earth")
		log.Println("New client:", so.Id())
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		server.ServeHTTP(w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("./asset")))

	go http.ListenAndServe(":5000", nil)

	for update := range updates {
		if update.Message == nil ||
			!update.Message.Chat.IsSuperGroup() && !update.Message.Chat.IsGroup() {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)
		server.BroadcastTo("earth", "message",
			update.Message.From.FirstName,
			update.Message.Text,
			update.Message.Chat.Title)
	}
}
