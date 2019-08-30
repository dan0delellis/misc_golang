package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "strings"
    "os"
    //"reflect"
)

var MakeTemplate = flag.Bool("make-template", false, "Write a template file")

func main() {
    flag.Parse()
    settingsFile := "settings.json"

    if *MakeTemplate {
        emptyJson := makeEmptySettings()
        writeJson(emptyJson, "template.json")
        os.Exit(0)
    } else {
        settings := parseSettingsJson(settingsFile)
        fmt.Println(settings)
    }




}

func parseSettingsJson(file string) (settings Settings) {
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

    json.Unmarshal(jsonBytes, &settings)
    return
}

func writeJson(jsonData Settings, fileName string) {
    if strings.HasSuffix(fileName, ".json") != true {
        fileName = strings.Join([]string{fileName, ".json"}, "")
    }
    outData, err := json.MarshalIndent(jsonData, "", "  ")
    if err != nil {
        fmt.Printf("Nope can't marshal that, %s\n", err)
        return
    }
    err2 := ioutil.WriteFile(fileName, outData, 0644)
    if err2 != nil {
        fmt.Printf("Failed to write file %s, %s\n", fileName, err)
    }

}

func makeEmptySettings() Settings {
    empty := Settings{
        Video{
            true,
            "ex-480p, 720p, 1080p, 4k",
             "ex-2000k",
            false,
        },
        Audio{
            true,
            "ex-vorbis, lame, aac, flac",
            "ex- 2, 5.1",
            "ex- loudnorm, might just make this a boolean 'UseLoudnorm'",
            "ex- 200k",
            false,
        },
        Subtitles{
            false,
            "ex-file.srt, though I need to figure out how to handle subtitles",
            12,
            "ex-ff00ff",
            "ex-white, black, red, etc",
        },
        Time{
            0,
            0,
        },
        Ready{
            false,
            "if 'JustCopy' is set as true on either audio or video settings, all other settings will be ignored.  Loudnorm2pass will be ignored if audiofilter is not set to 'loudnorm'.  Subtitles are hard to work with and i might delete that setting",
            "list of files to operate on. * will do all files in the directory.  'donkeyboner*' will match all files that start with 'donkeyboner'. '*.mkv' will operate on all mkv files.",
        },
    }
    return empty
}

type Settings struct {
    Video   Video   `json:"video"`
    Audio   Audio   `json:"audio"`
    Subtitles Subtitles `json:"subtitles"`
    Time   Time   `json:"time"`
    Ready   Ready   `json:"ready"`
}
type Video struct {
    JustCopy   bool  `json:"justCopy"`
    Resolution  string `json:"resolution"`
    VideoBitrate string `json:"videoBitrate"`
}
type Audio struct {
    JustCopy   bool  `json:"justCopy"`
    AudioCodec  string `json:"audioCodec"`
    AudioChannels string `json:"audioChannels"`
    AudioFilter  string `json:"audioFilter"`
    AudioBitrate string `json:"auidioBitrate"`
    Loudnorm2Pass bool  `json:"loudnorm2Pass"`
}
type Subtitles struct {
    BurnInSubtitles bool  `json:"burnInSubtitles"`
    SubtitleFile  string `json:"subtitleFile"`
    FontSize    int  `json:"fontSize"`
    FontColorHex  string `json:"fontColorHex"`
    FontColorWord  string `json:"fontColorWord"`
}
type Time struct {
    TimeSkipIntro int `json:"timeSkipIntro"`
    TotalTime   int `json:"totalTime"`
}
type Ready struct {
    Completed bool `json:"completed"`
    Notes string `json:"notes"`
}
