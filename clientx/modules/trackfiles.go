package modules

import (
    "fmt"
    "os"
    "os/user"
    "path/filepath"  // Make sure this is imported
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



func ScanDir(root string) {
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
    fmt.Println(userdetails(), Timestamp(), "New file/folder:", path)
    // Extract just the filename from the full path
    filename := filepath.Base(path)  // This turns "storage/abc" into "abc"

    fmt.Println("ANDREW")
    Sendfiles(filename)  // Pass just the filename
} else if !prev.ModTime.Equal(info.ModTime()) {
    op := "modified"
    if info.IsDir() != prev.IsDir {
        op = "replaced"
    }
    fmt.Println(userdetails(), Timestamp(), op+":", path)
    // Also fix the modified case
    filename := filepath.Base(path)
    Sendfiles(filename)
}
        return nil
    })

     // Detect deletions
for path := range fileStates {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        fmt.Println(userdetails(), Timestamp(), "DELETION DETECTED:", path)
        delete(fileStates, path)
        
        // Extract just the filename
        filename := filepath.Base(path)
        
        // Call DeleteFile with just the filename
        DeleteFile(filename)
        Logsupload()
    }
}
}

func Timestamp() string {
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
