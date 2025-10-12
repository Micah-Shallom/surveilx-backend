// controllers/websocket_controller.go - UPDATED
package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"survielx-backend/connections"
	"survielx-backend/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	// Add handshake timeout
	HandshakeTimeout: 10 * time.Second,
}

func WSHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !parsedToken.Valid {
		log.Println("Invalid token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Println("Invalid token claims")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	log.Println("WS connection established for user:", userID)

	// CRITICAL: Store the connection BEFORE defer
	connections.StoreClient(userID, conn)

	// Set up connection parameters
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		log.Println("Pong received from user:", userID)
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Send connection acknowledgment
	ackMsg := map[string]interface{}{
		"type":    "connection_ack",
		"message": "Connected successfully",
		"user_id": userID,
	}
	if err := conn.WriteJSON(ackMsg); err != nil {
		log.Println("Failed to send connection ack:", err)
		connections.DeleteClient(userID)
		conn.Close()
		return
	}

	// Start ping ticker in a goroutine
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("Ping error for user", userID, ":", err)
					return
				}
				log.Println("Ping sent to user:", userID)
			case <-done:
				return
			}
		}
	}()

	// Clean up when done
	defer func() {
		close(done)
		connections.DeleteClient(userID)
		conn.Close()
		log.Println("WS connection closed for user:", userID)
	}()

	// Read messages from client
	for {
		var message map[string]interface{}
		err := conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", userID, err)
			}
			break
		}

		// Handle ping from client
		if msgType, ok := message["type"].(string); ok && msgType == "ping" {
			pongMsg := map[string]interface{}{
				"type": "pong",
			}
			if err := conn.WriteJSON(pongMsg); err != nil {
				log.Println("Failed to send pong:", err)
				break
			}
			log.Println("Pong sent to user:", userID)
			continue
		}

		// Handle user responses
		rawMessage, _ := json.Marshal(message)
		services.HandleUserResponse(userID, rawMessage)
	}
}
