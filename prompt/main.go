package main
import (
    "fmt"
    "strings"
    "crypto/md5"
    "os"
)

const clear = "\\[\\033[00m\\]"

func main() {
    prompt := getPrompt()
    fmt.Print(prompt)
}

func getPrompt() string {
    hostname, _ := os.Hostname()
    hash := md5.Sum([]byte(hostname))

    offset := 4

    colorTime := hash[offset+0]
    colorUser := hash[offset+1]
    colorHost := hash[offset+2]
    colorPwd := hash[offset+3]

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

func foreground(i uint8, a string) string {
    return fmt.Sprintf("\\[\\033[38;05;%dm\\]%s%s", i, esc(a), clear)
}

func esc(a string) string {
    return fmt.Sprintf("\\%s", a)
}
