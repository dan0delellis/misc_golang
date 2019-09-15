package main

import (
    "fmt"
    taglib "github.com/wtolson/go-taglib"
    "time"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    musicRoot := "/mnt/music/"

    var exclusions Exclusions

    exclusions.Artist = []string{"Armik", "Johannes Linstead", "Elvis", "Foreigner"}
    exclusions.Genre = []string{"Guitar", "Classical"}
    exclusions.minTime = "2m"
    exclusions.maxTime = "10m"


    fileList := getFiles(musicRoot, exclusions)

    var filesToKeep []string

    for _,file := range fileList {
        track := scanFile(file)
        if keepTrack(track, exclusions) {
            filesToKeep = append(filesToKeep,track.FilePath)
        }
    }



    for _,keep := range filesToKeep {
        fmt.Println(keep)
    }
}

func keepTrack(file trackInfo, drop Exclusions) bool {
    var err error
    minTime, err := time.ParseDuration(drop.minTime)
    maxTime, err := time.ParseDuration(drop.maxTime)

    if err != nil {
        os.Exit(1)
    }

    for _,artist := range drop.Artist {
        if strings.Contains(file.Artist, artist) { return false }
    }

    for _,album := range drop.Album {
        if strings.Contains(file.Album, album) { return false }
    }

    for _,title := range drop.Title {
        if strings.Contains(file.Title, title) { return false }
    }

    for _,comment := range drop.Comment {
        if strings.Contains(file.Comment, comment) { return false }
    }

    for _,genre := range drop.Genre {
        if strings.Contains(file.Genre, genre) { return false }
    }

    if file.Duration < minTime || file.Duration > maxTime {
        return false
    }

    if file.Bitrate < drop.minBitrate {
        return false
    }

    return true
}

func getFiles(root string, drop Exclusions) (fileList []string) {
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }
//        fileInfo := scanFile(path)
//        if keepTrack(fileInfo, drop) {
        fileList = append(fileList, path)
//        }
        return nil

    })
    if err != nil {
        os.Exit(1)
    }

    return fileList


}

func scanFile(filename string) (songData trackInfo) {
    songFile, err := taglib.Read(filename)

    if err != nil {
        return
    }

    songData.FilePath = filename
    songData.Artist = songFile.Artist()
    songData.Album = songFile.Album()
    songData.Title = songFile.Title()
    songData.Duration = songFile.Length()
    songData.Comment = songFile.Comment()
    songData.Genre = songFile.Genre()
    songData.Bitrate = songFile.Bitrate()

    songFile.Close()

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
}
/*
type Exclusions struct {
    Artist      []string
    Album       []string
    Title       []string
    Comment     []string
    Genre       []string
    minTime     string
    maxTime     string
    minBitrate  int
}
*/
type Exclusions struct {
	Root       string   `json:"root"`
	ExArtists  []string `json:"ex_Artists"`
	ExAlbums   []string `json:"ex_Albums"`
	ExTitles   []string `json:"ex_Titles"`
	ExComments []string `json:"ex_Comments"`
	ExGenres   []string `json:"ex_Genre"`
	MinTime    string   `json:"min_time"`
	MaxTime    string   `json:"max_time"`
	MinBitRate int      `json:"min_bit_rate"`
}
