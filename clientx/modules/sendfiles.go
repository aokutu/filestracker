package modules

import (
    "fmt"
    "os/exec"
    "path/filepath"
    "net"
    "log"
    "os"
)




// Sendfiles uploads file + creates necessary folders on destination
func Sendfiles(relativePath string) {  // e.g. "A/B/C/photo.jpg"
    uploadURL := Getaddress()

    fullPath := filepath.Join("storage", relativePath)

    fmt.Printf("DEBUG - Upload URL: %s\n", uploadURL)
    fmt.Printf("DEBUG - Relative path: %s\n", relativePath)
    fmt.Printf("DEBUG - Full local path: %s\n", fullPath)

    // Check if file exists
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        fmt.Printf("ERROR: File not found: %s\n", fullPath)
        return
    }

    // Use curl with two form fields:
    // - filepath = folder structure + filename
    // - file     = actual file content
    cmd := exec.Command("curl",
        "-F", "filepath="+relativePath,           // Tell server the desired path
        "-F", "file=@"+fullPath,                  // The file itself
        uploadURL,
    )

    fmt.Printf("DEBUG - Curl command: curl -F filepath=%s -F file=@%s %s\n", relativePath, fullPath, uploadURL)

    output, err := cmd.CombinedOutput()
    fmt.Printf("DEBUG - Curl output:\n%s\n", string(output))

    if err != nil {
        fmt.Printf("Error uploading %s: %v\n", relativePath, err)
    } else {
        fmt.Printf("✅ Successfully uploaded %s (with folder structure)\n", relativePath)
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