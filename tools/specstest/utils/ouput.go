package utils

import (
	"os"
)

func StringOutput(fileName string, result string) error {
	filePath := "./report/" + fileName
	// filePath := "./tools/specstest/report/" + fileName
	fout, err := os.Create(filePath)
	defer fout.Close()

	if err != nil {
		return err
	} else {
		fout.WriteString(result)
		return nil
	}
}
