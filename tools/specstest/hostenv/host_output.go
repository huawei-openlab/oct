package hostenv

import (
	"log"
	"os"
)

func HostOutput(outputFile string, outPutString string) {
	outputPath := "/tmp/testtool/" + outputFile + ".json"
	fout, err := os.Create(outputPath)
	defer fout.Close()

	if err != nil {
		log.Fatal("Pkg hostsetup.HostOutput err = %v", err)
	} else {
		fout.WriteString(outPutString)
	}
}
