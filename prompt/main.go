package main
import (
    "fmt"
    "strings"
//    "math/rand"
)

const clear = "\\[\\033[00m\\]"

func main() {
    prompt := getPrompt()
    fmt.Print(prompt)
}

func getPrompt() string {
    colorTime := 4895418 % 256
    colorUser := 5475 % 256
    colorHost := 789784634 % 256
    colorPwd := 4574878 % 256
    parts := []string {
        //time
        "(" + foreground(colorTime,"t") + ") ",

        //username
        foreground(colorUser,"u"),

        // @
        "@",

        //hostname
        foreground(colorHost,"h"),

        // :
        ": ",

        //working dir
        foreground(colorPwd,"w"),

        //newline
        "\n",

        //prompt char
        esc("$ "),

    }
    return strings.Join(parts, "")
}

func foreground(i int,a string) string {
    return fmt.Sprintf("\\[\\033[38;05;%dm\\]%s%s",i,esc(a),clear)
}

func esc(a string) string {
    return fmt.Sprintf("\\%s", a)
}

