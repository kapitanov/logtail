package main

import (
	"os"
	"fmt"
	"flag"
)

var (
	file *LogFile
	linesToPrint = flag.Int("tail", -1, "Tail to print at start")
)

func main() {
	flag.Parse()
	if !flag.Parsed() {
		fmt.Fprintf(os.Stderr, "Invalid command line!\n")
		flag.PrintDefaults()
		return
	}
	
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "No file is specified\n")
		flag.PrintDefaults()
		return
	}
		
	filename := flag.Arg(0)
	
	file, err := NewLogFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewLogFile() failed! %s\n", err)
		panic(err)
	}
	defer file.Close()
		
	err = file.Monitor(*linesToPrint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Monitor() failed! %s\n", err)
		panic(err)
	}
		
	shutdown := make(chan bool)
	<- shutdown
}