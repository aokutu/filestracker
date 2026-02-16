package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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

func main() {
	// Define the endpoint for file upload
	http.HandleFunc("/backupdir", uploadHandler)

	// Start the server
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
