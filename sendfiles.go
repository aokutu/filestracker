package main

import (
    "fmt"
    "os/exec"
    "path/filepath"
)




func Sendfiles(filename string) {
    uploadURL := Getaddress()
    
    // DEBUG: Print the URL being used
    fmt.Printf("DEBUG - Upload URL: %s\n", uploadURL)
    
    fullPath := filepath.Join("storage", filename)
    fmt.Printf("DEBUG - Full file path: %s\n", fullPath)
    
    // Build the curl command as a string for debugging
    cmdStr := fmt.Sprintf("curl -F file=@%s %s", fullPath, uploadURL)
    fmt.Printf("DEBUG - Curl command: %s\n", cmdStr)
    
    cmd := exec.Command("curl", "-F", "file=@"+fullPath, uploadURL)
    output, err := cmd.CombinedOutput()
    
    fmt.Printf("DEBUG - Curl output: %s\n", string(output))
    
    if err != nil {
        fmt.Printf("Error uploading %s: %v\n", filename, err)
    } else {
        fmt.Printf("Successfully uploaded %s\n", filename)
    }
}