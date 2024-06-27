// handlers/room.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/livekit/protocol/livekit"

	"LiveKitBackend/services"
)

type RoomHandler struct {
	roomService services.RoomService
}

func NewRoomHandler(roomService services.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateRoomHandler called")

	var roomName struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&roomName)
	if err != nil {
		log.Printf("Failed to parse request body: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	room, err := h.roomService.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:         roomName.Name,
		EmptyTimeout: 60 * 60 * 60,
	})
	if err != nil {
		log.Printf("Failed to create room: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
		Room    string `json:"room"`
	}{
		Message: "Room created successfully",
		Room:    room.Name,
	}

	log.Printf("Room created: %s", room.Name)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *RoomHandler) ListRoomsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ListRoomsHandler called")

	rooms, err := h.roomService.ListRooms(context.Background(), &livekit.ListRoomsRequest{})
	if err != nil {
		log.Printf("Failed to list rooms: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Rooms listed: %v", rooms.Rooms)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rooms.Rooms)
}

func (h *RoomHandler) DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteRoomHandler called")

	var roomName struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&roomName)
	if err != nil {
		log.Printf("Failed to parse request body: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	_, err = h.roomService.DeleteRoom(context.Background(), &livekit.DeleteRoomRequest{
		Room: roomName.Name,
	})
	if err != nil {
		log.Printf("Failed to delete room: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Room deleted successfully",
	}

	log.Printf("Room deleted: %s", roomName.Name)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *RoomHandler) CreateRoomTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateRoomTokenHandler called")

	var req struct {
		Room     string `json:"room"`
		Identity string `json:"identity"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Failed to parse request body: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	token, err := h.roomService.CreateRoomToken(req.Room, req.Identity, time.Duration(60)*time.Hour)
	if err != nil {
		log.Printf("Failed to create room token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	log.Printf("Room token created for room: %s, identity: %s", req.Room, req.Identity)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
