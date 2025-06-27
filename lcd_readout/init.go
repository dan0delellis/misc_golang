package main

import (
	"github.com/augustoroman/serial_lcd"
	"regexp"
)

const cols = 16
const rows = 2

const newlineToCarriage = true

// InitLCD is a function to open the given path at a standard baud rate
// It sets the dimensions of the display to prescribed values. Curently 16x2 because that is the device that I have, and also because it's easier to understand what's going on when a 20x4 display is only showing 16x2 characters, than when a 16x2 display is operating on the factory default setting of 20x4
// It wraps the serial_lcd.LCD, which is really just a io.ReadWriteCloser, in a custom type so I can add functions to it
// It retains functions of both io.ReadWriteCloser and serial_lcd.LCD types because we, as a species, just can't get way from the concept of inheritance
//
// If newlineToCarriage is set to true, it will create a regex pointer that will match all single-line breaks and replace them with \r\n, as that is the sequence that is interpreted as a single line break on unix and windows
// The display will be cleared, and the background color, contrast, and brightness will be set to prescribed values
// Both the blinking block and static underline cursors will be turned off
func InitLCD(path string) (d Display, err error) {
	if newlineToCarriage {
		// single lines: \r, \n, \r\n
		// double lines: \r\r, \n\n, \n\r
		singleLineBreak = regexp.MustCompile(singleLinePattern)
	}
	var l serial_lcd.LCD
	l, err = serial_lcd.Open(path, 9600)
	if err != nil {
		return
	}
	d = Display{l}
	d.Clear()
	d.NoCursor()
	d.SetSize(cols, rows)
	d.SetAutoscroll(false)
	d.SetBrightness(BrightnessMap["faint"])
	d.SetContrast((212))
	d.ColorKeyword("gold")
	d.Raw(serial_lcd.COMMAND, 0x44)

	return
}

// TODO: make this a struct that is a [] of [16]byte but only allows `rows` elements of it to be populated
type Display struct {
	serial_lcd.LCD
}

// BlinkyBlock and BlinkyBlockOff will both turn off the underline cursor

// BlinkyBlock enables a blinking block cursor at the current insertion point
func (d Display) BlinkyBlock() {
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.UNDERLINE_CURSOR_OFF))
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.BLOCK_CURSOR_ON))
}

// BlinkyBlockOff stops the blinking block cursor at the current insertion point.
func (d Display) BlinkyBlockOff() {
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.UNDERLINE_CURSOR_OFF))
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.BLOCK_CURSOR_OFF))
}

//StaticUnderline and StaticUnderlineOff will both turn off the blinking block cursor

// StaticUnderline turns on the non-blinking underscore at the current insertion point
// There is no built-in command to make this cursor blink
func (d Display) StaticUnderline() {
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.BLOCK_CURSOR_OFF))
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.UNDERLINE_CURSOR_ON))
}

// StaticUnderlineOff turns off the non-blinking underscore at the current insertion point
// There is no built-in command to make this cursor blink
func (d Display) StaticUnderlineOff() {
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.UNDERLINE_CURSOR_OFF))
	d.Raw(serial_lcd.COMMAND, byte(serial_lcd.BLOCK_CURSOR_OFF))
}

// NoCursor turns off both the blinking block cursor and the static underline cursor
func (d Display) NoCursor() {
	d.StaticUnderlineOff()
	d.BlinkyBlockOff()
}

// BothCursors enables both the blinking block cursor and non-blinking underline cursor
func (d Display) BothCursors() {
	d.StaticUnderline()
	d.BlinkyBlock()
}
