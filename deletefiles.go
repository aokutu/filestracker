package main

import (
    "fmt"
    "os/exec"
)


 

func DeleteFile(filename string) {
    serverurl := Getaddress() + "?file="

    cmd := exec.Command("curl", "-X", "DELETE", serverurl+filename)

    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Failed to delete %s: %v\n", filename, err)
        fmt.Printf("Curl output: %s\n", string(output)) // 👈 shows real error
    } else {
        fmt.Printf("Successfully deleted %s from backupdir\n", filename)
    }
}  