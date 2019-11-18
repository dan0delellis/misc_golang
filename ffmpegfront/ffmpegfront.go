package main

import (
    "encoding/json"
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "strings"
    "os"
    "os/exec"
    //"reflect"
)

var MakeTemplate = flag.Bool("make-template", false, "Write a template file")
var inFile = flag.String("infile", "", "File to process with ffmpeg")
var outFile = flag.String("outfile", "", "File to write output to")
var settingsFile = flag.String("settings", "", "settings json file to read.")


func resolutionMap(res string) (fullRes string) {
    resolutions := map[string]string {
        "480p" : "640:480",
        "720p" : "1280:720",
        "1080p" : "1920:1080",
        "4k"    : "3840:2160",
    }

    if resolutions[res] != "" {
        fullRes = resolutions[res]
        return
    }
    fmt.Printf("%s is not a preprogramed resolution\n", res)
    os.Exit(1)
    return ""
}

func main() {
    flag.Parse()

    if *MakeTemplate {
        emptyJson := makeEmptySettings()
        writeJson(emptyJson, "template.json")
        os.Exit(0)
    }


    if ( *inFile == "" ) || ( *outFile == "" ) || ( *settingsFile == "" ) {
        fmt.Println("Need the following flags to be used:\n\t-infile [file to process]\n\t-outfile [output target]\n\t-settings [settings json to use]\n\nOr, call with the make-template flag for it to spit out a template JSON to fill in")
        os.Exit(1)
    }

    settings := parseSettingsJson(*settingsFile)

    args := []string{"-i", *inFile}

    if settings.Time.TimeSkipIntro != 0 {
        args = append(args, []string{"-ss", fmt.Sprintf("%d",settings.Time.TimeSkipIntro)}...)
    }

    if settings.Time.TotalTime != 0 {
        args = append(args, []string{"-t", fmt.Sprintf("%d",settings.Time.TotalTime)}...)
    }

    if(settings.Audio.JustCopy) {
        args = append(args, []string{"-c:a", "copy"}...)
    } else {
        audioArgs := parseAudioSettings(settings.Audio, *inFile)
        args = append(args, audioArgs...)
    }

    fmt.Println(args)

}

func parseAudioSettings(a Audio, file string) (args []string) {
    flags := make(map[string]string)

    if (a.AudioCodec != "") {
        flags["-c:a"] = a.AudioCodec
    } else {
        flags["-c:a"] = "aac"
    }

    if (a.AudioChannels != "") {
        flags["-ac"] = a.AudioChannels
    }

    if (a.AudioBitrate != "") {
        flags["-b:a"] = a.AudioBitrate
    } else {
        flags["-b:a"] = "192k"
    }

    if (a.AudioFilter == "loudnorm" || a.Loudnorm2Pass) {
        if (a.Loudnorm2Pass) {
            lnJson := getLoudnormJson(file)
            flags["-af"] = fmt.Sprintf("loudnorm=I=-16:TP=-1.5:LRA=11:measured_I=%s:measured_LRA=%s:measured_TP=%s:measured_thresh=%s:offset=%s:linear=true", lnJson.OutputI, lnJson.OutputLra, lnJson.OutputTp, lnJson.OutputThresh, lnJson.TargetOffset)
        } else {
            flags["-af"] = "loudnorm"
        }
    }

    for key, val := range(flags) {
        args = append(args, []string{key,val}...)
    }

    return
}

func getLoudnormJson(file string) (lnJson loudnormValues) {
    args := []string{"-i", file, "-t", "10", "-af", "loudnorm=I=-16:TP=-1.5:LRA=11:print_format=json", "-f", "null", "-"} //those values are pretty standard and I feel OK having them hardcoded.
    cmd := exec.Command("ffmpeg", args...)
    var errb bytes.Buffer
    cmd.Stderr = &errb
    err := cmd.Run()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    lines := strings.Split(errb.String(),"\n")
    jsonString := strings.Join(lines[len(lines)-13:len(lines)-1]," ") //The JSON data is the last 12 lines before some text in a bracket.  It would be wise to implement some form of json scanning algorithm, or deleting any text outside brackets
    jsonByte := []byte(jsonString)

    err = json.Unmarshal(jsonByte,&lnJson)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return


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

    err = json.Unmarshal(jsonBytes, &settings)
    if err != nil {
        fmt.Printf("Unable to parse json file: %v\n", err)
    }
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
        },
        Audio{
            true,
            "ex-vorbis, lame, aac, flac",
            "ex- 2, 5.1",
            "ex- loudnorm, might just make this a boolean 'UseLoudnorm' because what other filter am I likely to use?",
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

type loudnormValues struct {
	InputI            string `json:"input_i"`
	InputTp           string `json:"input_tp"`
	InputLra          string `json:"input_lra"`
	InputThresh       string `json:"input_thresh"`
	OutputI           string `json:"output_i"`
	OutputTp          string `json:"output_tp"`
	OutputLra         string `json:"output_lra"`
	OutputThresh      string `json:"output_thresh"`
	NormalizationType string `json:"normalization_type"`
	TargetOffset      string `json:"target_offset"`
}

/*example json:
settings:
{
  "video": {
    "justCopy": false,
    "resolution": "720p",
    "videoBitrate": "2000k"
  },
  "audio": {
    "justCopy": false,
    "audioCodec": "vorbis",
    "audioChannels": "2",
    "audioFilter": "loudnorm",
    "auidioBitrate": "200k",
    "loudnorm2Pass": true
  },
  "subtitles": {
    "burnInSubtitles": false,
    "subtitleFile": "ooga.srt",
    "fontSize": 12,
    "fontColorHex": "0xffffff",
    "fontColorWord": "white"
  },
  "time": {
    "timeSkipIntro": 15,
    "totalTime": 3600
  },
  "ready": {
    "completed": true,
    "notes": "",
  }
}

loudnorm sample output:
{
    "input_i" : "-18.33",
    "input_tp" : "-7.93",
    "input_lra" : "20.70",
    "input_thresh" : "-30.40",
    "output_i" : "-23.95",
    "output_tp" : "-7.03",
    "output_lra" : "7.60",
    "output_thresh" : "-34.44",
    "normalization_type" : "dynamic",
    "target_offset" : "-0.05"
}


parse multiple jsons from string data:
https://play.golang.org/p/6XAdq6N0PAD
*/


