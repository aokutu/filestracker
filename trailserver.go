package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"os/user"
	
)

func main() {

	fmt.Println("Starting server...")

	listener, err := net.Listen("tcp", "localhost:3002")
	if err != nil {
		fmt.Println("SERVER NOT STARTED")
		os.Exit(1)
	}

	fmt.Println("Server listening on port 3002")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		fmt.Println("Client connected:", conn.RemoteAddr())

		Writelogs("ANDREW")

		conn.Close() // close connection for now
	}
}


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


	 prevlogs, err := os.ReadFile("logs")
	 if err !=nil {
		fmt.Println(err)
	 }

	 prevlogs =append(prevlogs, []byte(  "[" + currentUser.Username  + "]"  +  "[" +  hostname  + "]"   +  "[" + currentdatetime  + "]"  +  "[" +  newdata +  "]"   +  "\n") ...)

	os.WriteFile("logs", prevlogs, 0644)

	 fmt.Print(prevlogs)

}