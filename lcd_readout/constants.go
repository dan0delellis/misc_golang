package main

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
    BrightnessNames[0]:	uint8(0),
    BrightnessNames[1]:	uint8(16),
    BrightnessNames[2]:	uint8(32),
    BrightnessNames[3]:	uint8(64),
    BrightnessNames[4]:	uint8(96),
    BrightnessNames[5]:	uint8(128),
    BrightnessNames[6]:	uint8(160),
    BrightnessNames[7]:	uint8(192),
    BrightnessNames[8]:	uint8(255),
}
