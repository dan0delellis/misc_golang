package main

import (
    _ "flag"
    "fmt"
    _ "strings"
    "os"
    _ "time"
    "slices"
    "maps"

    "gocv.io/x/gocv"
    "gocv.io/x/gocv/contrib"
)

const (
    //sortDir = "/mnt/media/photos/sort/20250120_pleasantonridge/bird_branch_solo"
    sortDir = "/mnt/media/photos/sort/20250111/jpg"
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

    filesMap := make(map[string]FileHashes)
    for _, v := range dirList {
        var temp FileHashes

        if v.IsDir() { continue }

        info, infoErr := v.Info()
        if infoErr != nil {
            msg = fmt.Errorf("failed getting file metadata for %s: %w", v.Name(), infoErr)
            return
        }
        if info.Size() < 100000 { continue }
        temp.Name = info.Name()

        //k := info.ModTime().Format(RFC3339Micro) + " " +  info.Name()
        k := info.Name()
        temp.HashList = hashList()
        filesMap[k] = temp
    }

    keys := slices.Sorted(maps.Keys(filesMap))

    var prev FileHashes
    for i, file := range keys {
        k := filesMap[file]

        k.ImgMat = gocv.IMRead(sortDir + "/" + k.Name, gocv.IMReadColor)
        if k.ImgMat.Empty() {
            msg = fmt.Errorf("failed reading %s", k.Name)
            return
        }
        defer k.ImgMat.Close()

        fmt.Printf("\n%s: ", k.Name)
        for j, h := range k.HashList {
            res := gocv.NewMat()
            h.Compute(k.ImgMat, &res)
            if res.Empty() {
                msg = fmt.Errorf("failed computing hash for file: %s", k.Name)
                return
            }
            if i > 0 {
                sim := h.Compare(res, prev.Results[j])
                fmt.Printf("%T: %g, ", h, sim)
            }
            k.Results = append(k.Results, res)
        }
        prev = k
    }
}

func hashList() ([]contrib.ImgHashBase) {
    return []contrib.ImgHashBase{contrib.AverageHash{}, contrib.PHash{}, contrib.ColorMomentHash{} }

}

type FileHashes struct {
    Name    string
    ImgMat   gocv.Mat
    HashList []contrib.ImgHashBase
    Results  []gocv.Mat
}
