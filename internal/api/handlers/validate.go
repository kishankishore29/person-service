package handlers

import "github.com/google/uuid"

// validateUUID Returns a boolean value to represent the validity of the passed uuid string
func validateUUID(in string) bool {
	_, err := uuid.Parse(in)
	return err == nil
}
