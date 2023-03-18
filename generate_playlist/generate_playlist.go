package main

import (
    "fmt"
    taglib "github.com/wtolson/go-taglib" //requires packae libtagc0-dev; requires `go get github.com/wtolson/go-taglib`
    "flag"
    "time"
    "os"
    "path/filepath"
    "strings"
    "io/ioutil"
    "encoding/json"
)

func main() {
    settingsFile := flag.String("settings", "settings.json", "settings json file to use")
    flag.Parse()

    exclusions := parseSettingsJson(*settingsFile)

    if exclusions.Playlist == "" {
        exclusions.Playlist = "playlist.m3u8"
    }

    fmt.Println(exclusions)

    fmt.Println("getting files")
    fileList := getFiles(exclusions.Root)

    var filesToKeep []string

    for _,file := range fileList {
        track := scanFile(file)
	if (track.FilePath == "!!!!I AM NOT AN AUDIO FILE WITH ID TAGS!!!!!!") {
	    continue
	}
        if !keepTrack(track, exclusions) {
            track.FilePath = strings.Join([]string{"#",track.FilePath}, " ")
        }
	filesToKeep = append(filesToKeep, track.FilePath)


    }

    writeFile(filesToKeep,exclusions.Playlist)
    os.Exit(0)


}

func keepTrack(file trackInfo, drop Exclusions) bool {
    var err error
    minTime, err := time.ParseDuration(drop.MinTime)
    maxTime, err := time.ParseDuration(drop.MaxTime)
    fmt.Println(file.FilePath)

    if err != nil {
        os.Exit(1)
    }

    for _,artist := range drop.ExArtists {
        if strings.Contains(file.Artist, artist) {
        fmt.Printf("matched artist: %s\n", file.FilePath)
            return false
        }
    }

    for _,album := range drop.ExAlbums {
        if strings.Contains(file.Album, album) {
        fmt.Printf("matched album: %s\n", file.FilePath)
            return false
        }
    }

    for _,title := range drop.ExTitles {
        if strings.Contains(file.Title, title) {
        fmt.Printf("matched title: %s\n", file.FilePath)
            return false
        }
    }

    for _,comment := range drop.ExComments {
        if strings.Contains(file.Comment, comment) {
        fmt.Printf("matched comment: %s\n", file.FilePath)
            return false
        }
    }

    for _,genre := range drop.ExGenres {
        if strings.Contains(file.Genre, genre) {
        fmt.Printf("matched genre: %s\n", file.FilePath)
            return false
        }
    }

    if file.Duration < minTime || file.Duration > maxTime {
        fmt.Printf("bad duration: %s\n", file.FilePath)
        return false
    }

    if file.Bitrate < drop.MinBitrate {
        fmt.Printf("bad bitrate: %s\n", file.FilePath)
        return false
    }

    return true
}

func getFiles(root string) (fileList []string) {
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }
        fileList = append(fileList, path)
        fmt.Println(path)
        return nil

    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return fileList
}

func scanFile(filename string) (songData trackInfo) {
    songFile, err := taglib.Read(filename)
    songData.FilePath = "!!!!I AM NOT AN AUDIO FILE WITH ID TAGS!!!!!!"

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

func parseSettingsJson(file string) (settings Exclusions) {
    jsonFile, err := os.Open(file)
    if err != nil {
        fmt.Printf("unable to open json file %s: %v\n", file, err)
        os.Exit(1)
    }

    jsonBytes, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        fmt.Printf("unable to read json file %s: %v\n", file, err)
        os.Exit(1)
    }

    err = json.Unmarshal(jsonBytes, &settings)
    if err != nil {
        fmt.Printf("Unable to parse json file: %v\n", err)
    }
    return
}

func writeFile(playlist []string, file string) {
    str := strings.Join(playlist, "\n")
    f, err := os.Create(file)
    if err != nil {
        fmt.Println(str)
        fmt.Println("Unable to create file, so you can copy/paste the above")
        fmt.Println(err)
        os.Exit(1)
    }

    _, err = f.WriteString(str)
    if err != nil {
        fmt.Println(str)
        fmt.Println("Unable to print data to file.  You can copy/paste the above")
        fmt.Println(err)
        f.Close()
        os.Exit(1)
    }

    f.Close()

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
        Playlist   string   `json:"playlist"`
	ExArtists  []string `json:"ex_Artists"`
	ExAlbums   []string `json:"ex_Albums"`
	ExTitles   []string `json:"ex_Titles"`
	ExComments []string `json:"ex_Comments"`
	ExGenres   []string `json:"ex_Genre"`
	MinTime    string   `json:"min_time"`
	MaxTime    string   `json:"max_time"`
	MinBitrate int      `json:"min_bit_rate"`
}
