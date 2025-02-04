package main

import (
    "fmt"
    "time"
)
const wait = 3
func main() {
    lcd, err := initLCD("/dev/ttyACM0")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer lcd.Close()
    lcd.Printf("hello!\ndo you like cheese")
    time.Sleep(wait*time.Second)
    lcd.Clear()
    lcd.SetBG(0,128,32)
    lcd.Print("now I am green")
    time.Sleep(wait*time.Second)
    lcd.Clear()
    lcd.SetBG(32,32,128)
    lcd.Print("good bye")
    time.Sleep(wait*time.Second)
    lcd.Clear()
    lcd.SetBG(0,0,0)
}
// Marqee will scroll the given text across one line.
func (d Display) Marqee(s string) {

}
