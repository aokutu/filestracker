package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"bufio"
	"strings"
)

func main() {
	fmt.Print("ENTER SERVER ADDRESS :")

	reader := bufio.NewReader(os.Stdin)

	Ipaddress,_ := reader.ReadString('\n')
	 Ipaddress = strings.TrimSpace(Ipaddress)

	// Hit your local endpoint
	resp, err := http.Get("http://" + Ipaddress + ":8080/readlogs")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading body: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(body),"\n")
	
}



