package main

import (
    "fmt"
    "log"
	"time"
	"strconv"
	"strings"
    "github.com/fsnotify/fsnotify"
	"os"
	 "os/user"
)

func main() {



	timestamp()
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    // Start listening for events
    go func() {
        for {
            select {
            case event := <-watcher.Events:
               // fmt.Println("Event:", event)
                if event.Op&fsnotify.Create == fsnotify.Create {
                    fmt.Println(userdetails(),timestamp() ,"New file created:", event.Name)
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    fmt.Println(userdetails(),timestamp() ,"File modified:", event.Name)
                }
                if event.Op&fsnotify.Remove == fsnotify.Remove {
                    fmt.Println(userdetails(),timestamp(),"File deleted:", event.Name)
                }
            case err := <-watcher.Errors:
                fmt.Println("Error:", err)
            }
        }
    }()

    // Watch your "storage" folder
    err = watcher.Add("./storage")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Watching ./storage folder…")
    select {} // block forever
}

func timestamp() string{
	 now := time.Now()
	var tracktime strings.Builder
	tracktime.WriteString(strconv.Itoa(now.Year()))
	tracktime.WriteString("-")
	tracktime.WriteString(strconv.Itoa(int(now.Month())))
	tracktime.WriteString("-")
	tracktime.WriteString(strconv.Itoa(now.Day()))
	tracktime.WriteString("-")
	tracktime.WriteString(strconv.Itoa(now.Hour()))
	tracktime.WriteString(":")
	tracktime.WriteString(strconv.Itoa(now.Minute()))
	tracktime.WriteString(":")
	tracktime.WriteString(strconv.Itoa(now.Second()))
	return tracktime.String()


}

func userdetails() string {

	var userdetails  strings.Builder
		hostname, err := os.Hostname()
if err != nil {
    fmt.Println("Error:", err)
}

userdetails.WriteString(hostname)

 u, err := user.Current()
    if err != nil {
        fmt.Println("Error getting user:", err)
       
    }

	userdetails.WriteString(":")
	userdetails.WriteString(u.Username)
	return userdetails.String()

}
