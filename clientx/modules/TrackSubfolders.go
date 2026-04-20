package modules

import (
	"fmt"
	"path/filepath"
	"os"
)

func ListFiles(Path string) {
	filepath.Walk(Path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error:", err)
			return nil
		}

		if info.IsDir() {
			fmt.Println("DIR :", p)
		} else {
			fmt.Println("FILE:", p)
		}

		return nil
	})
}
