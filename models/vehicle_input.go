package models

type RegisterVehicleInput struct {
	PlateNumber string `json:"plate_number" validate:"required"`
	Model       string `json:"model" validate:"required"`
	Color       string `json:"color" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=bus bike car"`
}

type LogVehicleInput struct {
	PlateNumber string `json:"plate_number" validate:"required"`
	IsEntry     bool   `json:"is_entry" validate:"required"`
}
