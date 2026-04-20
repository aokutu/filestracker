package modules

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"
)

func ListFiles(Path string) {
	FilePath := []string{}

	filepath.Walk(Path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error:", err)
			return nil
		}

		if info.IsDir() {
			//fmt.Println("DIR :", p)
				FilePath = strings.Split(p, "/")
				fmt.Println(FilePath)
		} else {
			//fmt.Println("FILE:", p)
			FilePath = append(FilePath,p)
			fmt.Println(p)
		}

		return nil
	})
}
