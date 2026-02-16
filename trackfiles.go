package main

import (
    "fmt"
    "log"
    "os"
    "os/user"
    "path/filepath"
    "sync"
    "time"
)

type FileInfo struct {
    ModTime time.Time
    IsDir   bool
}

var (
    fileStates = make(map[string]FileInfo)
    mu         sync.RWMutex
)

func main() {
	
    timestamp()
    root := "./storage"
    if err := os.MkdirAll(root, 0755); err != nil {
        log.Fatal(err)
    }
    scanDir(root)

    fmt.Println("Polling ./storage folder...")
    ticker := time.NewTicker(2 * time.Second) // Adjust interval as needed
    defer ticker.Stop()

    for range ticker.C {
        scanDir(root)
    }
}

func scanDir(root string) {
    mu.Lock()
    defer mu.Unlock()

    filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        info, err := d.Info()
        if err != nil {
            return err
        }

        prev, exists := fileStates[path]
        fileStates[path] = FileInfo{ModTime: info.ModTime(), IsDir: info.IsDir()}

        if !exists {
            fmt.Println(userdetails(), timestamp(), "New file/folder:", )

			Sendfiles( path) 
        } else if !prev.ModTime.Equal(info.ModTime()) {
            op := "modified"
            if info.IsDir() != prev.IsDir {
                op = "replaced"
            }
            fmt.Println(userdetails(), timestamp(), op+":", path)
			Sendfiles(path) 
        }
        return nil
    })

    // Detect deletions
    for path := range fileStates {
        _, err := os.Stat(path)
        if os.IsNotExist(err) {
            fmt.Println(userdetails(), timestamp(), "deleted:", path)
            delete(fileStates, path)
			DeleteFile(path)
        }
    }
}

func timestamp() string {
    now := time.Now()
    return fmt.Sprintf("%d-%02d-%02d-%02d:%02d:%02d",
        now.Year(), now.Month(), now.Day(),
        now.Hour(), now.Minute(), now.Second())
}

func userdetails() string {
    hostname, _ := os.Hostname()
    u, _ := user.Current()
    return fmt.Sprintf("%s:%s", hostname, u.Username)
}
