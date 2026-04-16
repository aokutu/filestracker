package modules

import (
	"fmt"
	"os"
	"time"
	"os/user"
	
)

func Writelogs(newdata string){

	 currentUser, err := user.Current()
    if err != nil {
        fmt.Println("Error getting current user:", err)
        return
    } 

    // Get computer/hostname
    hostname, err := os.Hostname()
    if err != nil {
        fmt.Println("Error getting hostname:", err)
        return
    }

	currenttime := time.Now()
    currentdatetime := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]",
        currenttime.Year(), currenttime.Month(), currenttime.Day(),
        currenttime.Hour(), currenttime.Minute(), currenttime.Second())


	 var prevlogs []byte

	 prevlogs =append(prevlogs, []byte(  "[" + currentUser.Username  + "]"  +  "[" +  hostname  + "]"   +  "[" + currentdatetime  + "]"  +  "[" +  newdata +  "]"   +  "\n") ...)

	os.WriteFile("logs", prevlogs, 0644)

	 fmt.Print(prevlogs)

	 

}