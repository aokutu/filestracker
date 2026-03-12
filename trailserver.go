package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	fmt.Println("Starting server...")

	listener, err := net.Listen("tcp", "localhost:3001")
	if err != nil {
		fmt.Println("SERVER NOT STARTED")
		os.Exit(1)
	}

	fmt.Println("Server listening on port 3000")

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

	currenttime := time.Now()
    currentdatetime := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]",
        currenttime.Year(), currenttime.Month(), currenttime.Day(),
        currenttime.Hour(), currenttime.Minute(), currenttime.Second())


	 prevlogs, err := os.ReadFile("logs")
	 if err !=nil {
		fmt.Println(err)
	 }

	 prevlogs =append(prevlogs, []byte(currentdatetime + newdata + "\n") ...)

	os.WriteFile("logs", prevlogs, 0644)

	 fmt.Print(prevlogs)

}