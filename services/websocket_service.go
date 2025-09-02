package services

import (
	"encoding/json"
	"fmt"
	"log"
	"survielx-backend/connections"
	"survielx-backend/database"
	"survielx-backend/models"
	"survielx-backend/utility"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func notifyUserForExitConfirmation(pendingID, userID, token string) {
	var (
		 pending models.PendingVehicleExit
		 vehicle models.Vehicle
	)
	connInterface, ok := connections.Clients.Load(userID)
	if !ok {
		log.Println("User not connected:", userID)
		return
	}
	conn := connInterface.(*websocket.Conn)

	tx := database.DB.Where("id = ?", pendingID).First(&pending)
	if tx.Error != nil{
		log.Default().Println("unable to fetch pending vehicle information", tx.Error.Error())
		return
	}

	tx = database.DB.Where("plate_number = ?", pending.PlateNumber).First(&vehicle)
	if tx.Error != nil {
		log.Default().Println("unable to fetch vehicle information", tx.Error.Error())
		return
	}

	msg := map[string]any{
		"type":       "exit_confirmation",
		"message":    "Are you the one leaving the premises?",
		"pending_id": pendingID,
		"token":      token,
		"plateNumber": pending.PlateNumber,
		"vehicleName": vehicle.Model,
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("WS write error:", err)
	}
}

func HandleUserResponse(userID string, rawMessage []byte) {
	var resp struct {
		PendingID string `json:"pending_id"`
		Token     string `json:"token"`
		Confirmed bool   `json:"confirmed"`
	}
	if err := json.Unmarshal(rawMessage, &resp); err != nil {
		return
	}

	var pending models.PendingVehicleExit
	err := database.DB.Where("id = ?", resp.PendingID).First(&pending).Error
	if err != nil || pending.UserID != userID || pending.ResponseToken != resp.Token || pending.Status != "pending" {
		log.Println("Invalid pending exit response")
		return
	}

	if resp.Confirmed {
		activity := models.VehicleActivity{
			PlateNumber: pending.PlateNumber,
			VisitorType: models.VisitorTypeRegistered,
			IsEntry:     false,
			Timestamp:   pending.Timestamp,
			ExitPointID: &pending.ExitPointID,
			VehicleID:   &pending.VehicleID,
		}
		if err := database.DB.Create(&activity).Error; err != nil {
			log.Println("Failed to log exit:", err)
		}
		pending.Status = "confirmed"
	} else {
		pending.Status = "denied"
		notifySecurity(pending.PlateNumber, pending.Timestamp, pending.ExitPointID)
	}
	database.DB.Save(&pending)
}

func handleExitTimeout(pendingID string, db *gorm.DB, exitPointID string) {
	time.Sleep(20 * time.Second)
	var pending models.PendingVehicleExit
	err := db.Where("id = ?", pendingID).First(&pending).Error
	if err == nil && pending.Status == "pending" {
		pending.Status = "timed_out"
		db.Save(&pending)
		notifySecurity(pending.PlateNumber, pending.Timestamp, exitPointID)
	}
}

func notifySecurity(plateNumber string, timestamp time.Time, apID string) {
	var (
		user models.User
		exitPoint models.AccessExitPoint
	)


	tx := database.DB.Model(&models.User{}).Where("role = ?", "security").Find(&user)
	if tx.Error != nil {
		log.Println("Failed to find security user:", tx.Error)
		return
	}

	tx = database.DB.Model(&models.AccessExitPoint{}).Where("id = ?", apID).First(&exitPoint)
	if tx.Error != nil {
		log.Println("Failed to find exit point:", tx.Error)
		return
	}

	securityUserID := user.ID
	connInterface, ok := connections.Clients.Load(securityUserID)
	if ok {
		conn := connInterface.(*websocket.Conn)

		alert := map[string]interface{}{
			"type": "security_alert",
			"data": map[string]interface{}{
				"id":           utility.GenerateUUID(),
				"plate_number": plateNumber,
				"reason":       "Suspicious exit attempt detected",
				"timestamp":    timestamp.Format(time.RFC3339),
				"location":     exitPoint.Name,
			},
		}

		if err := conn.WriteJSON(alert); err != nil {
			log.Printf("Failed to send security alert: %v", err)
			sendEmailToSecurity(fmt.Sprintf("Alert: Suspicious exit for %s", plateNumber))
		}
	} else {
		// Fallback: Email/SMS to security
		sendEmailToSecurity(fmt.Sprintf("Alert: Suspicious exit for %s", plateNumber))
	}
}

func sendEmailToSecurity(content string) {
	log.Println("Sending email to security:", content)
}

func SendNotification(userID string, message []byte) error {
	if conn, ok := connections.GetClient(userID); ok {
		return conn.WriteMessage(websocket.TextMessage, message)
	}
	return fmt.Errorf("client not found")
}
