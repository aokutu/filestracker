package modules

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
)

func DeleteFile(path string) {
	// Normalize slashes and remove "./" prefix
	relPath := filepath.ToSlash(path)
	if len(relPath) > 0 && relPath[:2] == "./" {
		relPath = relPath[2:]
	}

	// URL-encode to handle spaces/special characters
	encodedPath := url.QueryEscape(relPath)

	serverurl := Getaddress() + "?file="

	// Build curl DELETE command
	//cmd := exec.Command("curl", "-X", "DELETE", serverurl+encodedPath)
    

    cmd := exec.Command(
	"curl",           // executable
	"-X",             // first argument
	"DELETE",         // second argument
	serverurl+encodedPath, // third argument (full URL)
) 
	
fmt.Println("File :",encodedPath)
	 cmd = exec.Command("curl","-X", "DELETE","http://localhost:8080/delete?file=" + encodedPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to delete %s: %v\n", path, err)
		fmt.Printf("Curl output: %s\n", string(output))
	} else {
		fmt.Printf("Successfully deleted %s from backupdir\n", path)
	}
}