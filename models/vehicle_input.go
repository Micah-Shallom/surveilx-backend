package models

type RegisterVehicleInput struct {
	PlateNumber string `json:"plate_number" validate:"required"`
	Model       string `json:"model" validate:"required"`
	Color       string `json:"color" validate:"required"`
}

type LogVehicleInput struct {
	PlateNumber string `json:"plate_number" validate:"required"`
	IsEntry     bool   `json:"is_entry"`
}
