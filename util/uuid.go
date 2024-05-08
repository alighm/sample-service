package util

import hu "github.com/hashicorp/go-uuid"

// NewUUID creates a new UUID RFC spec
func NewUUID() string {
	u, _ := hu.GenerateUUID()
	return u
}

// IsUUID checks if the provided id is a valid UUID
func IsUUID(s string) bool {
	_, err := hu.ParseUUID(s)
	return err == nil
}
