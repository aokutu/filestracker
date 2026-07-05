package modules

import "fmt"
import "os/exec"

func Logsupload(){
 out, _ := exec.Command(
        "curl",
        "-H", "X-API-Key: 4664" ,
        "-X", "PUT",
        "-T", "logs",
        "http://localhost:8080/logs",
    ).CombinedOutput() 

 fmt.Println( string(out))
	
}