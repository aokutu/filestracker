package main 
import (
	"fmt"
	"filestrackerserer/modules"
	"net/http"
	"log"
	"os"
	"bufio"
)

 
func main() {
	// Start TCP server in the background


	// Register HTTP endpoints
	http.Handle("/backupdir",KeyloggerMiddleware(http.HandlerFunc(modules.UploadHandler )))
	http.Handle("/delete", KeyloggerMiddleware(http.HandlerFunc(modules.DeleteHandler )))
	http.Handle("/logs",   KeyloggerMiddleware(http.HandlerFunc(modules.LogsHandler)) )
	http.Handle("/readlogs", KeyloggerMiddleware(http.HandlerFunc(Readlogs )))
	http.Handle("/apikey", KeyloggerMiddleware(http.HandlerFunc(Keylogger)))


	// Start HTTP server
	fmt.Println("HTTP server runninig on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP server error:", err)
	}
} 


//curl -H "X-API-Key: 4664" http://localhost:8080/apikey
func KeyloggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "4664" {
			log.Println("INVALID KEY ")
			return 
		} else {
			log.Println("ACCESS  PERMITED")
			// Pass to the next handler
			next.ServeHTTP(w, r)
		}
		
		
		
	})
}


func Keylogger(w http.ResponseWriter, r *http.Request) {
}

func Readlogs(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")
	log.Println(ua)
	if ua != "Go-http-client/1.1" {
			http.Error(w, "INVALID  ACCESS TOOL ", http.StatusForbidden)
		return 
	}
	Filestring(w)  // pass w to the function
}

 func Filestring(w http.ResponseWriter) {
	file, err := os.Open("logs")
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Print from last line to first
	for i := len(lines) - 1; i >= 0; i-- {
		w.Write([]byte(lines[i] + "\n"))
	}
} 