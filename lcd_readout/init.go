package main
import (
    "github.com/augustoroman/serial_lcd"
)

const cols = 16
const rows = 2

//initLCD is a function to open the given path at a standard baud rate
//it wraps the object, which is really just a io.ReadWriteCloser, in a custom type so I can add functions to it
func initLCD(path string) (d Display, err error) {
    var l serial_lcd.LCD
    l, err = serial_lcd.Open(path, 9600)
    if err != nil {
        return
    }
    d = Display{l}
    d.Clear()
    d.SetSize(cols, rows)
    d.SetAutoscroll(true)
    d.SetBrightness(32)
    d.SetContrast(uint8(212))
    d.SetBG(225, 0, 0)
    return
}

// Height returns the number of rows the LCD was set with
// The interface is write-only so this really just returns the unexported constant
func (d Display) Height() int { return rows }

// Width returns the number of rows the LCD was set with
// The interface is write-only so this really just returns the unexported constant
func (d Display) Width() int {return cols}

type Display struct{serial_lcd.LCD}
