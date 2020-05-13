package main

import (
    "fmt"
    "os"
    "bufio"
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
        append(fishies, temp)
    }

}

func parseFish(s string) (f fish) {
    p := strings.Split(s, "\t")
    f.Name = p[0]
    f.Price =


}

func getPrice(s string) (i int) {


}

type Fish struct {
    Name string
    Price int
    Location    float32
    Size        float32
    Times       []int
    Months      []int
}
