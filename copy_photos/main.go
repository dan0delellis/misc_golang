package main
//TODO: flags
//TODO: report to redis for status readout
//TODO: image similarity comparison
//TODO: clean up the process() func
//TODO: verify nikon fs with that 512 byte zero file
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

func mountPointName() string {
    return fmt.Sprintf("%s/%s.%x",tempDir,tempPrefix, time.Now().Unix())
}

func main() {
    verbose = true
    var rc int
    var err error

    defer func() {
        os.Exit(rc)
    }()

    photosUid, photosGid, err = getIds(photosDir)
    if err != nil {
        fmt.Printf("Error getting uid/gid of target root dir: %v", err)
        rc=1
        return
    }
    debug("user ids of target dir:", photosUid, photosGid)

    mountPoint := mountPointName()

    fsRoot, err := findAndMountDisk(blkidCache, mountPoint)
    if err != nil {
        fmt.Printf("Error finding or mounting an applicable block device: %v\n", err)
        rc = 1
        return
    }
    debug("found fsroot:", fsRoot)
    defer func() {
        err = mount.Unmount(mountPoint)
        if err != nil {
            fmt.Printf("Error unmounting filesystem: %v\n", err)
            rc = 1
            return
        }
        debug("umonted disk")
    }()

    if err != nil {
        fmt.Printf("Error reading target dir (%s) info: %v", photosDir, err)
        rc=1
        return
    }

    var fileQueue []TargetFile
//    err = fs.WalkDir(fsRoot, ".", walkies)
    err = fs.WalkDir(fsRoot, ".", func(path string, entry fs.DirEntry, err error) error {return walkies(&fileQueue, path, entry, err)} )
    if err != nil {
        rc = 1
        fmt.Printf("Error traversing filesystem: %v\n", err)
        return
    }
    debug("done walking files")
    for _, v := range fileQueue {
        fmt.Println(v.SourceFile, v.ArchiveFile, v.SortFile)
    }

    debug("forcing owner/perms")
    photoDirRoot := os.DirFS(photosDir)

    err = fs.WalkDir(photoDirRoot, ".", setPermissions)
    if err != nil {
        fmt.Printf("Error forcing owner/perms of content: %v\n", err)
        rc = 1
        return
    }
    debug("done")
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
