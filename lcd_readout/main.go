package main

import (
	"fmt"
	"golang.org/x/image/colornames"
	"time"
)

const wait = 3

func main() {
	lcd, err := InitLCD("/dev/ttyACM0")
	defer lcd.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	lcd.demo()
}

func (lcd Display) demo() {
	lcd.Print("Here")
	time.Sleep(300 * time.Millisecond)
	lcd.Print(" we")
	time.Sleep(300 * time.Millisecond)
	lcd.Print(" GO!\n")
	time.Sleep(300 * time.Millisecond)

	lcd.ColorKeyword("green")
	lcd.Print("color test")
	time.Sleep(time.Second)
	i := 0
	for k, v := range colornames.Map {
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
	lcd.Clear()
	lcd.ColorKeyword("gold")

	lcd.Print("brightness test")
	time.Sleep(time.Second)

	for _, v := range BrightnessNames {
		lcd.Clear()
		lcd.Print(v)
		lcd.BrightnessKeyword(v)
		lcd.Home()
		time.Sleep(500 * time.Millisecond)
	}

	lcd.BrightnessKeyword("dim")
	lcd.Marquee(`Lorem ipsum dolor sit amet`)
	lcd.Home()
	lcd.BrightnessKeyword("moderate")
	lcd.ColorKeyword("magenta")
	lcd.Print("Later, tater!")
	time.Sleep(time.Second)
	lcd.BrightnessKeyword("off")
	lcd.Clear()

}
