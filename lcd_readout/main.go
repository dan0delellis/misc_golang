package main

import (
    "fmt"
    "time"
    "golang.org/x/image/colornames"
)
const wait = 3
func main() {
    lcd, err := InitLCD("/dev/ttyACM0")
    defer lcd.Close()
    if err != nil {
        fmt.Println(err)
        return
    }

    i := 0
    for k,v := range colornames.Map {
        i++
        lcd.Clear()
        lcd.Print(k)
        lcd.SetBG(v.R, v.G, v.B)
        lcd.Home()
        time.Sleep(700 * time.Millisecond)
        if i > 10 {
            break
        }
    }

    if err != nil {
        fmt.Println(err)
    }

    for _, v := range BrightnessNames {
        lcd.Clear()
        lcd.Print(v)
        lcd.BrightnessKeyword(v)
        lcd.Home()
        time.Sleep(500 * time.Millisecond)
    }

    lcd.Marquee(`Lorem ipsum dolor sit amet`)
}
