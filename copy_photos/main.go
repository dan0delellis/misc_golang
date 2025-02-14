package main
//TODO: make sure intermediary directories have the right ownership
//TODO: clean up the process() func
//TODO: verify nikon fs with that 512 byte zero file
//TODO: clean up mnt dir
//TODO: close opened files
import (
    "fmt"
    "os"
    "io/fs"
    "time"
    "github.com/moby/sys/mount"
    "syscall"
)

const tempDir = "temp"
const tempPrefix = "camera_sd"
const blkidCache = "/run/blkid/blkid.tab"

const dateDirFormat = "20060102"

var verbose bool

var photosUid int
var photosGid int
var mountPoint string

func mountPointName() string {
    return fmt.Sprintf("%s/%s.%x",tempDir,tempPrefix, time.Now().Unix())
}

func main() {
    verbose = true
    var rc int
    defer func() {
        os.Exit(rc)
    }()

    mountPoint = mountPointName()

    var err error

    fsRoot, err := findAndMountDisk(blkidCache, mountPoint)
    if err != nil {
        fmt.Printf("Error finding or mounting an applicable block device: %v\n", err)
        rc = 1
        return
    }
    debug("found fsroot:", fsRoot)

    photosUid, photosGid, err = getIds(photosDir)
    if err != nil {
        fmt.Printf("Error reading target dir (%s) info: %w", photosDir, err)
        rc=1
        return
    }
    debug("user ids of target dir:", photosUid, photosGid)

    err = fs.WalkDir(fsRoot, ".", walkies)
    if err != nil {
        rc = 1
        fmt.Printf("Error traversing filesystem: %v\n", err)
        return
    }
    debug("done walking files")

    err = mount.Unmount(mountPoint)
    if err != nil {
        fmt.Printf("Error unmounting filesystem: %v\n", err)
        rc = 1
        return
    }
    debug("umonted disk")
}

func getIds(path string) (user, group int, err error) {
    d, err := os.Open(path)
    if err != nil {
        return
    }
    stat, err := d.Stat()
    if err != nil {
        return
    }
    fileSys := stat.Sys()

    a := fileSys.(*syscall.Stat_t)
    user = int(a.Uid)
    group = int(a.Gid)
    return
}
func debug(a ...any) {
    if verbose {
        fmt.Println(a)
    }
}
