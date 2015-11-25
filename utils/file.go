package utils

import (
	"fmt"
	"os"
)

func FileExist(filename string) bool {
	stat, err := os.Stat(filename)
	return err == nil || os.IsExist(err) || (!stat.IsDir())
}

func DirExist(path string) bool {
	stat, err := os.Stat(path)
	return err == nil || os.IsExist(err) || stat.IsDir()
}

func RemoveFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	} else {
		return fmt.Errorf("File/Dir %v NOT exist, cannot be removed", path)
	}
	return nil
}
