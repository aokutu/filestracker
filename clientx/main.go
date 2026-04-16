package main

import (
	"fmt"
	"clientfilestracker/modules"
	"log"
	"net"
	"os"
	"time"
)

func startTCP() {
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

		modules.Writelogs("ANDREW")

		conn.Close()
	}
}

func main() {

	fmt.Println("Starting server...")

	// 🔥 run TCP server in background
	go startTCP()

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