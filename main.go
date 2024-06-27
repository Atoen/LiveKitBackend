package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"LiveKitBackend/handlers"
	"LiveKitBackend/services"
	"github.com/gorilla/mux"
	lksdk "github.com/livekit/server-sdk-go"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fmt.Println("Starting HTTP server...")

	host := "http://localhost:7880"
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")
	client := lksdk.NewRoomServiceClient(host, apiKey, apiSecret)
	roomService := services.NewRoomServiceClient(client, apiKey, apiSecret)

	roomHandler := handlers.NewRoomHandler(roomService)

	router := mux.NewRouter()
	router.HandleFunc("/create-room", roomHandler.CreateRoomHandler).Methods("POST")
	router.HandleFunc("/list-rooms", roomHandler.ListRoomsHandler).Methods("GET")
	router.HandleFunc("/delete-room", roomHandler.DeleteRoomHandler).Methods("POST")
	router.HandleFunc("/create-room-token", roomHandler.CreateRoomTokenHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":501", router))
}
