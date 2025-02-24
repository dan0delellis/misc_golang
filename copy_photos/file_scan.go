package main

import (
    "io/fs"
    "fmt"
    "os"
)

const photosDir     = "horkhorkhork"
const sortDir       = "sort"
const archiveDir    = ".copy_target"
const minSize       = 1_000

func getLinkDirs() []string {
    return []string{".copy_target", "sort"}
}

//operates on a single dir entry, this function is specific to source files on the mounted SD card
func findFiles(queue *[]TargetFile, path string, d fs.DirEntry, err error) (error) {
    debug("working file:", path)
    var task TargetFile

    if err != nil {
        return fmt.Errorf("error on path:%v", err)
    }

    if ( d.IsDir() || ! d.Type().IsRegular() ) {
        return nil
    }

    task, err = prepareCopy(path,d)

    if err != nil {
        return fmt.Errorf("error operating on %s: %s\n", path, err)
    }

    switch task.Action {
        case NoAction:
            debug(path+" needs no action")
            return nil
        case NeedsCopy:
            debug(path+" will be copied")
            //continue
        case NeedsVerify:
            debug(path+" will be copied, if target file is an incomplete copy")
            //continue
        case Conflict:
            debug("target does not appear to be an incomplete copy of "+path)
            return fmt.Errorf("conflict detected: %s has different contents than %s",task.TargetFile, path)
        default:
            debug("i have no idea what do with:", path, task.TargetFile)
            return fmt.Errorf("unhandled status for copying %s to %s: %d", path, task.TargetFile, task.Action)
    }

    task.SourceFile = path
    *queue = append(*queue, task)
    return nil
}

func prepareCopy(p string, d fs.DirEntry) (target TargetFile, err error) {
    var tgt TargetFile
    srcErr := tgt.Generate(d, getLinkDirs())
    if srcErr != nil {
        err = fmt.Errorf("Error generating target paths from source file (%s) info: %v", d.Name(), srcErr)
        return
    }
    debug("got src file meta:", d.Name())
    if tgt.SourceInfo.Size() < minSize {
        debug("this file is below minsize", tgt.SourceInfo.Size(), minSize)
        return
    }

    //does the file already exist?
    fileInfo, eNoEnt := os.Stat(tgt.TargetFile)
    if eNoEnt == nil {
        //no error means a file is already there
        tgt.TargetStat = fileInfo
        tgt.Action = compareSrcTgt(tgt.SourceInfo, fileInfo)
    } else {
        tgt.Action = NeedsCopy
    }

    debug(fmt.Sprintf("Op level for file %s is %d", tgt.TargetFile, tgt.Action))

    target = tgt
    return
}

//This could be made into a function attached to the TargetFile type
func copyFromDisk(mp string, target TargetFile) (err error) {
    //Dont create directories until we know they are required
    dirErr := target.MakePaths()
    if dirErr != nil {
        err = dirErr
        return
    }

    debug("made target/src dir", target.ArchivePath)

    //copy from mounted filesystem into archive
    raw, readErr := readData(mp+"/"+target.SourceFile, target.SourceInfo.Size())
    if readErr != nil {
        err = readErr
        return
    }
    debug("correct number of bytes read:", target.SourceInfo.Size())

    writeErr := writeData(raw, target.ArchiveFile)
    if writeErr != nil {
        err = writeErr
        return
    }
    debug("wrote file ", target.TargetFile)

    //hardlink archive to human-friendly
    linkErr := os.Link(target.ArchiveFile, target.SortFile)
    if linkErr != nil {
        err = fmt.Errorf("Error making hardlink (%s) copy of (%s): %v", target.SortFile, target.ArchiveFile, linkErr)
        return
    }
    debug("made hardlink", target.SortFile)

    //force atime/mtime/ctime to be `mtime.Unix()` in syscall.Utime
    timeErr := os.Chtimes(target.ArchiveFile, target.SourceInfo.ModTime(), target.SourceInfo.ModTime())
    if timeErr != nil {
        err = fmt.Errorf("error setting atime/mtime for copied file (%s) to match source: %v", target.ArchiveFile, timeErr)
        return
    }
    debug("set time")

    return
}

func readData(src string, expectedSize int64) (raw []byte, err error) {
    debug("reading file", src)
    raw, err = os.ReadFile(src)
    if err != nil {
        err = fmt.Errorf("Error trying to read file contents of %s", src, err)
        return
    }
    readSize := int64(len(raw))
    if readSize != expectedSize {
        err = fmt.Errorf("Wrong number of bytes read (%d) from %s (%d expected)", readSize, src, expectedSize)
        return
    }
    debug("raw data is of correct expected size:", readSize)
    return
}

func writeData(data []byte, target string) (err error) {
    err = os.WriteFile(target, data, 0664)
    if err != nil {
        err = fmt.Errorf("Error trying to write %s: %v", target, err)

        delErr := os.Remove(target)
        if delErr != nil {
            err = fmt.Errorf("%v. Also failed to delete partial file: %v", err, delErr)
        }
        return
    }

    return
}

func compareSrcTgt(src, tgt fs.FileInfo) int {
    if tgt.Size() < src.Size() {
        debug("target is smaller than source, possible incomplete copy")
        //This may have been an incomplete filecopy
        return NeedsVerify
    }
    if tgt.Size() == src.Size() && tgt.ModTime() == src.ModTime() {
        debug("target and source have same size/mtime")
        //These are likely the same file
        return NoAction
    }

    //some kind of conflict
    debug("target is larger than source and/or has different mtime")
    return Conflict
}

func setPermissions(path string, uid, gid int, bitMode fs.FileMode) (error) {
    chOwnErr := os.Chown(path, uid, gid)
    if chOwnErr != nil {
        return fmt.Errorf("error trying to change owner of copied file (%s) to match dir owner:group (%d:%d)", path, uid, gid)
    }
    debug("changed owner of ", path)

    chmodErr := os.Chmod(path, bitMode)
    if chmodErr != nil {
        return fmt.Errorf("Error setting permissions on %s to %v: %v", path, bitMode, chmodErr)
    }
    debug("set permissions on ", path)

    return nil
}
