package utilities

import "fmt"

var colors = map[string]string{
	"blue":    "\x1b[34;1m",
	"cyan":    "\x1b[36;1m",
	"end":     "\x1b[0m",
	"magenta": "\x1b[35;1m",
	"red":     "\x1b[31;1m",
	"white":   "\x1b[37;1m",
	"yellow":  "\x1b[33;1m",
}

// ColorPrintln is a simple method that outputs colored text to StdOut
func ColorPrintln(text string, color string) {
	fmt.Printf("%s%s%s\n", colors[color], text, colors["end"])
}
