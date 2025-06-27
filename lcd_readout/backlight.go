package main

import (
	"fmt"
	"golang.org/x/image/colornames"
)

var BrightnessNames = []string{
	"off",
	"faint",
	"dim",
	"soft",
	"moderate",
	"bright",
	"vivid",
	"brilliant",
	"max",
}

var BrightnessMap = map[string]uint8{
	BrightnessNames[0]: 0,
	BrightnessNames[1]: 16,
	BrightnessNames[2]: 32,
	BrightnessNames[3]: 64,
	BrightnessNames[4]: 96,
	BrightnessNames[5]: 128,
	BrightnessNames[6]: 160,
	BrightnessNames[7]: 192,
	BrightnessNames[8]: 255,
}

// SVGColorKeyword sets the background color to the string provided, if said string has a corresponding value in the image/colornames library
// If the provided string is not a known color name, BG will not be changed
func (d Display) ColorKeyword(name string) (err error) {
	values, ok := colornames.Map[name]
	if ok {
		d.SetBG(values.R, values.G, values.B)
	} else {
		err = fmt.Errorf("%s is not a known color name. Please email www-svg@w3.org if you think it should be!")
	}
	return
}

// BrightnessKeyword sets the backlight level to one of the prescribed brightness levels, as listed in 'BrightnessNames'
// If the provided name is not listed, it will return an error
func (d Display) BrightnessKeyword(name string) (err error) {
	value, ok := BrightnessMap[name]
	if ok {
		d.SetBrightness(value)
	} else {
		err = fmt.Errorf("%s is not a prescribed brightness level. This is the list of supported names: %s", BrightnessNames)
	}
	return
}
