package main
import (
    "bufio"
    "encoding/json"
    "os"
    "os/exec"
    "strconv"
    "time"
    log "github.com/sirupsen/logrus"

)

func main() {
    logFile := "/var/log/status"
    mac := "f0:5f:77:f7:21:39"
    oldStatus := getPreviousStatus(logFile)
    danHome := obsessivelyCheckForDan(mac)
    if danHome != oldStatus {
        _ = logStatus(danHome, logFile)
    }

//      if danHome {
//              do stuff (send shutdown signal, log)
}

func getPreviousStatus(logFile string) (oldStatus bool) {
    lastLog := getLastLogLine(logFile)
    lastLogJson := LogJson{}
    json.Unmarshal([]byte(lastLog), &lastLogJson)

    oldStatus,_ = strconv.ParseBool(lastLogJson.Status)

    return

}

func getLastLogLine(f string) (s string) {
    s = ""
    file, err := os.Open(f)
    if err != nil {
        return
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if scanner.Text() != "" {
            s = scanner.Text()
        }
    }
    return
}

func pingPhone(mac string) (danStatus bool) {
    cmd := exec.Command("/usr/bin/l2ping", "-c1", "-d0", mac)
    err := cmd.Run()

    if err != nil {
        danStatus=false
        return
    }
    danStatus = true
    return

}

func obsessivelyCheckForDan( mac string ) (danStatus bool) {
    for retries:=1; retries<5; retries++ {
        danStatus = pingPhone(mac)
        if danStatus {
            break
        }
        time.Sleep(9*time.Second)
    }
    return

}

func logStatus(danStatus bool, logFile string) bool {
    var file, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
    if err != nil {
        //fmt.Println("Could not open logfile:" + err.Error())
        return false
    }
    log.SetFormatter(&log.JSONFormatter{})
    log.SetOutput(file)

    var status string
    if danStatus {
        status = "true"
    } else {
        status = "false"
    }
    logData := log.Fields{
        "status": status,
    }
    log.WithFields(logData).Info()
    return true

}

type LogJson struct {
	Level  string `json:"level"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
	Time   string `json:"time"`
}
