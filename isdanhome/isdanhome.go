# do something like this https://play.golang.org/p/8e1zNrs4a1C

package main
import (
    "fmt"
    "os/exec"
    "time"
    "bytes"

)

func main() {
//    mac := "40:4e:36:bb:c3:10"
    danHome := true


        if danHome {
            fmt.Println("dan must be home")
        } else {
            fmt.Println("dan is not home")
        }
        fmt.Println("I sleep")
        time.Sleep(10*time.Second)
}

func pingPhone(mac string) (danStatus bool) {
    cmd := exec.Command("/usr/bin/l2ping", "-c1", "-d0", mac)
    var errbuf bytes.Buffer
    cmd.Stderr = &errbuf
    err := cmd.Run()

    if err != nil {
        fmt.Println(errbuf.String())
        danStatus=false
        return
    }
    danStatus = true
    return

}

