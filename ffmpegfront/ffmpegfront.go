package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    //"reflect"
)

func main() {
    exampleJson := Settings{}
    exampleJson.Video.Resolution = "720p"
    exampleJson.Video.VideoBitrate = "2000k"

    exampleJson.Audio = Audio{false, "vorbis", "2", "loudnorm", "200k", true}

    exampleJson.Subtitles = Subtitles{false, "ooga.srt", 12, "0xffffff", "white"}
    exampleJson.Time = Time{15, 3600}
    exampleJson.Ready.Completed = true

    emptyJson := Settings{}
    emptyJson.Video = Video{false, "", ""}
    emptyJson.Audio = Audio{true, "", "", "", "", false}
    emptyJson.Subtitles = Subtitles{false, "", 0, "", ""}
    emptyJson.Time = Time{0, 0}
    emptyJson.Ready = Ready{false}

    emptyJsonfile, _ := json.MarshalIndent(emptyJson, "", " ")
    ioutil.WriteFile("empty.json", emptyJsonfile, 0666)

    exampleJsonFile, _ := json.MarshalIndent(exampleJson, "", " ")
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
}
