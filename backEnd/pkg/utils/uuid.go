package utils

import (
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

func ExtractUUIDFromUrl(path string, desiredUrl string) (string, string) {
	if strings.HasPrefix(path, "/"+desiredUrl+"/") {
		id := strings.TrimPrefix(path, "/"+desiredUrl+"/")
		return id, ""
	} else {
		return "", "not found"
	}
}
