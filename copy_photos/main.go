package main
import (
    "fmt"
    "os"
    "golang.org/x/net/html"
    "slices"
    "strings"
    "github.com/moby/sys/mount"
    "io/fs"
    "time"
)

const blkidCache = "/run/blkid/blkid.tab"
const fsTypeKey = "type"
const mode = "ro"

const fsLabelKey = "label"
const partLabel = "nikon"

const tempDir = "temp"
const tempPrefix = "camera_sd"

const dateDirFormat = "20060102"
const targetDir = "mnt/media/photos/sort/"

func fsTypes() ([]string) {
    return []string{"exfat", "fat32"}
}

func mountPointName() string {
    return fmt.Sprintf("%s/%s.%x",tempDir,tempPrefix, time.Now().Unix())
}

func main() {
    var rc int
    defer func() {
        os.Exit(rc)
    }()

    mountPoint := mountPointName()
    fmt.Println(mountPoint)

    fsRoot, err := findAndMountDisk(blkidCache, mountPoint)
    if err != nil {
        fmt.Printf("Error finding or mounting an applicable block device: %v\n", err)
        rc = 1
        return
    }

    err = fs.WalkDir(fsRoot, ".", walkies)
    if err != nil {
        rc = 1
        fmt.Printf("Error traversing filesystem: %v\n", err)
        return
    }

    err = mount.Unmount(mountPoint)
    if err != nil {
        fmt.Printf("Error unmounting filesystem: %v\n", err)
        rc = 1
        return
    }
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
    fmt.Println("target disk ", targetDisk,", filesystem ", fsType)

    err = os.Mkdir(mountPoint, 0755)
    if err != nil {
        err = fmt.Errorf("Unable to create temp mountpoint: %v", err)
        return
    }

    err = mount.Mount(targetDisk, mountPoint, fsType, mode)
    if err != nil {
        err = fmt.Errorf("Unable to mount disk to mountpoint: %v", err)
        return
    }

    targetFS = os.DirFS(mountPoint)
    return
}

func walkies(path string, d fs.DirEntry, err error) (errr error) {
    if errr != nil {
        errr = fmt.Errorf("error on path:%v",errr)
        return
    }
    if d.IsDir() {
        return
    }
    err = process(path,d)
    if errr != nil {
        errr = fmt.Errorf("error operating on %s: %s\n", path, err)
        return
    }
    return
}

func process(p string, d fs.DirEntry) (err error) {
    size, mtime, err := stat(d)
    if err != nil {
        return
    }
    dateDir := mtime.Local().Format(dateDirFormat) + "/"
    fmt.Println("cp ", p, size, targetDir + dateDir + d.Name())
    return
}

func stat(d fs.DirEntry) (size int64, mtime time.Time, err error) {
    stat, err := d.Info()
    if err != nil {
        return
    }
    size = stat.Size()
    mtime = stat.ModTime()
    return
}

func getDiskPath(p string) (diskId, diskFS string, err error) {
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
        ok, fs := validateAttrs(n); if ok {
            diskId = n.Data
            diskFS = fs
        }
    }
    if diskFS == "" {
        err = fmt.Errorf("Failed to identify disk filesystem: %s", diskId)
    }
    return
}

func validateAttrs(n *html.Node) (ok bool, fs string) {
    if n.Type == html.TextNode && len(n.Parent.Attr) > 0 {
        var hasCorrectLabel, hasCorrectFS bool

        for _, v := range n.Parent.Attr {
            if v.Key == fsLabelKey {
                if strings.HasPrefix(strings.ToLower(v.Val), partLabel) {
                    hasCorrectLabel = true
                }
            }
            if v.Key == fsTypeKey {
                if slices.Contains(fsTypes(), v.Val) {
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
