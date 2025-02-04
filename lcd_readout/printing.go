package main
import (
    "fmt"
)

// Printf writes the content to the display, formatted with the provided string and the args Printf expects
func (d Display) Printf(s string, a ...any) (n int, err error) {
    return fmt.Fprintf(d, s, a...)
}

// Write writess the provided string to the display
func (d Display) Print(s string) (n int, err error) {
    return fmt.Fprintf(d, s)
}
