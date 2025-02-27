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

func findAndMountDisk(cache, mountPoint string, fstypes, labels []string) (targetFS fs.FS, err error) {
    targetDisk, fsType, err := getDiskPath(cache, fstypes, labels)
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

func validateAttrs(n *html.Node, fstypes, prefixes []string) (ok bool, fs string) {
    if n.Type == html.TextNode && len(n.Parent.Attr) > 0 {
        var hasCorrectLabel, hasCorrectFS bool

        debug(fmt.Sprintf("examining attributes for %s", n.Data))
        for _, v := range n.Parent.Attr {
            debug(fmt.Sprintf("comparing attr %s to %s", v.Key, fsLabelKey))
            if v.Key == fsLabelKey {
                hasCorrectLabel = findLabelWithPrefix(prefixes, v.Val)
            }

            if v.Key == fsTypeKey {
                hasCorrectFS, fs = findFsType(fstypes, v.Val)
            }
            debug("found fs with expected label:", hasCorrectLabel, "; found fs with expected type:", hasCorrectFS)
            if hasCorrectLabel && hasCorrectFS && fs != "" {
                ok = true
            }
        }
        debug(fmt.Sprintf("done with  %s", n.Data))
    }

    return
}

func findFsType(fstypes []string, fs string) (bool, string) {
    if slices.Contains(fstypes, fs) {
        return true, fs
    }

    return false, ""
}

func findLabelWithPrefix(prefixes []string, label string) bool {
    for _, prefix := range prefixes {
        debug(fmt.Sprintf("is %s a prefix to '%s'", prefix, label))
        if strings.HasPrefix(strings.ToLower(label), prefix) {
            return true
        }
    }
    return false
}
