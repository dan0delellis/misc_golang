package main

import (
    "fmt"
    _"time"
)

func main() {
    lcd, err := initLCD("/dev/ttyACM0")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer lcd.Close()
    lcd.Printf("% -5s %x %f", "fun", 145, 1.34)
}
// Writef writes the content to the display, formatted with the provided string and the args Printf expects
func (d Display) Printf(s string, a ...any) (n int, err error) {
    return fmt.Fprintf(d, s, a...)
}

// Write writess the provided string to the display
func (d Display) Print(s string) (n int, err error) {
    return fmt.Fprintf(d, s)
}

// marqee will scroll the given text across one line.
func (d Display) Marqee(s string) {

}
