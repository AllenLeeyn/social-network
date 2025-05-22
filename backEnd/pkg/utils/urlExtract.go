package utils

import (
	"strings"
)

func ExtractFromUrl(path string, desiredUrl string) (string, string) {
	if strings.HasPrefix(path, "/"+desiredUrl+"/") {
		extractedStr := strings.TrimPrefix(path, "/"+desiredUrl+"/")
		return extractedStr, ""
	} else {
		return "", "not found"
	}
}
