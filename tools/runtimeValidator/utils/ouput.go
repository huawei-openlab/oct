package utils

import (
	"os"
)

func StringOutput(fileName string, result string) error {
	filePath := "./report/" + fileName
	// filePath := "./tools/runtimeValidator/report/" + fileName
	fout, err := os.Create(filePath)
	defer fout.Close()

	if err != nil {
		return err
	} else {
		fout.WriteString(result)
		return nil
	}
}

func SpecifyOutput(filePath string, result string) error {
	fileName := filePath + "linuxspec.json"
	fout, err := os.Create(fileName)
	defer fout.Close()
	if err != nil {
		return err
	} else {
		fout.WriteString(result)
		return nil
	}
}
