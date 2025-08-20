package models

type RegisterVehicleInput struct {
	PlateNumber string `json:"plate_number" validate:"required"`
	Model       string `json:"model" validate:"required"`
	Color       string `json:"color" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=bus bike car"`
}

type LogVehicleInput struct {
	PlateNumber  string `json:"plate_number" validate:"required"`
	IsEntry      bool   `json:"is_entry"`
	Type         string `json:"type" validate:"oneof=bus bike car"`
	EntryPointID string `json:"entry_point_id,omitempty"`
	ExitPointID  string `json:"exit_point_id,omitempty"`
}

type VehicleFilters struct {
	PlateNumber string
	Model       string
	Color       string
	Type        string
}

type VehicleOwnerProfileResponse struct {
	Vehicle    Vehicle                   `json:"vehicle"`
	Owner      VehicleOwnerInfo          `json:"owner"`
	Activities []VehicleActivityResponse `json:"activities"`
}

type VehicleOwnerInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	UserName string `json:"username"`
	FullName string `json:"full_name"`
}

type ActivityReportResponse struct {
	Data       ActivityReportData `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type ActivityReportData struct {
	Activities []VehicleActivityResponse `json:"activities"`
	Summary    map[string]interface{}    `json:"summary"`
	DateRange  map[string]string         `json:"date_range"`
}
