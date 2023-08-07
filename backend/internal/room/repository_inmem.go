package room

import (
	"context"
	"errors"
)

type repository struct {
	hub *Hub
}

func NewRepository(hub *Hub) Repository {
	return &repository{hub: hub}
}

func (r *repository) CreateRoom(ctx context.Context, room *Room) (*Room, error) {
	r.hub.Rooms[room.ID] = room

	return room, nil
}

func (r *repository) DeleteRoom(ctx context.Context, id string) error {
	if _, ok := r.hub.Rooms[id]; !ok {
		return errors.New("Room does not exist")
	}

	delete(r.hub.Rooms, id)
	return nil
}

func (r *repository) GetRooms(ctx context.Context) ([]*Room, error) {
	rooms := make([]*Room, 0, len(r.hub.Rooms))
	for _, r := range r.hub.Rooms {
		rooms = append(rooms, r)
	}

	return rooms, nil
}

func (r *repository) GetClients(ctx context.Context, roomId string) ([]*Client, error) {
	clients := make([]*Client, 0, len(r.hub.Rooms[roomId].Clients))
	for _, c := range r.hub.Rooms[roomId].Clients {
		clients = append(clients, c)
	}

	return clients, nil
}
