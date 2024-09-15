package room

import (
	"gochatv1/config"

	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

type service struct {
	repository Repository
	config     *config.Config
	validate   *validator.Validate
	hub        *Hub
}

func NewService(repo Repository, cfg *config.Config, val *validator.Validate, hub *Hub) Service {
	return &service{
		repo,
		cfg,
		val,
		hub,
	}
}

type CreateRoomReq struct {
	Name string `json:"name" validate:"required,min=3"`
}

type CreateRoomRes struct {
	ID string `json:"id"`
}

func (s *service) CreateRoom(ctx context.Context, req *CreateRoomReq) (*CreateRoomRes, error) {
	err := s.validate.Struct(req)
	if err != nil {
		return nil, err
	}

	context, cancel := context.WithTimeout(ctx, s.config.DBTimeout)
	defer cancel()

	id := ulid.Make().String()
	newRoom := NewRoom(id, req.Name)
	room, err := s.repository.CreateRoom(context, newRoom)
	if err != nil {
		return nil, err
	}

	go newRoom.run(s.hub)

	res := &CreateRoomRes{ID: room.ID}

	return res, nil
}

type DeleteRoomReq struct {
	ID string `json:"id"`
}

func (s *service) DeleteRoom(ctx context.Context, req *DeleteRoomReq) error {
	context, cancel := context.WithTimeout(ctx, s.config.DBTimeout)
	defer cancel()

	err := s.repository.DeleteRoom(context, req.ID)
	if err != nil {
		return err
	}

	return nil
}

type GetRoomsRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *service) GetRooms(ctx context.Context) ([]GetRoomsRes, error) {
	context, cancel := context.WithTimeout(ctx, s.config.DBTimeout)
	defer cancel()

	rooms, err := s.repository.GetRooms(context)
	if err != nil {
		return nil, err
	}

	res := make([]GetRoomsRes, 0)
	for _, r := range rooms {
		res = append(res, GetRoomsRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return res, nil
}

type JoinRoomReq struct {
	Conn     *websocket.Conn
	UserID   string `json:"userId"   validate:"required"`
	RoomID   string `json:"roomId"   validate:"required"`
	Username string `json:"username" validate:"required"`
}

func (s *service) JoinRoom(ctx context.Context, req *JoinRoomReq) error {
	err := s.validate.Struct(req)
	if err != nil {
		return err
	}

	client := &Client{
		Conn:     req.Conn,
		Message:  make(chan *Message, 10),
		UserID:   req.UserID,
		RoomID:   req.RoomID,
		Username: req.Username,
	}

	msg := &Message{
		Content:  "New user has joined",
		RoomID:   req.RoomID,
		Username: req.Username,
	}

	s.hub.Rooms[req.RoomID].Register <- client
	s.hub.Rooms[req.RoomID].Broadcast <- msg

	go client.writeMessage()
	go client.readMessage(s.hub.Rooms[req.RoomID])

	return nil
}

type GetClientsReq struct {
	RoomID string `json:"id" validate:"required"`
}

type GetClientsRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (s *service) GetClients(ctx context.Context, req *GetClientsReq) ([]GetClientsRes, error) {
	err := s.validate.Struct(req)
	if err != nil {
		return nil, err
	}

	context, cancel := context.WithTimeout(ctx, s.config.DBTimeout)
	defer cancel()

	clients, err := s.repository.GetClients(context, req.RoomID)
	if err != nil {
		return nil, err
	}

	res := make([]GetClientsRes, 0)
	for _, c := range clients {
		res = append(res, GetClientsRes{
			ID:       c.UserID,
			Username: c.Username,
		})
	}

	return res, nil
}

func (r *Room) run(hub *Hub) {
	for {
		select {
		case client := <-r.Register:
			hub.Rooms[client.RoomID].Clients[client.UserID] = client

		case client := <-r.Unregister:
			if _, ok := hub.Rooms[client.RoomID].Clients[client.UserID]; ok {
				r.Broadcast <- &Message{
					Content:  "User left the chat",
					RoomID:   client.RoomID,
					Username: client.Username,
				}
				delete(hub.Rooms[client.RoomID].Clients, client.UserID)
				close(client.Message)
			}

		case msg := <-r.Broadcast:
			for _, client := range hub.Rooms[msg.RoomID].Clients {
				client.Message <- msg
			}
		}
	}
}
