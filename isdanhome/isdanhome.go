package main
import (
    "fmt"
    "os/exec"
    "time"
    "bytes"
    log "github.com/sirupsen/logrus"

)

func main() {
    mac := "40:4e:36:bb:c3:10"
    danHome := obsessivelyCheckForDan(mac)
    if danHome {
            fmt.Println("Dan is totally home")
        } else {
            fmt.Println("Dan is probably not home")
        }
}


func pingPhone(mac string) (danStatus bool) {
    cmd := exec.Command("/usr/bin/l2ping", "-c1", "-d0", mac)
    var errbuf bytes.Buffer
    cmd.Stderr = &errbuf
    err := cmd.Run()

    if err != nil {
        fmt.Print(errbuf.String())
        danStatus=false
        return
    }
    danStatus = true
    return

}

func obsessivelyCheckForDan( mac string ) (danStatus bool) {
    for retries:=1; retries<5; retries++ {
        danStatus = pingPhone(mac)
        time.Sleep(9*time.Second)
        if danStatus {
            break
        }
    }
    return

}

func logStatus(danStatus bool) {
    var mesg,time string
    if danStatus{
        mesg = "Dan came home!"
    } else {
        mesg = "Dan went away"
    }
    time = fmt.Println(time.Now().Format(time.RFC3339))


}
