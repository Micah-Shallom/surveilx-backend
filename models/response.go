package models

import "time"

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type VehicleResponse struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	User        UserResponse `json:"user"`
	PlateNumber string       `json:"plate_number"`
	Type        string       `json:"type"`
	Model       string       `json:"model"`
	Color       string       `json:"color"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type VehicleLogResponse struct {
	ID        string          `json:"id"`
	VehicleID string          `json:"vehicle_id"`
	Vehicle   VehicleResponse `json:"vehicle"`
	UserID    string          `json:"user_id"`
	User      UserResponse    `json:"user"`
	Timestamp time.Time       `json:"timestamp"`
	IsEntry   bool            `json:"is_entry"`
	CreatedAt time.Time       `json:"createdAt"`
}
