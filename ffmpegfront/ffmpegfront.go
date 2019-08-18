package main
import (
	"fmt",
	"flag",
	"encoding/json"
)

type Settings struct {
	Video     Video     `json:"video"`
	Audio     Audio     `json:"audio"`
	Subtitles Subtitles `json:"subtitles"`
	Time      Time      `json:"time"`
	Ready     Ready     `json:"ready"`
}
type Video struct {
	JustCopy     bool `json:"justCopy"`
	Resolution   string `json:"resolution"`
	VideoBitrate string `json:"videoBitrate"`
}
type Audio struct {
	JustCopy      bool `json:"justCopy"`
	AudioCodec    string `json:"audioCodec"`
	AudioChannels string `json:"audioChannels"`
	AudioFilter   string `json:"audioFilter"`
	AuidioBitrate string `json:"auidioBitrate"`
	Loudnorm2Pass string `json:"loudnorm2Pass"`
}
type Subtitles struct {
	SubtitleFile  string `json:"subtitleFile"`
	FontSize      int `json:"fontSize"`
	FontColorHex  string `json:"fontColorHex"`
	FontColorWord string `json:"fontColorWord"`
}
type Time struct {
	TimeSkipIntro string `json:"timeSkipIntro"`
	TotalTime     string `json:"totalTime"`
}
type Ready struct {
	Completed bool `json:"completed"`
}
