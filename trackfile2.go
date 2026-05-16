package main
///
import (
    "fmt"
    "log"
    "github.com/fsnotify/fsnotify"
)

func main() {
    // Create new watcher
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    // Start listening for events
    done := make(chan bool)
    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                fmt.Println("Event:", event)
                
                // Check event type
                if event.Op&fsnotify.Write == fsnotify.Write {
                    fmt.Println("Modified file:", event.Name)
                }
                if event.Op&fsnotify.Create == fsnotify.Create {
                    fmt.Println("Created file:", event.Name)
                }
                if event.Op&fsnotify.Remove == fsnotify.Remove {
                    fmt.Println("Deleted file:", event.Name)
                }
                if event.Op&fsnotify.Rename == fsnotify.Rename {
                    fmt.Println("Renamed file:", event.Name)
                }
                if event.Op&fsnotify.Chmod == fsnotify.Chmod {
                    fmt.Println("Permissions changed:", event.Name)
                }

            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("Error:", err)
            }
        }
    }()

    // Add a file or directory to watch
    err = watcher.Add(".")  
    if err != nil {
        log.Fatal(err)
    }

    // Watch multiple paths
    // err = watcher.Add("/another/path")
    
    <-done // Block forever
}