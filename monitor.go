package main

import (
	"os"
	"fmt"
	"io"
	"bufio"
	"path/filepath"
	"time"
)

const WATCH_TIMEOUT time.Duration = time.Duration(1*1000)

type LogFile struct {
	Cursor int64
	FileName string
}

func NewLogFile(filename string) (logFile *LogFile, err error) {
	logFile = new(LogFile)
	logFile.FileName = filepath.Clean(filename)
			
	return logFile, nil
}

func (logFile *LogFile) Close() { }

func (logFile *LogFile) Monitor(maxLinesToPrint int) error {
	logFile.initialize(maxLinesToPrint)	
	go logFile.monitorBg()	
	return nil
}

func (logFile *LogFile) monitorBg() {
	for {		
		logFile.onModified()
		time.Sleep(WATCH_TIMEOUT)
	}
}

func (logFile *LogFile) open() (file *os.File, err error) {
	file, err = os.OpenFile(logFile.FileName, os.O_RDONLY, 0)	
	return
}

func (logFile *LogFile) initialize(maxLinesToPrint int) error {
	file, err := logFile.open()
	if err == os.ErrNotExist {
		// File doesn't exists
		logFile.Cursor = 0
		return nil
	}
	
	if err != nil {	
		return err
	}
	
	defer file.Close()
		
	stat, err := file.Stat()
	if err != nil {
		return err
	}
						
	size := stat.Size()					
	
	from := int64(0)
	if maxLinesToPrint > 0 {		
		from = findTail(file, maxLinesToPrint)
	}
	printTail(file, from)
	
	logFile.Cursor = size	
	return nil
}

func printTail(file *os.File, from int64) (n int64) {	
	file.Seek(from, os.SEEK_SET)		
	reader := bufio.NewReader(file)		
	
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "ReadAt() failed! %s\n", err)
			return 0
		}	
		
		n += int64(len(line))
		text := string(line)
		PrintLine(text)
	}
}

func findTail(file *os.File, n int) int64 {
	file.Seek(0, os.SEEK_SET)		
	reader := bufio.NewReader(file)		
	
	offset  := int64(0)
	offsets := make(map[int]int64, n)	
	lineNum := 0
	
	for {		
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		
		if err != nil {
			return 0
		}	
		
		offsets[lineNum] = offset
		
		lineNum++
		offset += int64(len(line)) + 2
	}
	
	if len(offsets) <= n {
		
		return 0
	}
	
	return offsets[len(offsets) - n]
}

func (logFile *LogFile) onModified() {	
	file, err := logFile.open()
	if err == os.ErrNotExist {
		logFile.Cursor = 0
		return
	}
	
	if err != nil {	
		return
	}
	
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}
						
	newCursor := stat.Size()
	delta := newCursor - logFile.Cursor	
	if delta == 0 {
		return
	}
	
	if delta < 0 {
		PrintSeparator()
		logFile.Cursor = 0	
	}
				
	printTail(file, logFile.Cursor)
	logFile.Cursor = newCursor
}

