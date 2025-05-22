package utils

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

// essentially the same as calling uuid.NewV4
func GenerateUuid() (string, error) {
	// Create a Version 4 UUID.
	u2, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return u2.String(), nil
}

func ExtractUUIDFromUrl(path string, desiredUrl string) (string, error) {
	path = strings.Trim(path, "/")

	prefix := desiredUrl
	if !strings.HasPrefix(path, prefix) {
		return "", fmt.Errorf("path does not start with expected prefix")
	}

	remaining := strings.TrimPrefix(path, prefix)
	remaining = strings.Trim(remaining, "/")

	parts := strings.SplitN(remaining, "/", 2)
	id := parts[0]

	return id, nil
}
