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
