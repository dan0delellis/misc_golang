package main
import (
    "fmt"
)


var winderp = []byte{'\r','\n'}
var twospace = []byte{' ', ' '}

// Printf writes the content to the display, formatted with the provided string and the args Printf expects
func (d Display) Printf(s string, a ...any) (n int, err error) {
    return d.Print(fmt.Sprintf(s,a...))
}

// Write writes the provided string to the display
func (d Display) Print(s string) (n int, err error) {
    b := []byte(s)
    if newlineToCarriage {
        b = singleLineBreak.ReplaceAll(b, winderp)
    }
    return d.Write(b)
}
