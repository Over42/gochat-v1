package room

import (
	"gochat/config"

	"context"

	"github.com/go-playground/validator/v10"
)

type Room struct {
	ID         string
	Name       string
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

type Hub struct {
	Rooms map[string]*Room
}

type Service interface {
	CreateRoom(ctx context.Context, req *CreateRoomReq) (*CreateRoomRes, error)
	DeleteRoom(ctx context.Context, req *DeleteRoomReq) error
	GetRooms(ctx context.Context) ([]GetRoomsRes, error)
	JoinRoom(ctx context.Context, req *JoinRoomReq) error
	GetClients(ctx context.Context, req *GetClientsReq) ([]GetClientsRes, error)
}

type Repository interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	DeleteRoom(ctx context.Context, id string) error
	GetRooms(ctx context.Context) ([]*Room, error)
	GetClients(ctx context.Context, roomId string) ([]*Client, error)
}

func NewRoom(id string, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func Init(cfg *config.Config, val *validator.Validate) *Handler {
	hub := NewHub()
	roomRep := NewRepository(hub)
	roomSvc := NewService(roomRep, cfg, val, hub)
	roomHdl := NewHandler(roomSvc, cfg)
	return roomHdl
}
