package main
import (
    "fmt"
    "io/fs"
    "os"
    "golang.org/x/net/html"
    "github.com/moby/sys/mount"
    "slices"
    "strings"
)

const fsTypeKey = "type"
const mountMode = "ro"

const fsLabelKey = "label"
const partLabel = "nikon"

func fsTypes() ([]string) {
    return []string{"exfat", "fat32"}
}

func getDiskPath(p string, fstypes, labels []string) (diskId, diskFS string, err error) {
    devs, err := os.Open(p)
    defer devs.Close()

    if err != nil {
        err = fmt.Errorf("Failed opening %s:%s", p, err)
        return
    }

    data, err := html.Parse(devs)
    if err != nil {
        err = fmt.Errorf("Failed parsing dev data: %s", err)
        return
    }

    for n := range data.Descendants() {
        ok, fs := validateAttrs(n, fstypes, labels); if ok {
            diskId = n.Data
            diskFS = fs
        }
    }
    if diskFS == "" {
        err = fmt.Errorf("Failed to identify disk filesystem: %s", diskId)
    }
    return
}

func findAndMountDisk(cache, mountPoint string) (targetFS fs.FS, err error) {
    targetDisk, fsType, err := getDiskPath(blkidCache)
    if err != nil {
        err = fmt.Errorf("Unable to locate applicable disk: %v", err)
        return
    }
    if targetDisk == "" {
        return
    }
    debug("target disk ", targetDisk,", filesystem ", fsType)

    err = os.Mkdir(mountPoint, 0755)
    if err != nil {
        err = fmt.Errorf("Unable to create temp mountpoint: %v", err)
        return
    }
    debug("made mountpoint", mountPoint)

    err = mount.Mount(targetDisk, mountPoint, fsType, mountMode)
    if err != nil {
        err = fmt.Errorf("Unable to mount disk to mountpoint: %v", err)
        return
    }
    debug("mounted disk")

    targetFS = os.DirFS(mountPoint)
    return
}

func validateAttrs(n *html.Node, fstypes, labels []string) (ok bool, fs string) {
    if n.Type == html.TextNode && len(n.Parent.Attr) > 0 {
        var hasCorrectLabel, hasCorrectFS bool

        for _, v := range n.Parent.Attr {
            if v.Key == fsLabelKey {
                for _, label := range labels
                    if strings.HasPrefix(strings.ToLower(v.Val), label) {
                        hasCorrectLabel = true
                    }
                }
            }
            if v.Key == fsTypeKey {
                if slices.Contains(fstypes, v.Val) {
                    hasCorrectFS = true
                    fs = v.Val
                }
            }
            if hasCorrectLabel && hasCorrectFS && fs != "" {
                ok = true
            }
        }
    }

    return
}
