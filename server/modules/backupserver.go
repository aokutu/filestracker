package modules

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// ---------------------- HTTP HANDLERS ----------------------

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form with max memory
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        http.Error(w, "cannot parse form: "+err.Error(), 400)
        return
    }

    // Get filepath FIRST (before extracting file)
    filePath := r.FormValue("filepath")
    if filePath == "" {
        http.Error(w, "missing filepath field", 400)
        return
    }

    // Clean the path
    filePath = filepath.Clean("/" + filePath)
    filePath = strings.TrimPrefix(filePath, "/")

    // Now get the file
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "cannot read file: "+err.Error(), 400)
        return
    }
    defer file.Close()

    // Build destination
    dstPath := filepath.Join("backupdir", filePath)
    dir := filepath.Dir(dstPath)

    fmt.Printf("DEBUG: filepath=%s, dstPath=%s, dir=%s, header.Filename=%s\n", 
        filePath, dstPath, dir, header.Filename)

    // Create directories
    err = os.MkdirAll(dir, 0755)
    if err != nil {
        http.Error(w, "cannot create directory: "+err.Error(), 500)
        return
    }

    // Create destination file
    dst, err := os.Create(dstPath)
    if err != nil {
        http.Error(w, "cannot create file: "+err.Error(), 500)
        return
    }
    defer dst.Close()

    // Copy with size check
    written, err := io.Copy(dst, file)
    if err != nil {
        http.Error(w, "cannot save file: "+err.Error(), 500)
        return
    }

    fmt.Printf("DEBUG: wrote %d bytes to %s\n", written, dstPath)

    fmt.Fprintf(w, "uploaded: %s (%d bytes)\n", filePath, written)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the filename from query parameter
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}

	// Build the full path inside backupdir
	filePath := "backupdir/" + filename

	// Attempt to remove the file
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Unable to delete file", http.StatusInternalServerError)
		}
		return
	}

	Writelogs("DELETED " + filename)
	// Respond with success message
	fmt.Fprintf(w, "File '%s' deleted successfully!", filename)
}

func LogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the backup directory exists
	err := os.MkdirAll("backupdir", 0755)
	if err != nil {
		http.Error(w, "Unable to create backup directory", http.StatusInternalServerError)
		return
	}

	// Open the logs file in append mode (create if it doesn't exist)
	file, err := os.OpenFile("backupdir/logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Unable to open logs file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Copy the request body (uploaded file) into the file, appending to existing content
	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Failed to append logs file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logs appended successfully"))
}

// ---------------------- TCP HANDLER ----------------------

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read until newline
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed or error:", err)
			log.Println(line)
			return
		}
		fmt.Println("Received from client:", line)
	}
}

// ---------------------- MAIN ----------------------

func Writelogs(newdata string) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// Get computer/hostname
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return
	}

	currenttime := time.Now()
	currentdatetime := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]",
		currenttime.Year(), currenttime.Month(), currenttime.Day(),
		currenttime.Hour(), currenttime.Minute(), currenttime.Second())

	prevlogs := []byte{}

	prevlogs = append(prevlogs, []byte("["+currentUser.Username+"]"+"["+hostname+"]"+"["+currentdatetime+"]"+"["+newdata+"]"+"\n")...)

	os.WriteFile("logs", prevlogs, 0644)

	cmd := exec.Command(
		"curl",
		"-X", "PUT",
		"-T", "logs",
		"http://localhost:8080/logs",
	)

	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
