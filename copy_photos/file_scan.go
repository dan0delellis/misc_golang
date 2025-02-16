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

//operates on a single dir entry, this function is specific to source files on the mounted SD card
func walkies(path string, d fs.DirEntry, err error) (error) {
    debug("working file", path)
    if err != nil {
        return fmt.Errorf("error on path:%v", err)
    }

    if ( d.IsDir() || ! d.Type().IsRegular() ) {
        return nil
    }

    err = process(path,d)
    if err != nil {
        return fmt.Errorf("error operating on %s: %s\n", path, err)
    }
    return nil
}

//This is a walkdir func
func setPermissions(path string, d fs.DirEntry, err error) (error) {
    if err != nil {
        return fmt.Errorf("error setting permissions on path:%v", err)
    }
    if ! ( d.Type().IsRegular() || d.IsDir() ) {
        return nil
    }
    path = photosDir + "/" + path

    chOwnErr := os.Chown(path, photosUid, photosGid)
    if chOwnErr != nil {
        return fmt.Errorf("error trying to change owner of copied file (%s) to match dir owner:group (%d:%d)", path, photosUid, photosGid)
    }
    debug("changed owner")

    var bitMode fs.FileMode
    if d.IsDir() {
        bitMode = 0775
    } else {
        bitMode = 0664
    }

    chmodErr := os.Chmod(path, bitMode)
    if chmodErr != nil {
        return fmt.Errorf("Error setting permissions on %s to %v: %w", path, bitMode, chmodErr)
    }
    debug("set permissions")

    return nil
}

func process(p string, d fs.DirEntry) (err error) {
    var target TargetFile
    srcErr := target.Generate(d)
    if srcErr != nil {
        err = fmt.Errorf("Error generating target paths from source file (%s) info: %w", d.Name(), srcErr)
        return
    }
    debug("got src file meta", d.Name())
    if target.Info.Size() < minSize {
        debug("this file is below minsize")
        return
    }

    //does the file already exist?
    archiveFileInfo, eNoEnt := os.Stat(target.ArchiveFile)
    if eNoEnt == nil {
        //no error means a file is already there
        //is it the same file? check stat of sourcefile
        if compareSrcTgt(target.Info, archiveFileInfo) {
            return
        }
    }

    debug("target does not exist:", target.ArchiveFile)

    //Dont create directories until we know they are required
    dirErr := target.MakePaths()
    if dirErr != nil {
        err = dirErr
        return
    }

    debug("made target/src dir", target.ArchivePath)

    //copy from mounted filesystem into archive
    raw, readErr := readData(mountPoint+"/"+p, target.Info.Size())
    if readErr != nil {
        err = readErr
        return
    }
    debug("correct number of bytes read")

    writeErr := writeData(raw, target.ArchiveFile)
    if writeErr != nil {
        err = writeErr
        return
    }
    debug("wrote file")

    //hardlink archive to human-friendly
    linkErr := os.Link(target.ArchiveFile, target.SortFile)
    if linkErr != nil {
        err = fmt.Errorf("Error making hardlink (%s) copy of (%s): %w", target.SortFile, target.ArchiveFile, linkErr)
        return
    }
    debug("made hardlink", target.SortFile)

    //force atime/mtime/ctime to be `mtime.Unix()` in syscall.Utime
    timeErr := os.Chtimes(target.ArchiveFile, target.Info.ModTime(), target.Info.ModTime())

    if timeErr != nil {
        err = fmt.Errorf("error setting atime/mtime for copied file (%s) to match source: %w", target.ArchiveFile, timeErr)
        return
    }
    debug("set time")

    return
}

func readData(src string, expectedSize int64) (raw []byte, err error) {
    raw, err = os.ReadFile(src)
    if err != nil {
        err = fmt.Errorf("Error trying to read file contents of %s", src, err)
        return
    }
    debug("read raw")
    readSize := int64(len(raw))
    if readSize != expectedSize {
        err = fmt.Errorf("Wrong number of bytes read (%d) from %s (%d expected)", readSize, src, expectedSize)
        return
    }
    return
}

func writeData(data []byte, target string) (err error) {
    err = os.WriteFile(target, data, 0664)
    if err != nil {
        err = fmt.Errorf("Error trying to write %s: %w", target, err)

        delErr := os.Remove(target)
        if delErr != nil {
            err = fmt.Errorf("%w. Also failed to delete partial file: %w", err, delErr)
        }
        debug(err)
        return
    }

    return
}

func compareSrcTgt(src, tgt fs.FileInfo) bool {
    if tgt.Size() == src.Size() && tgt.ModTime() == src.ModTime() {
        //These are likely the same file
        return true
    }
    debug(fmt.Sprintf("Target file size/mtime (%d/%v)does not match source file (%d/%v), overwriting", tgt.Size(), tgt.ModTime(), src.Size(), src.ModTime()))
    return false
}
