package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var verbose bool

var inputDir = "/mydata/input"
var outputDir = "/mydata/output"

func main() {
	flag.BoolVar(&verbose, "v", true, "verbose output")
	flag.Parse()

	print("Copying input to output")
	inputFiles, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range inputFiles {
		src := filepath.Join(inputDir, fi.Name())
		dst := filepath.Join(outputDir, fi.Name())
		err := copyFile(src, dst)
		if err != nil {
			log.Fatal(err)
		}
	}

	print("Done")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func print(v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Println(v...)
}
