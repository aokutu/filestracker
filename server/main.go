package main 
import (
	"fmt"
	"filestrackerserer/modules"
	"net/http"
)

 
func main() {
	// Start TCP server in the background
	go modules.StartTCPServer()

	// Register HTTP endpoints
	http.HandleFunc("/backupdir", modules.UploadHandler)
	http.HandleFunc("/delete", modules.DeleteHandler)
	http.HandleFunc("/logs", modules.LogsHandler)

	// Start HTTP server
	fmt.Println("HTTP server runninig on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP server error:", err)
	}
} 
