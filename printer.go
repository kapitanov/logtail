package main

import (
	"fmt"
	"strings"
	"github.com/ttacon/chalk"
)

func Error(format string, a... interface{}) {
	fmt.Println(chalk.Red.Color(fmt.Sprintf(format, a...)))
}

func PrintSeparator() { }

func PrintLine(line string) {
	color := detectPrimaryColor(line)
	fmt.Println(color.Color(line))
}

func detectPrimaryColor(line string) chalk.Color {
	if strings.Contains(line, "TRACE") {
		return chalk.Blue	
	}

	if strings.Contains(line, "DEBUG") {
		return chalk.Blue	
	}
	
	if strings.Contains(line, "INFO") {
		return chalk.Cyan	
	}
	
	if strings.Contains(line, "WARN") {
		return chalk.Yellow
	}
	
	if strings.Contains(line, "ERROR") {
		return chalk.Red
	}
	
	if strings.Contains(line, "FATAL") {
		return chalk.Red
	}
	
	return chalk.White
}