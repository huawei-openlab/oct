package hostenv

import (
	// "archive/tar"
	"fmt"
	// "io"
	"log"
	"os"
	"os/exec"
	// "path"
)

func CreateBoundle() error {
	var cmd *exec.Cmd
	var err error
	fmt.Println("Starting to create boundle...")
	fmt.Println("	Going to create bind path for hostend-containerend communicating...")
	goPath := os.Getenv("GOPATH")
	bindPath := goPath + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/containerend/"
	err = os.MkdirAll(bindPath, 0777)
	if err != nil {
		log.Fatalf("CreateBoundle create bindpath %v errr %v", bindPath, err)
	}
	fmt.Println("	done")
	fmt.Println("	Creating root filesystem for boundle, maybe need several mins...")
	cmd = exec.Command("/bin/sh", "-c", "docker pull ubuntu:14.04")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("CreateBoundle pull image error, %v", err)
	}
	cmd = exec.Command("/bin/sh", "-c", "docker export $(docker create ubuntu) > ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("CreateBoundle export filesystem err, %v", err)
	}

	destDir := "./rootfs"
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		log.Fatalf("CreateBoundle create rootfs dir err, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "tar -C "+destDir+" -xf ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("CreateBoundle extarct rootfs err, %v", err)
	}
	/*destDir := "./rootfs"
	srcDir := "ubuntu.tar"
	untarPkg(destDir, srcDir)*/

	fmt.Println("Create boundle done")
	return nil
}

/*func untarPkg(destDir string, srcDir string) {
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		log.Fatalf("CreateBoundle create rootfs dir err, %v", err)
	}
	fp, err := os.Open(srcDir)
	if err != nil {
		log.Fatalf("CreateBoundle open ubuntu.tar err, %v", err)
	}
	tr := tar.NewReader(fp)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		} else if err != nil {
			panic(err)
		}
		//handleError(err)
		// Check if it is diretory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create diretory before create file
			os.MkdirAll(destDir+"/"+path.Dir(hdr.Name), os.ModePerm)
			// Write data to file
			fw, _ := os.Create(destDir + "/" + hdr.Name)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				panic(err)
			}
		}
	}
}*/
