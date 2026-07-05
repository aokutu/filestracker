package modules

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings" // add this
)

func DeleteFile(path string) {
	// Normalize path
	relPath := filepath.ToSlash(path)
	relPath = strings.TrimPrefix(relPath, "./")

	// Strip the storage/ prefix that the server doesn't know about
	relPath = strings.TrimPrefix(relPath, "storage/")

	// Build server URL
	baseURL := "http://localhost:8080/delete"

	// Add query parameter safely
	q := url.Values{}
	q.Set("file", relPath)
	finalURL := baseURL + "?" + q.Encode()

	fmt.Println("Deleting via:", finalURL)

	// Run curl command
	cmd := exec.Command("curl", "-H", "X-API-Key: 4664" , "-X", "DELETE", "-s", "-w", "\nHTTP_CODE:%{http_code}", finalURL)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to delete %s: %v\n", path, err)
		fmt.Printf("Curl output: %s\n", string(output))
		return
	}

	fmt.Printf("Successfully deleted %s from backupdir\n", path)
	fmt.Printf("Response: %s\n", string(output))
}
