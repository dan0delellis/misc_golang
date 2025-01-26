package main
import (
    "fmt"
    "os"
    "io/fs"
    "time"
    "github.com/moby/sys/mount"
    "syscall"
)

func mountPointName() string {
    return fmt.Sprintf("%s/%s.%x",tempDir,tempPrefix, time.Now().Unix())
}

func main() {
    var rc int
    defer func() {
        os.Exit(rc)
    }()

    mountPoint := mountPointName()

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
    hork, _ := os.Open(archiveDir)
    dork, _ := hork.Stat()
    ooga := dork.Sys().(*syscall.Stat_t)
    fmt.Printf("%+v\n", ooga)
}

func process(p string, d fs.DirEntry) (err error) {
    size, mtime, err := stat(d)
    if err != nil {
        return
    }

    dateDir := mtime.Local().Format(dateDirFormat) + "/"
    archivePath := fmt.Sprintf("%s/%s", archiveDir, dateDir)
    archiveFile := fmt.Sprintf("%s/%s", archivePath, d.Name())

    //does the file already exist?
    _, eNoEnt := os.Stat(archiveFile)
    if eNoEnt == nil {
        //no error means the file is already there
        return
    }

    //create archive+date dir if needed
    err = os.MkdirAll(archivePath, 0775)
    if err != nil {
        return
    }
    fmt.Println("cp ", p, size, archiveDir + dateDir + d.Name())
    //copy from mounted filesystem into archive
        //delete partial archive file if there was a failure

    //create human-friendly directory
    //hardlink archive to human-friendly
    //force file+dirs to be owned by the dir root owner
    //force atime/mtime/ctime to be `mtime.Unix()` in syscall.Utime

    return
}
