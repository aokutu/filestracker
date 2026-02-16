package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	//"path"
)

func uploadFile(filePath, uploadURL string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffer to hold the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create the form file field
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return err
	}

	// Write the file contents into the multipart part
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Close the writer to finalize the multipart data
	writer.Close()

	// Send the POST request with the file
	req, err := http.NewRequest("POST", uploadURL, &requestBody)
	if err != nil {
		return err
	}

	// Set the content type for the multipart data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func  Sendfiles( filePath  string) {


	//filePath := "storage/test" // Replace with your file path
	uploadURL := "http://127.0.0.1:8080/backupdir" // Replace with the URL of the remote server
	

	err := uploadFile(filePath, uploadURL)
	if err != nil {
		fmt.Println("Error uploading file:", err)
	} else {
		fmt.Println("File uploaded successfully!")
	}


}

func Sendx(){
	fmt.Println("xxxx")
}