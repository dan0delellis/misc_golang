package main

import (
    "time"
    "io/fs"
    "fmt"
    "syscall"
    "os"
)

const photosDir     = "mnt/media/photos"
const sortDir       = "sort"
const archiveDir    = ".copy_target"

//operates on a single dir entry, this function is specific to source files on the mounted SD card
func walkies(leafFile string, d fs.DirEntry, err error) (errr error) {
    debug("working file", leafFile)
    if errr != nil {
        errr = fmt.Errorf("error on path:%v",errr)
        return
    }

    if ( d.IsDir() || ! d.Type().IsRegular() ) {
        return
    }

    err = process(leafFile,d)
    if errr != nil {
        errr = fmt.Errorf("error operating on %s: %s\n", leafFile, err)
        return
    }
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

func (t *TargetFile) Generate(f fs.DirEntry) ( i fs.FileInfo, e error ){
    i, e = f.Info()
    if e != nil {
        return
    }
    dateDir := i.ModTime().Local().Format(dateDirFormat) + "/"

    t.ArchivePath = fmt.Sprintf("%s/%s/%s", photosDir, archiveDir, dateDir)
    t.ArchiveFile = fmt.Sprintf("%s/%s", t.ArchivePath, f.Name())

    t.SortPath    = fmt.Sprintf("%s/%s/%s", photosDir, sortDir, dateDir)
    t.SortFile    = fmt.Sprintf("%s/%s", t.SortPath, f.Name())

    return
}

//This should create the dir paths for ArchivePath and SortPath
func (t *TargetFile) MakePaths() (err error) {
    return
}

// A TargetFile struct stores the paths for the initial copy (archive) of a copied photo, and the non-hidden sortable file
// the sortable file will be a hardlink to the archive file.
type TargetFile struct {
    ArchivePath, ArchiveFile string
    SortPath, SortFile string
}

func process(p string, d fs.DirEntry) (err error) {
    var target TargetFile
    srcInfo, srcErr := target.Generate(d)
    if err != nil {
        err = fmt.Errorf("Error generating target paths from source file (%s) info: %w", d.Name(), srcErr)
        return
    }
    debug("got src file meta", d.Name())

    //does the file already exist?
    //move this to its own function
    afStat, eNoEnt := os.Stat(target.ArchiveFile)
    if eNoEnt == nil {
        //no error means a file is already there

        //is it the same file? check stat of sourcefile
        if afStat.Size() == srcInfo.Size() && afStat.ModTime() == srcInfo.ModTime() {
            //these are likely the same file, and the target file is complete
            return
        } else {
            debug(fmt.Sprintf("target file %s already exists but (%d:%s) does not match size:mtime of source (%d:%s)", target.ArchiveFile, afStat.Size(), afStat.ModTime().Local().Format(time.RFC3339), srcInfo.Size(), srcInfo.ModTime().Local().Format(time.RFC3339)))
            debug("overwriting")
        }
    }
    debug("target does not exist:", target.ArchiveFile)

    //create archive+date dir if needed
    //move this to its own function
    archDirErr := os.MkdirAll(target.ArchivePath, 0775)
    if archDirErr != nil {
        err = fmt.Errorf("Error creating archive directory %s: %w", target.ArchivePath, archDirErr)
        return
    }
    debug("made target dir", target.ArchivePath)
    //copy from mounted filesystem into archive
    raw, readErr := os.ReadFile(mountPoint+"/"+p)
    if readErr != nil {
        err = fmt.Errorf("Error trying to copy contents of %s: %w", d.Name(), readErr)
        debug(err)
        return
    }
    debug("read raw")
    if int64(len(raw)) != srcInfo.Size() {
        err = fmt.Errorf("Wrong number (%d) of bytes read from %s (%d expected)", len(raw), d.Name(), srcInfo.Size())
        return
    }
    debug("correct bytes read")
    writeErr := os.WriteFile(target.ArchiveFile, raw, 0664)
    if writeErr != nil {
        err = fmt.Errorf("Error trying to write %s: %w", target.ArchiveFile, writeErr)

        delErr := os.Remove(target.ArchiveFile)
        if delErr != nil {
            err = fmt.Errorf("%w. Also failed to delete partial file: %w", err, delErr)
        }
        debug(err)
        return
    }
    debug("wrote file")

    //create human-friendly directory
    //this should be part of the above function that makes the archive dir
    sortDirErr := os.MkdirAll(target.SortPath, 0775)
    if archDirErr != nil {
        err = fmt.Errorf("Error creating sorting directory %s: %w", target.SortPath, sortDirErr)
        return
    }
    debug("made sort dir:", target.SortPath)

    //hardlink archive to human-friendly
    linkErr := os.Link(target.ArchiveFile, target.SortFile)
    if linkErr != nil {
        err = fmt.Errorf("Error making hardlink (%s) copy of (%s): %w", target.SortFile, target.ArchiveFile)
        return
    }
    debug("made hardlink", target.SortFile)
    //force file+dirs to be owned by the dir root owner
    chOwnErr := os.Chown(target.ArchiveFile, photosUid, photosGid)
    if chOwnErr != nil {
        err = fmt.Errorf("error trying to change owner of copied file (%s) to match dir owner:group (%d:%d)", target.ArchiveFile, photosUid, photosGid)
        return
    }
    debug("changed owner")

    //force atime/mtime/ctime to be `mtime.Unix()` in syscall.Utime
    timeErr := syscall.Utime(target.ArchiveFile, &syscall.Utimbuf{Actime: srcInfo.ModTime().Unix(), Modtime:srcInfo.ModTime().Unix()})
    if timeErr != nil {
        err = fmt.Errorf("error setting atime/mtime for copied file (%s) to match source: %w", target.ArchiveFile, timeErr)
        return
    }
    debug("set time")

    return
}
