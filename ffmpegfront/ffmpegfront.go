package main

import (
    "encoding/json"
    "log"
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "regexp"
    "strings"
    "os"
    "os/exec"
    "time"
    //"reflect"
)

var MakeTemplate = flag.Bool("make-template", false, "Write a template file")
var argsOnly = flag.Bool("args-only", false, "Output the arguments instead of executing ffmpeg with them.")
var inFile = flag.String("infile", "", "File to process with ffmpeg")
var outFile = flag.String("outfile", "", "File to write output to")
var settingsFile = flag.String("settings", "", "settings json file to read.")
var logFile = flag.String("logfile", fmt.Sprintf("/mnt/tlstc/4tb/download/rencode/logs/%s.log", inFile), "log file to write to")


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
    log.Printf("%s is not a preprogramed resolution\n", res)
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
        log.Println("Need the following flags to be used:\n\t-infile [file to process]\n\t-outfile [output target]\n\t-settings [settings json to use]\n\nOr, call with the make-template flag for it to spit out a template JSON to fill in")
        os.Exit(1)
    }

    f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
	log.Println(err)
    }
    defer f.Close()
    log := log.New(f, "ffmpegfront", log.LstdFlags)

    settings := parseSettingsJson(*settingsFile)

    args := []string{"-i", *inFile}

    if !settings.Ready.NoOverwrite {
        args = append(args, "-y")
    }

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

    if(settings.Video.JustCopy) {
        args = append(args, []string{"-c:v", "copy"}...)
    } else {
        videoArgs := parseVideoSettings(settings.Video, settings.Subtitles, *inFile)
        args = append(args, videoArgs...)
    }

//This needs to happen last before executing the command:
    args = append(args, *outFile)

    log.Printf("executing with these arguments: %v", args)
    cmd := exec.Command("/usr/bin/ffmpeg", args...)
    startTime := time.Now()
    output, err2 := cmd.CombinedOutput()
    log.Printf("finished with exit status: %v", err)
    if err2 != nil {
        log.Printf("output: %s", string(output))
    }
    duration := time.Since(startTime)
    log.Printf("Time elapsed: %s\n", duration)
}

func parseVideoSettings(v Video, s Subtitles, f string) (args []string) {
//subtitles options look like this: `-vf "subtitles=subs.srt:force_style='FontName=ubuntu,Fontsize=24,PrimaryColour=&H0000ff&'"`, so this string needs to get built :/
//also subtitle and scaling need to be part of the same filter so thats just great
    if !v.SoftwareEncode {
        args = append(args, []string{"-c:v", "h264_omx"}...)
    }

    if v.VideoBitrate != "" {
        args = append(args, []string{"-b:v", v.VideoBitrate}...)
    }

    if v.Resolution != "" || s.BurnInSubtitles {
        filter := ""
        if v.Resolution != "" {
        var res string
            regex := regexp.MustCompile(`^[0-9]*:[0-9]*$`)
            if regex.MatchString(v.Resolution) {
                res = v.Resolution
            } else {
                res = resolutionMap(v.Resolution)
            }

            filter = fmt.Sprintf("%sscale=%s", filter, res)
        }

        if s.BurnInSubtitles {
            var subFile string
            filter = fmt.Sprintf("%s, subtitles='", filter)

            if s.SubtitleFile == "" {
                subFile = f
            } else {
                subFile = s.SubtitleFile
            }
            filter = fmt.Sprintf("%s%s", filter, subFile)

            if s.SubtitleStyle != "" {
                filter = fmt.Sprintf("%s:force_style=%s", filter, s.SubtitleStyle)
            }

            filter = fmt.Sprintf(`%s'`, filter)
        }
        filter = fmt.Sprintf(`%s`, filter)
        args = append(args, []string{"-vf", filter}...)
    }

    return
}

func parseAudioSettings(a Audio, file string) (args []string) {
    var codec, bitrate, filter string

    if (a.AudioCodec != "") {
        codec = a.AudioCodec
    } else {
        codec = "aac"
    }

    args = append(args, []string{"-c:a", codec}...)

    if (a.AudioChannels != "") {
        args = append(args,[]string{"-ac", a.AudioChannels}...)
    }

    if (a.AudioBitrate != "") {
        bitrate = a.AudioBitrate
    } else {
        bitrate = "192k"
    }

    args = append(args, []string{"-b:a", bitrate}...)

    if (a.AudioFilter == "loudnorm" || a.Loudnorm2Pass) {
        if (a.Loudnorm2Pass) {
            lnJson := getLoudnormJson(file)
            filter = fmt.Sprintf("loudnorm=I=-16:TP=-1.5:LRA=11:measured_I=%s:measured_LRA=%s:measured_TP=%s:measured_thresh=%s:offset=%s:linear=true", lnJson.OutputI, lnJson.OutputLra, lnJson.OutputTp, lnJson.OutputThresh, lnJson.TargetOffset)
        } else {
            filter = "loudnorm"
        }
    }

    args = append(args, []string{"-filter:a", filter}...)

    return
}

func getLoudnormJson(file string) (lnJson loudnormValues) {
    args := []string{"-i", file, "-t", "10", "-af", "loudnorm=I=-16:TP=-1.5:LRA=11:print_format=json", "-f", "null", "-"} //those values are pretty standard and I feel OK having them hardcoded.
    cmd := exec.Command("ffmpeg", args...)
    var errb bytes.Buffer
    cmd.Stderr = &errb
    err := cmd.Run()
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    lines := strings.Split(errb.String(),"\n")
    jsonString := strings.Join(lines[len(lines)-13:len(lines)-1]," ") //The JSON data is the last 12 lines before some text in a bracket.  It would be wise to implement some form of json scanning algorithm, or deleting any text outside brackets
    jsonByte := []byte(jsonString)

    err = json.Unmarshal(jsonByte,&lnJson)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    return


}

func parseSettingsJson(file string) (settings Settings) {
    jsonFile, err := os.Open(file)
    if err != nil {
        log.Printf("unable to open json file %s: %v\n", file, err)
        os.Exit(1)
    }

    jsonBytes, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        log.Printf("unable to read json file %s: %v\n", file, err)
        os.Exit(1)
    }

    err = json.Unmarshal(jsonBytes, &settings)
    if err != nil {
        log.Printf("Unable to parse json file: %v\n", err)
    }
    return
}

func writeJson(jsonData Settings, fileName string) {
    if strings.HasSuffix(fileName, ".json") != true {
        fileName = strings.Join([]string{fileName, ".json"}, "")
    }
    outData, err := json.MarshalIndent(jsonData, "", "  ")
    if err != nil {
        log.Printf("Nope can't marshal that, %s\n", err)
        return
    }
    err2 := ioutil.WriteFile(fileName, outData, 0644)
    if err2 != nil {
        log.Printf("Failed to write file %s, %s\n", fileName, err)
    }

}

func makeEmptySettings() Settings {
    empty := Settings{
        Video{
            true,
            false,
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
            "ex-file.srt, file.mkv.  It will burn the first subtitle track if given a video file. If you want to burn in a different track, then you'll need to extract it from the video file and specify it.  If you need more complicated options, do it manually ¯\\_(ツ)_/¯",
            "styles look like this: 'FontName=ubuntu,Fontsize=24,PrimaryColour=&H0000ff&' note that the hex is BRG because fuck you that's why",
        },
        Time{
            0,
            0,
        },
        Ready{
            false,
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
    SoftwareEncode bool `json:"softwareEncode"`
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
    SubtitleStyle string `json:"subtitleStyle"`
}
type Time struct {
    TimeSkipIntro int `json:"timeSkipIntro"`
    TotalTime   int `json:"totalTime"`
}
type Ready struct {
    NoOverwrite bool `json:"noOverwrite"`
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
