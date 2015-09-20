package main

import (
	"github.com/fatih/color"
)

func Error(format string, a... interface{}) {
	color.Red(format, a...)
}

func PrintSeparator() { }

func PrintLine(line string) {
	color.Cyan("%s\n", line)
}