package utility

import (
	"fmt"

	"github.com/google/uuid"
)

func ValidateUUID(uuidStr string) error {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return fmt.Errorf("invalid UUID format: %v", err)
	}
	return nil
}
