package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

func main() {
    var fishies []Fish

    tsv, _ := os.Open("fishies")
    scn := bufio.NewScanner(tsv)
    var lines []string

    for scn.Scan() {
        lines = append(lines, scn.Text())
    }

    for _, val := range(lines) {
        var temp Fish
        temp = parseFish(val)
        fishies = append(fishies, temp)
    }
    fmt.Println(fishies)

}

func parseFish(s string) (f Fish) {
    p := strings.Split(s, "\t")
    f.Name      = p[0]
    f.Price     = parseCost(p[1])
    f.Location  = getLocation(p[2])
    f.Size      = getSize(p[3])
    f.Times     = parseTimes(p[4])
    f.Months    = parseMonths(p[5])
    return
}

type Fish struct {
    Name        string
    Price       int64
    Location    Location
    Size        Shadow
    Times       [24]bool
    Months      [12]bool
}

type Location struct {
    Main        string
    Sub         string
}

type Shadow struct {
    Size        int64
    Fin         bool
}

func parseCost(s string) (c int64) {
    c,_ = strconv.ParseInt(s, 10, 64)
    return
}

func getLocation(code string)  (location Location) {
    loc := make(map[string]Location)
    loc["0"] = Location{"river", "all"}
    loc["0.1"] = Location{"river", "mouth"}
    loc["0.2"] = Location{"river", "cliff"}
    loc["1"] = Location{"lake", "all"}
    loc["2"] = Location{"sea", "all"}
    loc["2.1"] = Location{"sea", "pier"}
    loc["2.2"] = Location{"sea", "rain"}

    location = loc[code]
    return
}

func getSize(code string) (size Shadow) {
    s := strings.Split(code, ".")
    size.Size,_ = strconv.ParseInt(s[0], 10, 64)
    if len(s) ==2 {
        size.Fin = true
    } else {
        size.Fin = false
    }
    return
}

func parseTimes(s string) (t [24]bool) {
    p := strings.Split(s, ",")

    for _, val := range(p) {
        x,_ := strconv.ParseInt(val, 10, 8)
        t[x] = true
    }

    return
}

func parseMonths(s string) (m [12]bool) {
    p := strings.Split(s, ",")

    for key, val := range(p) {
        if val == "TRUE" {
            m[key] = true
        }
    }
    return
}
