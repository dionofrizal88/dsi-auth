package util

import (
	"fmt"
	"net/http"
)

// GetType is a func to return its content as bytes and the content type.
func GetType(fileBytes []byte) ([]byte, string, error) {
	if len(fileBytes) == 0 {
		return nil, "", fmt.Errorf("file content is empty")
	}

	// Detect the content type
	contentType := http.DetectContentType(fileBytes)

	return fileBytes, contentType, nil
}
