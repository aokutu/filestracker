package modules

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Sendfiles uploads a single file or recursively uploads all files in a folder.
// Folder structure is preserved: "A/B/filename" → backupdir/A/B/filename
func Sendfiles(relativePath string) {
	relativePath = filepath.ToSlash(relativePath)
	relativePath = strings.TrimPrefix(relativePath, "./")
	relativePath = strings.TrimPrefix(relativePath, "storage/")

	fullPath := filepath.Join("storage", relativePath)

	info, err := os.Stat(fullPath)
	if err != nil {
		fmt.Printf("ERROR: Cannot access %s: %v\n", fullPath, err)
		return
	}

	// If it's a directory, walk all files inside it
	if info.IsDir() {
		fmt.Printf("📁 Uploading folder: %s\n", relativePath)

		fileCount := 0
		err = filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("ERROR accessing %s: %v\n", path, err)
				return nil
			}
			if d.IsDir() {
				return nil
			}

			// Get relative path from storage root (preserves folder structure)
			relPath, err := filepath.Rel("storage", path)
			if err != nil {
				fmt.Printf("ERROR calculating relative path for %s: %v\n", path, err)
				return nil
			}
			relPath = filepath.ToSlash(relPath)

			fmt.Printf("  📄 Found file #%d: %s\n", fileCount+1, relPath)
			uploadSingleFile(relPath)
			fileCount++
			return nil
		})
		if err != nil {
			fmt.Printf("ERROR walking directory: %v\n", err)
		}
		fmt.Printf("✅ Finished uploading folder: %s (%d files)\n", relativePath, fileCount)
		return
	}

	// Single file — preserve any folder prefix in the relative path
	fmt.Printf("📄 Uploading single file: %s\n", relativePath)
	uploadSingleFile(relativePath)
}

// uploadSingleFile handles the actual curl upload.
// relativePath includes folder structure, e.g. "A/B/filename"
func uploadSingleFile(relativePath string) {
	uploadURL := Getaddress()
	fullPath := filepath.Join("storage", relativePath)

	fmt.Printf("    → Uploading: filepath=%s | local=%s\n", relativePath, fullPath)

	// Verify file exists locally before uploading
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("    ❌ ERROR: Local file not found: %s\n", fullPath)
		return
	}

	cmd := exec.Command("curl",
		"-s", 
		"-H", "X-API-Key: 4664" ,                              // silent mode (less noise)
		"-w", "\nHTTP_CODE:%{http_code}\n", // get HTTP status
		"-F", "filepath="+relativePath, // A/B/filename — server creates backupdir/A/B/filename
		"-F", "file=@"+fullPath,
		uploadURL,
	)

	output, err := cmd.CombinedOutput()
	fmt.Printf("    ← Server response: %s\n", string(output))

	if err != nil {
		fmt.Printf("    ❌ Error uploading %s: %v\n", relativePath, err)
	} else {
		fmt.Printf("    ✅ Uploaded %s\n", relativePath)
	}

	Sendlogs()
}

func Sendlogs() {
	conn, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Println("Connection failed:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(userdetails() + "\t" + Timestamp()))
	if err != nil {
		log.Println("Write failed:", err)
		return
	}
}
