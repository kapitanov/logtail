package main

import (
	"fmt"
	"strings"
	"regexp"
	color "github.com/daviddengcn/go-colortext"
)

func Error(format string, a... interface{}) {
	defer color.ResetColor()
	
	color.Foreground(color.Red, true)
	fmt.Println(fmt.Sprintf(format, a...))
}

func PrintSeparator() { }

func PrintLine(line string) {
	tokens    := tokenizer.Tokenize(line)
	tokens.Print()
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




type LogTokenType int
const (
	TOKEN_PLAIN    LogTokenType = 1
	TOKEN_DATETIME LogTokenType = 2
	TOKEN_SEPARATOR LogTokenType = 3
	TOKEN_DEBUG LogTokenType = 4
	TOKEN_INFO LogTokenType = 5
	TOKEN_WARN LogTokenType = 6
	TOKEN_ERROR LogTokenType = 7
)

type LogToken struct {
	Text string
	Type LogTokenType
}

type LogTokenizerRule struct {
	Expr *regexp.Regexp
	Type LogTokenType
}

type LogTokenizer struct {
	rules []LogTokenizerRule
}

func NewLogTokenizer() LogTokenizer {
	t := LogTokenizer{}
	t.rules = make([]LogTokenizerRule, 0)
	
	t.rules = append(t.rules, createRule(`[0-9]{4,4}(\.|\-)[0-9]{2,2}(\.|\-)[0-9]{2,2}\s[0-9]{2,2}\:[0-9]{2,2}\:[0-9]{2,2}\.[0-9]+`, TOKEN_DATETIME))
	t.rules = append(t.rules, createRule(`[0-9]{4,4}(\.|\-)[0-9]{2,2}(\.|\-)[0-9]{2,2}\s[0-9]{2,2}\:[0-9]{2,2}\:[0-9]{2,2}`, TOKEN_DATETIME))
	t.rules = append(t.rules, createRule(`[0-9]{2,2}\:[0-9]{2,2}\:[0-9]{2,2}\.[0-9]+`, TOKEN_DATETIME))
	t.rules = append(t.rules, createRule(`[0-9]{2,2}\:[0-9]{2,2}\:[0-9]{2,2}`, TOKEN_DATETIME))
	t.rules = append(t.rules, createRule(`\-\>`, TOKEN_SEPARATOR))
	t.rules = append(t.rules, createRule(` at `, TOKEN_SEPARATOR))
	t.rules = append(t.rules, createRule(`TRACE`, TOKEN_DEBUG))
	t.rules = append(t.rules, createRule(`DEBUG`, TOKEN_DEBUG))
	t.rules = append(t.rules, createRule(`INFO`, TOKEN_INFO))
	t.rules = append(t.rules, createRule(`WARN`, TOKEN_WARN))
	t.rules = append(t.rules, createRule(`ERROR`, TOKEN_ERROR))
	t.rules = append(t.rules, createRule(`FATAL`, TOKEN_ERROR))
	
	return t	
}

var tokenizer LogTokenizer = NewLogTokenizer()

func createRule(pattern string, tokenType LogTokenType) LogTokenizerRule {
	expr, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	return LogTokenizerRule { expr, tokenType }
}

type LogTokens []LogToken

func (t LogTokenizer) Tokenize(text string) LogTokens {
	tokens := make([]LogToken, 0)
	for {
		for _, rule := range t.rules {
			idx := rule.Expr.FindStringSubmatchIndex(text)
			if idx != nil && idx[0] == 0 {
				textPart := text[idx[0]:idx[1]]
				text = text[idx[1] :]			
				tokens = append(tokens, LogToken { textPart, rule.Type })
				continue
			}
		}
					
		if len(text) < 1 {
			break
		}
		
		tokens = append(tokens, LogToken { text[0:1], TOKEN_PLAIN})		
		if len(text) <= 1 {
			break
		}
		text = text[1 : ]		
	}
		
	return LogTokens(tokens)
}

func (tokens LogTokens) Print() {
	defer color.ResetColor()
	
	var defColor  color.Color = color.White
	var defBright bool        = true
	
	for _, token := range tokens {
		if token.Type >= TOKEN_DEBUG {
			defColor, defBright, _ = token.Type.getFgColor()
		}
	}	
	
	for _, token := range tokens {
		fg, bright, specific := token.Type.getFgColor()
		if !specific {
			fg     = defColor 
			bright = defBright
		}
		
		color.Foreground(fg, bright)
		fmt.Print(token.Text)
	}
	
	fmt.Print("\n")
}

func (t LogTokenType) getFgColor() (color.Color, bool, bool) {
	if t == TOKEN_DATETIME {
		return color.Cyan, true, true
	}

	if t == TOKEN_SEPARATOR {
		return color.Cyan, true, true
	}
	
	if t == TOKEN_DEBUG {
		return color.White, false, true
	}
	
	if t == TOKEN_INFO {
		return color.Green, true, true
	}
	
	if t == TOKEN_WARN {
		return color.Yellow, true, true
	}
	
	if t == TOKEN_ERROR {
		return color.Red, true, true
	}
	
	return color.White, true, false
}
