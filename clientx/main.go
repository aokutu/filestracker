package main

import (
	"fmt"
	"clientfilestracker/modules"
	"log"
	_"net"
	"os"
	"time"
	"bufio"
)

func main() {

	fmt.Print("ENTER  SERVER IP ADDRESS:PORTNUMBER :")
	Reader := bufio.NewReader(os.Stdin)
	IpAddress,_ := Reader.ReadString('\n')
	
 err := os.WriteFile("serveraddress",  []byte(IpAddress), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	modules.ListFiles("storage")

	fmt.Println("Starting server...")

	// 🔥 run TCP server in background

	// ensure storage exists
	root := "./storage"
	if err := os.MkdirAll(root, 0755); err != nil {
		log.Fatal(err)
	}

	// initial scan
	modules.ScanDir(root)

	fmt.Println("Polling ./storage folder...")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		modules.ScanDir(root)
	}


}