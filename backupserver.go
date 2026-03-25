package main

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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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

 
func deleteHandler(w http.ResponseWriter, r *http.Request) {
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


	 Writelogs( "DELETED  " +  filename   )    
	// Respond with success message
	fmt.Fprintf(w, "File '%s' deleted successfully!", filename)
} 



// ---------------------- TCP HANDLER ----------------------

func handleConnection(conn net.Conn) {
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

func startTCPServer() {
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
		go handleConnection(conn) // handle each client concurrently
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


	 prevlogs, err := os.ReadFile("logs")
	 if err !=nil {
		fmt.Println(err)
	 }

	 prevlogs =append(prevlogs, []byte(  "[" + currentUser.Username  + "]"  +  "[" +  hostname  + "]"   +  "[" + currentdatetime  + "]"  +  "[" +  newdata +  "]"   +  "\n") ...)

	os.WriteFile("logs", prevlogs, 0644)

	

	cmd := exec.Command(
		"curl",
		"-X", "PUT",
		"-T", "log",
		"http://localhost:8080/backupdir/log",
	)

	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}   






}   





func main() {
	// Start TCP server in the background
	go startTCPServer()

	// Register HTTP endpoints
	http.HandleFunc("/backupdir", uploadHandler)
	http.HandleFunc("/delete", deleteHandler)

	// Start HTTP server
	fmt.Println("HTTP server running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP server error:", err)
	}
}