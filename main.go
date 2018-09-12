package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/const-subject/alpin-bot/socket"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	token := flag.String("token", "none", "Token by bot")
	port := flag.Int("port", 5000, "Socket port")
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

	b, err := ioutil.ReadFile("markers.json")
	if err != nil {
		panic(err.Error())
	}

	var markers []socket.Marker
	err = json.Unmarshal(b, &markers)
	if err != nil {
		panic(err.Error())
	}

	mm := socket.NewMarkerManager()
	for _, marker := range markers {
		log.Println("Marker:", marker.Name)
		mm.Add(marker)
	}

	s := socket.NewSocket(updates, *port, mm)

	s.Start()
}
