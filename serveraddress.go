package main 

import "os"
import "strings"
import "fmt"


func Getaddress() string {
    serverbyte, err := os.ReadFile("serveraddress")
    if err != nil {
        fmt.Println("Error reading serveraddress file, using default")
        return "http://127.0.0.1:8080/backupdir"  // Make sure it has the full path
    }
    address := strings.TrimSpace(string(serverbyte))
    
    // Make sure it has http:// prefix
    if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
        address = "http://" + address
    }
    
    // Make sure it ends with /backupdir
    if !strings.HasSuffix(address, "/backupdir") {
        address = address + "/backupdir"
    }
    
    fmt.Printf("DEBUG - Using server address: %s\n", address)
    return address
}