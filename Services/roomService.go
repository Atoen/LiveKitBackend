// services/room_service.go
package services

import (
	"context"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
	"time"
)

type RoomService interface {
	CreateRoom(ctx context.Context, req *livekit.CreateRoomRequest) (*livekit.Room, error)
	ListRooms(ctx context.Context, req *livekit.ListRoomsRequest) (*livekit.ListRoomsResponse, error)
	DeleteRoom(ctx context.Context, req *livekit.DeleteRoomRequest) (*livekit.DeleteRoomResponse, error)
	CreateRoomToken(room, identity string, validFor time.Duration) (string, error)
}

type RoomServiceClient struct {
	client    *lksdk.RoomServiceClient
	apiKey    string
	apiSecret string
}

func NewRoomServiceClient(client *lksdk.RoomServiceClient, apiKey, apiSecret string) *RoomServiceClient {
	return &RoomServiceClient{client: client, apiKey: apiKey, apiSecret: apiSecret}
}

func (r *RoomServiceClient) CreateRoom(ctx context.Context, req *livekit.CreateRoomRequest) (*livekit.Room, error) {
	return r.client.CreateRoom(ctx, req)
}

func (r *RoomServiceClient) ListRooms(ctx context.Context, req *livekit.ListRoomsRequest) (*livekit.ListRoomsResponse, error) {
	return r.client.ListRooms(ctx, req)
}

func (r *RoomServiceClient) DeleteRoom(ctx context.Context, req *livekit.DeleteRoomRequest) (*livekit.DeleteRoomResponse, error) {
	return r.client.DeleteRoom(ctx, req)
}

func (r *RoomServiceClient) CreateRoomToken(room, identity string, validFor time.Duration) (string, error) {
	at := auth.NewAccessToken(r.apiKey, r.apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(validFor)

	return at.ToJWT()
}
