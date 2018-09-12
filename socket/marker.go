package socket

import "sync"

type Marker struct {
	Name   string  `json:"name"`
	ChatId int64   `json:"chat_id"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
}

type MarkerManager struct {
	mu   sync.Mutex
	data []Marker
}

func NewMarkerManager() *MarkerManager {
	return &MarkerManager{}
}

func (mm *MarkerManager) Add(marker Marker) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.data = append(mm.data, marker)
}

func (mm *MarkerManager) Get() []Marker {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	return mm.data
}
