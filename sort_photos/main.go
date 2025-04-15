package main

import (
    _ "flag"
    "fmt"
    _ "strings"
    "os"
    _ "time"
    "slices"
    "maps"

    _ "gocv.io/x/gocv"
    _ "gocv.io/x/gocv/contrib"
)

const (
    sortDir = "/mnt/media/photos/sort/20250111"
    RFC3339Micro = "2006-01-02T15:04:05.999999"
)

func main() {
    var ec int
    var msg error
    defer func() {
        if msg != nil {
            ec = 1
            fmt.Println(msg)
        }
        os.Exit(ec)
    } ()
    workDir, dirErr := os.Open(sortDir)
    if dirErr != nil {
        msg = fmt.Errorf("Failed to open target directory: %w", dirErr)
        return
    }
    dirList, listErr := workDir.ReadDir(0)
    if listErr != nil {
        msg = fmt.Errorf("Failed to list target directory: %w", dirErr)
        return
    }

    filesMap := make(map[string]string)
    for _, v := range dirList {
        if v.IsDir() { continue }

        info, infoErr := v.Info()
        if infoErr != nil {
            msg = fmt.Errorf("failed getting file metadata for %s: %w", v.Name(), infoErr)
            return
        }
        if info.Size() < 1000000 { continue }

        k := info.ModTime().Format(RFC3339Micro) +"/"+ info.Name()
        filesMap[k] = info.Name()
    }

    keys := slices.Sorted(maps.Keys(filesMap))

    for _, k := range keys {
        fmt.Println(k)
    }

}
