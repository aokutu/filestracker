package modules

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"time"
	"os/exec"
	
)

// ---------------------- HTTP HANDLERS ----------------------



func UploadHandler(w http.ResponseWriter, r *http.Request) {

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "cannot read file", 400)
        return
    }
    defer file.Close()

    // IMPORTANT: ensure folder exists
    os.MkdirAll("backupdir", 0755)

    dstPath := "backupdir/" + header.Filename

    dst, err := os.Create(dstPath)
    if err != nil {
        http.Error(w, "cannot create file", 500)
        return
    }
    defer dst.Close()

    _, err = io.Copy(dst, file)
    if err != nil {
        http.Error(w, "cannot save file", 500)
        return
    }

    fmt.Fprintf(w, "uploaded: %s\n", header.Filename)
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


	 Writelogs( "DELETED " +  filename   )    
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



func Writelogs(newdata string){

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

	 prevlogs =append(prevlogs, []byte(  "[" + currentUser.Username  + "]"  +  "[" +  hostname  + "]"   +  "[" + currentdatetime  + "]"  +  "[" +  newdata +  "]"   +  "\n") ...)

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




