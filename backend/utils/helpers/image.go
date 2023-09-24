package helpers

import (
	"fmt"
	"time"
)

func GenerateUniqueFilename(fileExt string) string {
	// Use timestamp
	timestamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("%d%s", timestamp, fileExt)
	return fileName
}
