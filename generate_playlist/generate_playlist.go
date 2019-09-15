package main

import (
    "fmt"
    taglib "github.com/wtolson/go-taglib"
    "time"
    "os"
    "path/filepath"

)

func main() {
    musicRoot := "/mnt/music"


    err := filepath.Walk(musicRoot,
        func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() == true { return err }
        fmt.Println(path)
        return nil
    })
    if err != nil {
        fmt.Println(err)
    }

/*
    songFile, err := taglib.Read(file)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    var track trackInfo

    track.FilePath = file
    track.Artist = songFile.Artist()

    fmt.Println(track)
*/
}


func scanFile(filename string) (err error, track trackInfo) {
    songFile, err := taglib.Read(filename)

    if err != nil {
        fmt.Println(err)
        return
    }

    var songData trackInfo
    songData.FilePath = filename
    songData.Artist = songFile.Artist()
    songData.Album = songFile.Album()
    songData.Title = songFile.Title()
    songData.Duration = songFile.Length()
    songData.Comment = songFile.Comment()
    songData.Genre = songFile.Genre()
    songData.Bitrate = songFile.Bitrate()
    songData.Year = songFile.Year()


    return
}

type trackInfo struct {
    FilePath    string
    Artist      string
    Album       string
    Title       string
    Comment     string
    Genre       string
    Duration    time.Duration
    Bitrate     int
    Year        int
}
