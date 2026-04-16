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
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create 'backupdir' directory if it doesn't exist
	err = os.MkdirAll("backupdir", os.ModePerm)
	if err != nil {
		http.Error(w, "Unable to create backupdir folder", http.StatusInternalServerError)
		return
	}

	// Save the file in the 'backupdir' folder
	outFile, err := os.Create("backupdir/" + header.Filename)
	if err != nil {
		http.Error(w, "Unable to create file in backupdir", http.StatusInternalServerError)
		return
	}

	Writelogs( "EDITED " +  header.Filename   ) 
	defer outFile.Close()

	// Copy the contents of the uploaded file to the new file
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	fmt.Fprintf(w, "File uploaded successfully to backupdir!")
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

func StartTCPServer() {
	listen, err := net.Listen("tcp", "localhost:3001")
	if err != nil {
		fmt.Println("TCP Listen error:", err)
		return
	}
	fmt.Println("TCP server running on :3001")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("TCP Accept error:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go HandleConnection(conn) // handle each client concurrently
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




