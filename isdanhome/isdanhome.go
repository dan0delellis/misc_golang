package main
import (
    "fmt"
    "os/system"
    "github.com/juju/fslock"

)

func main() {
    var mac := "40:4e:36:bb:c3:10"


    lockFile := fslock.New("/var/lock/isdanhome")
    lockErr := lockFile.TryLock()
    if lockErr != nil {
        fmt.Println("Couldn't get lockfile!  Shit's bad yo!")
    }

    time.Sleep(10 * time.Second)
    lockFile.Unlock()

}


