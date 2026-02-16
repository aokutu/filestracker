// fileops.go - Independent file operations (upload/delete)
package main

import (
    "fmt"

    "net/http"
    
)

const uploadURL = "http://127.0.0.1:8080/backupdir"

func DeleteFile(filename string) {
    deleteURL := uploadURL + "/" + filename
    
    req, err := http.NewRequest("DELETE", deleteURL, nil)
    if err != nil {
        fmt.Println("Error creating DELETE request:", err)
        return
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error deleting file:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == 200 || resp.StatusCode == 204 {
        fmt.Println("File deleted successfully:", filename)
    } else {
        fmt.Println("Delete failed for", filename, "- Status:", resp.Status)
    }
}
