package main

import (
    "fmt"
    "time"
)
const wait = 3
func main() {
    lcd, err := InitLCD("/dev/ttyACM0")
    if err != nil {
        fmt.Println(err)
        return
    }
    lcd.Clear()
    lcd.BlinkyBlock()
    lcd.Marquee("NOW WHAT")

    defer lcd.Close()
}


// Wait, that's illegal!
var fmtWidthStr = fmt.Sprintf("% " + fmt.Sprintf("%d",cols) + "s")

// Marquee will scroll the given text across one line. This requires that the width is set correctly or it will look strange
// Single line breaks, be them \r, \n, or \rn, are replaced with two spaces
// It will not move the cursor with the text. It will hopefully block all text input while running.
func (d Display) Marquee(s string) {
    s = string(singleLineBreak.ReplaceAll([]byte(s), twospace))

    //This is going to be the row of text we display
    buffer := fmt.Sprintf(fmtWidthStr,"")


    for i:=0; i < len(s); i++ {
        d.Home()
        d.Print(buffer)
        time.Sleep(300 * time.Millisecond)
    }
}

