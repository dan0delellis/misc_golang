package main
import (
    "fmt"
    "io/fs"
    "os"
    "golang.org/x/net/html"
    "github.com/moby/sys/mount"
    "strings"
)

const fsTypeKey = "type"
const fsLabelKey = "label"
const mountMode = "ro"
const nikonFile = "NIKON001.DSC"

func findAndMountDisks(opts *Opts, mountDir string) (targetFS fs.FS, mountedDirs []string, err error) {
    targetDisks, err := getDiskPaths(opts)
    if err != nil {
        err = fmt.Errorf("Unable to locate applicable disk: %v", err)
        return
    }
    if len(targetDisks) == 0 {
        debug("no valid disks found")
        return
    }
    debug("target disk list", targetDisks)

    for _, v := range targetDisks {
        err = v.Mount(mountDir, mountMode, opts)
        if err != nil {
            return
        }
    }

    targetFS = os.DirFS(mountDir)
    return
}

func getDiskPaths(opts *Opts) (devIDs []Dev, err error) {
    if opts.DevIDs == nil {
        //if no devid is specified, it will scan all devids
        opts.DevIDs = []string{""}
    } else {
        //if devids are specified, assume that whatever label/fstype is discovered is expected

        //label is validated by checking if discovered label has an expected label as a prefix
        //fstype is validated by checking if expected fs is empty string, or if discovered fs is an exact match of one of the expected filesystems
        opts.FsTypes = []string{""}
        opts.FsLabels = []string{""}
    }

    raw, err := os.Open(opts.BlkidCache)
    defer raw.Close()

    if err != nil {
        err = fmt.Errorf("Failed opening %s:%s", opts.BlkidCache, err)
        return
    }

    data, err := html.Parse(raw)
    if err != nil {
        err = fmt.Errorf("Failed parsing dev data: %s", err)
        return
    }

    for n := range data.Descendants() {
        ok, fs := validateAttrs(n, opts); if ok {
            devIDs = append(devIDs, Dev{DevID:n.Data, Filesystem:fs})
        }
    }
    return
}

func validateAttrs(n *html.Node, opts *Opts) (ok bool, fs string) {
    if n.Type == html.TextNode && len(n.Parent.Attr) > 0 {
        var hasCorrectLabel, hasCorrectFS bool

        debug(fmt.Sprintf("examining attributes for %s", n.Data))

        nameMatch, _ := findEmptyOrMatch(n.Data, opts.DevIDs)
        if !nameMatch {
            debugf("%s is not a specified devname", n.Data)
            return
        }

        for _, v := range n.Parent.Attr {
            debug(fmt.Sprintf("comparing attr %s to %s", v.Key, fsLabelKey))
            if v.Key == fsLabelKey {
                hasCorrectLabel = findLabelWithPrefix(opts.FsLabels, v.Val)
            }

            if v.Key == fsTypeKey {
                hasCorrectFS, fs = findFsType(v.Val, opts.FsTypes)
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

func findFsType(fs string, fstypes []string) (bool, string) {
    return findEmptyOrMatch(fs, fstypes)
}
func findEmptyOrMatch(discovered string, expected []string) (bool, string) {
    for _, exp := range expected {
        if exp == "" || exp == discovered {
            return true, discovered
        }
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

func (dev *Dev) Mount(location, mode string, opts *Opts) ( err error ) {
    mountPoint := location + "/" + dev.DevID
    err = os.MkdirAll(mountPoint, 0755)
    if err != nil {
        err = fmt.Errorf("Unable to create temp dir for mountpoint: %v", err)
        return
    }
    debug("made mountpoint", mountPoint)

    err = mount.Mount(dev.DevID, mountPoint, dev.Filesystem, mode)
    if err != nil {
        err = fmt.Errorf("Unable to mount disk to mountpoint: %v", err)
        return
    }
    debug("mounted disk")
    if opts.NikonFile {
        if opts.NikonFilePath == "" {
            opts.NikonFilePath = nikonFile
        }
    }
    if opts.NikonFilePath != "" {
        r, e := readData(mountPoint+"/"+nikonFile,512)
        if e != nil {
            err = fmt.Errorf("Called with --nikon or --nikonfile but %s not found in root of %s", nikonFile, dev.DevID)
            return
        }

        if !compareByteSlices(r, make([]byte,512)) {
            err = fmt.Errorf("called with --nikon or --nikonfile but %s at root of %s is not a 512-byte file of 0s", nikonFile, dev.DevID)
            return
        }
    }
    return
}

type Dev struct {
    DevID, Filesystem string
}
