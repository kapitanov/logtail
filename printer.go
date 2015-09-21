package main

import (
	"fmt"
	"strings"
	color "github.com/daviddengcn/go-colortext"
)

func Error(format string, a... interface{}) {
	defer color.ResetColor()
	
	color.Foreground(color.Red, true)
	fmt.Println(fmt.Sprintf(format, a...))
}

func PrintSeparator() { }

func PrintLine(line string) {
	defer color.ResetColor()
	color.Foreground(detectPrimaryColor(line), true)
	fmt.Println(line)
}

func detectPrimaryColor(line string) color.Color {
	if strings.Contains(line, "TRACE") {
		return color.Blue	
	}

	if strings.Contains(line, "DEBUG") {
		return color.Blue	
	}
	
	if strings.Contains(line, "INFO") {
		return color.Cyan	
	}
	
	if strings.Contains(line, "WARN") {
		return color.Yellow
	}
	
	if strings.Contains(line, "ERROR") {
		return color.Red
	}
	
	if strings.Contains(line, "FATAL") {
		return color.Red
	}
	
	return color.White
}