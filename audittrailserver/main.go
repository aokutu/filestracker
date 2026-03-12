package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "time"
    "sync"
)


  var lk sync.Mutex

func main() {
    // Try to read chat.txt once at startup (for initial display)

  

    lk.Lock()
    rawdata, err := os.ReadFile("chat.txt")
    lk.Unlock()
    if err != nil {
        fmt.Println("TRACK  FILE NOT CREATED")
        os.Exit(0)
    }

    listen, err := net.Listen("tcp", "localhost:3000")
    if err != nil {
        fmt.Println("AUDIT TRACK  SERVER NOT STARTED")
        os.Exit(0)
    }

    fmt.Println(string(rawdata)) // Initial chat display

}

// ----- Corrected handleClient -----
func handleClient(c net.Conn) {
    defer c.Close()

    reader := bufio.NewReader(c)

    // Read client name once
    name, err := reader.ReadString('\n')
    if err != nil {
        log.Println("Failed to read name from client:", c.RemoteAddr(), err)
        return
    }
    name = strings.TrimSpace(name)
    log.Printf("Client connected with name: %s", name)

    // Append greeting to chat.txt
    Updateuser(name)

    // ---- Broadcast loop for this client ----
    go func( ) {
        lastLen := 0
        for {
             lk.Lock()
            rawdata, err := os.ReadFile("chat.txt")
            
            if err != nil {
                log.Println("Failed to read chat.txt:", err)
                return
            }

            if len(rawdata) > lastLen {
                _, err := c.Write(rawdata[lastLen:])
                if err != nil {
                    log.Println("Client disconnected:", c.RemoteAddr(), err)
                    return
                }
                lastLen = len(rawdata)
            }
             lk.Unlock()
            time.Sleep(1 * time.Second)
        }
    }()

    // ---- Loop to read messages from this client ----
    for {
        
        message, err := reader.ReadString('\n') // BLOCKS until client sends a message
        if err != nil {
            log.Println("Client disconnected:", c.RemoteAddr(), err)
            return
        }

        message = strings.TrimSpace(message)
        if message == "" {
            continue // ignore empty messages
        }

        ProcessMessage(name, message,&lk) // append actual user message to chat.txt
    }
}

// ----- Updateuser function -----
func Updateuser(name string) {
    currenttime := time.Now()
    currentdatetime := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]",
        currenttime.Year(), currenttime.Month(), currenttime.Day(),
        currenttime.Hour(), currenttime.Minute(), currenttime.Second())

    name = "[" + name + "]"
    newdata := []byte(currentdatetime + " " + name + " : hello\n")


     lk.Lock()
     defer lk.Unlock()
    f, err := os.OpenFile("chat.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println("Failed to open chat.txt for append:", err)
        return
    }
    defer f.Close()
    
    
    if _, err := f.Write(newdata); err != nil {
        log.Println("Failed to write to chat.txt:", err)
    }
    
}

// ----- ProcessMessage function -----
func ProcessMessage(name string, message string,lk *sync.Mutex) {
    currenttime := time.Now()
    currentdatetime := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]",
        currenttime.Year(), currenttime.Month(), currenttime.Day(),
        currenttime.Hour(), currenttime.Minute(), currenttime.Second())

    name = "[" + name + "]"
    newdata := []byte(currentdatetime + " " + name + " : " + message + "\n")
     lk.Lock()
     defer lk.Unlock()
    f, err := os.OpenFile("chat.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println("Failed to open chat.txt for append:", err)
        return
    }
    defer f.Close()
    if _, err := f.Write(newdata); err != nil {
        log.Println("Failed to write to chat.txt:", err)
    }

    
}