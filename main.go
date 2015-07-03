package main

import (
	"flag"
	"fmt"
)

var verbose bool

func verb(v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Println(v...)
}

func main() {
	flag.BoolVar(&verbose, "v", true, "verbose output")
	flag.Parse()

	verb("Copying input to output")

	verb("Done")
}
