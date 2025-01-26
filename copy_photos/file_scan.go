package main

import (
    "time"
    "io/fs"
    "fmt"
    "syscall"
)

//operates on a single dir entry, this function is specific to source files on the mounted SD card
func walkies(path string, d fs.DirEntry, err error) (errr error) {
    if errr != nil {
        errr = fmt.Errorf("error on path:%v",errr)
        return
    }
    if d.IsDir() {
        return
    }
    err = process(path,d)
    if errr != nil {
        errr = fmt.Errorf("error operating on %s: %s\n", path, err)
        return
    }
    return
}

func stat(d fs.DirEntry) (size int64, mtime time.Time, err error) {
    stat, err := d.Info()
    if err != nil {
        return
    }
    size = stat.Size()
    mtime = stat.ModTime()

    fileSys := stat.Sys()
    a := fileSys.(*syscall.Stat_t)
    fmt.Println(a.Gid,a.Uid)
    return
}
