package main

import (
	"fmt"
	"clientfilestracker/modules"
	"log"
	_"net"
	"os"
	"time"
)

func main() {

	modules.ListFiles(".")

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