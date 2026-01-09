package room

import (
	"bingo/bingo"
	"bingo/config"
	"sync"
)

type RoomManager struct {
	rooms map[int]*Room
	mu    sync.RWMutex
}

type Room struct {
	ID        int
	Name      string
	Card      bingo.BingoCard
	Players   map[string]*Player
	isStarted bool
	mu        sync.RWMutex
}

func (rm *RoomManager) CreateRoom(name string) (*Room, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	roomId := len(rm.rooms) + 1
	room := &Room{
		ID:      roomId,
		Name:    name,
		Players: make(map[string]*Player),
	}
	rm.rooms[roomId] = room

	return room, nil
}

func (rm *RoomManager) DeleteRoom(roomId int) (bool, error) {
	rm.mu.RLock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomId]; !exists {
		return false, nil
	}
	// 試合中は削除不可
	if rm.rooms[roomId].isStarted {
		return false, nil
	}

	delete(rm.rooms, roomId)
	return true, nil
}

func (rm *RoomManager) GetRoom(roomId int) (*Room, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[roomId], nil
}

func (r *Room) AddPlayer(discordUserId string, name string, color string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RLocker().Unlock()

	var c config.Config
	maxPlayer := c.GetConfig().MaxPlayer

	if len(r.Players)+1 > maxPlayer {
		return false, nil
	}

	var progress [25]bool

	player := &Player{
		DiscordUserID: discordUserId,
		Name:          name,
		Color:         color,
		Progress:      progress,
		Rating:        0,
	}

	r.Players[discordUserId] = player
	return true, nil
}

func (r *Room) DeletePlayer(discordUserId string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Players[discordUserId]; !exists {
		return false, nil
	}

	delete(r.Players, discordUserId)
	return true, nil
}
func (r *Room) GetPlayer(discordUserId string) (*Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.Players[discordUserId] == nil {

	}

	return r.Players[discordUserId], nil
}
