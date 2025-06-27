package main

// Height returns the number of rows the LCD was set with
// The interface is write-only so this really just returns the unexported constant
func (d Display) Height() int { return rows }

// Width returns the number of rows the LCD was set with
// The interface is write-only so this really just returns the unexported constant
func (d Display) Width() int { return cols }
