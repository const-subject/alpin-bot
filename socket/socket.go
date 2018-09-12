package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"gopkg.in/telegram-bot-api.v4"
)

type Socket struct {
	updates tgbotapi.UpdatesChannel
	server  *socketio.Server
	port    int
	mm      *MarkerManager
	mu      sync.Mutex
	running bool
}

func NewSocket(updates tgbotapi.UpdatesChannel, port int, mm *MarkerManager) *Socket {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {

		j, err := json.Marshal(mm.Get())
		if err != nil {
			panic(err.Error())
		}
		so.Emit("load", string(j))
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

	return &Socket{
		updates: updates,
		server:  server,
		port:    port,
		mm:      mm,
	}
}

func (s *Socket) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go (func() { log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)) })()

	for update := range s.updates {
		s.mu.Lock()
		if !s.running {
			break
		}
		s.mu.Unlock()

		if update.Message == nil ||
			!update.Message.Chat.IsSuperGroup() && !update.Message.Chat.IsGroup() {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)
		s.server.BroadcastTo("earth", "message",
			update.Message.From.FirstName,
			update.Message.Text,
			update.Message.Chat.Title,
			update.Message.Chat.ID)
	}
}

func (s *Socket) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}
