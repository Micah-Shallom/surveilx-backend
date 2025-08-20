package models

type RegisterVehicleInput struct {
	PlateNumber      string      `json:"plate_number" validate:"required"`
	Model            string      `json:"model" validate:"required"`
	Color            string      `json:"color" validate:"required"`
	Type             string      `json:"type" validate:"required,oneof=bus bike car"`
}

type LogVehicleInput struct {
	PlateNumber  string `json:"plate_number" validate:"required"`
	IsEntry      bool   `json:"is_entry"`
	Type         string `json:"type" validate:"oneof=bus bike car"`
	EntryPointID string `json:"entry_point_id,omitempty"`
	ExitPointID  string `json:"exit_point_id,omitempty"`
}
