package main

import (
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"

	"github.com/gorilla/websocket"

	pb "ride-sharing/shared/proto/driver"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("Missing userID in WebSocket connection")
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message from rider %s: %s", userID, string(message))
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("Missing userID in WebSocket connection")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Printf("Missing packageSlug in WebSocket connection")
		return
	}

	driverService, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer driverService.Close()

	req := &pb.RegisterDriverRequest{
		DriverId:    userID,
		PackageSlug: packageSlug,
	}

	driver, err := driverService.Client.RegisterDriver(r.Context(), req)
	if err != nil {
		log.Printf("Failed to register driver: %v", err)
		return
	}

	defer func() {
		_, err := driverService.Client.UnregisterDriver(r.Context(), req)
		if err != nil {
			log.Printf("Failed to unregister driver: %v", err)
		}
		log.Println("driver unregistred: ", userID)
	}()

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driver,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message from driver %s: %s", userID, string(message))
	}
}
