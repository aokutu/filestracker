package main

import "fmt"
import "os/exec"

func Logsupload(){
 out, _ := exec.Command(
        "curl",
        "-X", "PUT",
        "-T", "logs",
        "http://localhost:8080/logs",
    ).CombinedOutput() 

 fmt.Println( string(out))
	
}