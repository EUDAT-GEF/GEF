package main

import (
	"os"
	"log"
	"io/ioutil"
	"fmt"
	"io"
)

const (
	sourceFolder = "/root/buffer"
	targetFolder = "/root/volume"
)

// copySingleFile copies a file
func copySingleFile(sourcePath string, destPath string) (err error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}

	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		sourceInfo, err := os.Stat(sourcePath)
		if err != nil {
			err = os.Chmod(destPath, sourceInfo.Mode())
		}

	}

	return
}

// copyDirectory recursively copies a folder
func copyDirectory(sourcePath string, destPath string) (err error) {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destPath, sourceInfo.Mode())
	if err != nil {
		return err
	}

	directory, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		sourceFileHandler := sourcePath + "/" + obj.Name()
		destinationFileHandler := destPath + "/" + obj.Name()

		if obj.IsDir() {
			err = copyDirectory(sourceFileHandler, destinationFileHandler)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = copySingleFile(sourceFileHandler, destinationFileHandler)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

// copyFromSource copies files and folders from a given location to another specified path
func copyFromSource(sourcePath string, destPath string) (error) {
	log.Println("Reading folder: " + sourcePath)
	doesExist, hasErrors := exists(sourcePath)
	if hasErrors == nil {
		if doesExist {
			files, _ := ioutil.ReadDir(sourcePath)
			for _, f := range files {
				if f.IsDir() == true {
					hasErrors = copyDirectory(sourcePath + "/" + f.Name(), destPath)
				}
				if hasErrors == nil {
					copySingleFile(sourcePath + "/" + f.Name(), destPath)
					log.Println("Copying a file: " + f.Name())
				}
			}
		}
	}

	return hasErrors
}

// exists checks if a given path exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}


func main() {
	log.Println("Source folder :" + sourceFolder)
	log.Println("Targer folder :" + targetFolder)

	err := copyFromSource(sourceFolder, targetFolder)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Directory has been copied")
	}
}
