package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
  //  "time"
)

func main() {
    // Connect to the server
    conn, err := net.Dial("tcp", "localhost:3000")
    if err != nil {
        fmt.Println("CONNECTION FAILED:", err)
        os.Exit(1)
    }
    defer conn.Close()  // FIXED: Move defer here!

    // Show logo
    logobyte, _ := os.ReadFile("logo")
    fmt.Println(string(logobyte))


    fmt.Println("[ENTER YOUR NAME]:")
    reader := bufio.NewReader(os.Stdin)
  

     var name string

for {
    input, _ := reader.ReadString('\n')
    name = strings.TrimSpace(input)

    if name == "" {
        fmt.Println("[ENTER YOUR NAME]:")
        continue
    }

    break
}

    // Send name ONCE with newline
    fmt.Printf("Sending name: %s\n", name)
	
    _, err = conn.Write([]byte(name + "\n"))  // FIXED: Single send + \n
    if err != nil {
        fmt.Println("Send failed:", err)
        return
    }

    fmt.Print(name, " : ")

    go Loadchats(conn)

    Chatmeg(name,conn)
    





}


func Loadchats(conn net.Conn) {


        // Read server responses
    for {
        buf := make([]byte, 1024)
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Connection closed:", err)
            return
        }
        fmt.Print(string(buf[:n]))
    }
}


func Chatmeg(name string, conn net.Conn) {

    reader := bufio.NewReader(os.Stdin)
    
    

    for {
    

        message, err := reader.ReadString('\n')
        
        if err != nil {
            fmt.Println("Input error:", err)
            return
        }

        message = strings.TrimSpace(message)

        // Send with newline
        _, err = conn.Write([]byte(message + "\n"))
        if err != nil {
            fmt.Println("Send error:", err)
            return
        }
    }
}




