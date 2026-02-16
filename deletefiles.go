package main

import (
    "fmt"
    "os/exec"
)

func DeleteFile(filename string) {
    serverurl :=Getaddress() + "?file=" 
    cmd := exec.Command("curl", "-X", "DELETE",serverurl+filename)
    
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Failed to delete %s: %v\n", filename, err)
    } else {
        fmt.Printf("Successfully deleted %s from backupdir\n", filename)
    }
}