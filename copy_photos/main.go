package main
//TODO: report to redis for status readout
//TODO: image similarity comparison
import (
    "fmt"
    "os"
    "io/fs"
    "time"
    "github.com/moby/sys/mount"
    "syscall"
)

var verbose bool

func mountPointName(opts *Opts) string {
    return fmt.Sprintf("%s/%s.%x",opts.TempDir, opts.TempPrefix, time.Now().Unix())
}

func main() {
    var rc int

    defer func() {
        os.Exit(rc)
    }()
    opts, err := parseFlags()
    if err != nil {
        //the flags library already `helpfully` prints the error for us, even if the error is generated by the -h flag, so there's no need to print it here
        rc=1
        return
    }
    verbose = opts.Verbose

    var photosUid, photosGid int

    photosUid, photosGid, err = getIds(&opts)
    if err != nil {
        fmt.Printf("Error getting uid/gid of target root dir: %v", err)
        rc=1
        return
    }
    debug("user/group ids of target dir:", photosUid, photosGid)

    mountPoint := mountPointName(&opts)

    fsRoot, mountedDirs, err := findAndMountDisks(&opts, mountPoint)
    defer func() {
        if !opts.KeepMounts {
            for _, v := range mountedDirs {
                debugf("unmounting %s", v)
                err = mount.Unmount(v)
                if err != nil {
                    fmt.Printf("Error unmounting filesystem: %v\n", err)
                    rc = 1
                    return
                }
            }
            debug("umonted disk")
        } else {
            debug("not unmounting disks")
        }
    }()

    if err != nil {
        fmt.Printf("Error finding or mounting an applicable block device: %v\n", err)
        rc = 1
        return
    }
    debug("found fsroot:", fsRoot)
    var fileQueue []TargetFile

    err = fs.WalkDir(fsRoot, ".", func(path string, entry fs.DirEntry, err error) error {
        return findFiles(&fileQueue, opts.RootPath, path, entry, err)
    })

    if err != nil {
        rc = 1
        fmt.Printf("Error traversing filesystem: %v\n", err)
        return
    }

    debug(fmt.Sprintf("found %d files", len(fileQueue)))
    for k, v := range fileQueue {
        fmt.Println(k+1, "of", len(fileQueue))
        err = v.CopyFromDisk(mountPoint)
        if err != nil {
            fmt.Println(err)
            rc=1
            return
        }
    }

    debug("forcing owner/perms")
    photoDirRoot := os.DirFS(opts.RootPath)

    err = fs.WalkDir(photoDirRoot, ".", func(path string, d fs.DirEntry, e error) error {
        if e != nil {
            return fmt.Errorf("encountered error setting permissions: %v", e)
        }
        if ! ( d.Type().IsRegular() || d.IsDir() ) {
            return nil
        }
        var bitMode fs.FileMode
        if d.IsDir() {
            bitMode = 0775
        } else {
            bitMode = 0664
        }

        return setPermissions(opts.RootPath + "/" + path, photosUid, photosGid, bitMode)
    })

    if err != nil {
        fmt.Printf("Error forcing owner/perms of content: %v\n", err)
        rc = 1
        return
    }
    debug("done")
}

func getIds(opts *Opts) (user, group int, err error) {
    if opts.UserID > -1 && opts.GroupID > -1 {
        user = opts.UserID
        group = opts.GroupID
        return
    }
    path := opts.RootPath
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
    user = max(opts.UserID, int(a.Uid))
    group = max(opts.GroupID, int(a.Gid))
    return
}

func debug(a ...any) {
    if verbose {
        fmt.Println(a)
    }
}

func debugf(s string, a ...any) {
    debug(fmt.Sprintf(s, a...))
}

type StatusReport struct {
    TTL int
    Detail []string



}
