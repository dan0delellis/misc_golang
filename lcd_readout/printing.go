package main

import (
	"fmt"
	"regexp"
	"time"
)

var winderp = []byte{'\r', '\n'}
var twospace = []byte{' ', ' '}
var singleLineBreak *regexp.Regexp
var singleLinePattern = `(\r\n?|\n)`

// Printf writes the content to the display, formatted with the provided string and the args Printf expects
func (d Display) Printf(s string, a ...any) (n int, err error) {
	return d.Print(fmt.Sprintf(s, a...))
}

// Write writes the provided string to the display
func (d Display) Print(s string) (n int, err error) {
	b := []byte(s)
	if newlineToCarriage {
		b = singleLineBreak.ReplaceAll(b, winderp)
	}
	return d.Write(b)
}

// Marquee will scroll the given text across the top line line. This requires that the width is set correctly or it will look strange
// Single line breaks, be them \r, \n, or \rn, are replaced with two spaces
// The cursor will always be one after the last character of the marquee. Right now that means position (1,2)
//
// Future will hopefully have the ability to specify a row and offset so the row can have a header
func (d Display) Marquee(s string) {
	// Wait, that's illegal!
	var fmtWidthStr = fmt.Sprint("% " + fmt.Sprintf("%d", cols) + "s")

	d.Clear()

	//buffer originally held the text that was printed on each move, but i instead have the cursor just write each character as it's generated
	//Still use it because I still need to pad the displayed string with that much whitespace
	buffer := []byte(fmt.Sprintf(fmtWidthStr, " "))

	//Replace newlines with two spaces, add two spaces to the end
	b := append(buffer, append(singleLineBreak.ReplaceAll([]byte(s), twospace), twospace...)...)

	for i := 0; i <= len(b); i++ {
		d.Home()

		for j := 0; j < len(buffer); j++ {
			//The cursor moves by one with each print
			d.Print(string(b[(j+i)%len(b)]))
		}

		time.Sleep(250 * time.Millisecond)
	}
}
