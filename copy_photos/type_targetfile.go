package main

import (
    "fmt"
    "io/fs"
    "os"
)

const dateDirFormat = "20060102"

//Values are defined in increasing amount of work needed
const (
    NoAction    =   iota    //No action will be taken
    NeedsCopy               //File will be copied
    NeedsVerify             //If target file is an incomplete copy of the original, it will be overwritten. otherwise, a conflict will be thrown
    Conflict                //Requires manual intervention
)

// A TargetFile struct stores the paths for the initial copy (archive) of a copied photo, and the non-hidden sortable file
// the sortable file will be a hardlink to the archive file.
// a more flexible struct would have the Sourcefile and info as named fields, and then a list of generic path/name pairs that would all be hardlinks from the initial copy
type TargetFile struct {
    SourceFile  string
    SourceInfo  fs.FileInfo

    TargetFile  string
    TargetStat  fs.FileInfo

    //depricate these
    ArchivePath, ArchiveFile string
    SortPath, SortFile string

    Links   []FileWithDirPath
    Action  int
}

type FileWithDirPath struct {
    Path    string
    File    string
}

func (t *TargetFile) Generate(rootPath string, f fs.DirEntry, linkDirs []string) ( e error ) {
    if len(linkDirs) < 1 {
        e = fmt.Errorf("No copy targets specified")
        return
    }
    t.SourceInfo, e = f.Info()
    if e != nil {
        return
    }

    t.Links = make([]FileWithDirPath, len(linkDirs))

    dateDir := t.SourceInfo.ModTime().Local().Format(dateDirFormat)
    for i, v := range linkDirs {
        t.Links[i].Path = rootPath + "/" +  v + "/" + dateDir
        t.Links[i].File = t.Links[i].Path + "/" + f.Name()
    }

    t.TargetFile = t.Links[0].File
    if len(t.Links) > 1 {
        t.Links = t.Links[1:]
    } else {
        t.Links = nil
    }

    t.ArchivePath = fmt.Sprintf("%s/%s/%s", rootPath, archiveDir, dateDir)
    t.ArchiveFile = fmt.Sprintf("%s/%s", t.ArchivePath, f.Name())

    t.SortPath    = fmt.Sprintf("%s/%s/%s", rootPath, sortDir, dateDir)
    t.SortFile    = fmt.Sprintf("%s/%s", t.SortPath, f.Name())

    return
}

func (target *TargetFile) CopyFromDisk(mp string) (err error) {
    debug("copying file from/to", target.SourceFile, target.TargetFile)
    err = target.MakePaths()
    if err != nil {
        err = fmt.Errorf("Failed creating target path: %v", err)
        return
    }

    var rawData []byte
    rawData, err = readData(mp+"/"+target.SourceFile, target.SourceInfo.Size())
    if err != nil {
        err = fmt.Errorf("Failed reading source file: %v", err)
        return
    }

    debug("correct number of bytes read:", target.SourceInfo.Size())

    if target.Action == NeedsVerify {
        var extantData []byte

        debug(fmt.Sprintf("size type is %T", target.TargetStat.Size()))
        debug("need to compare target file contents")
        extantData, err = readData( target.TargetFile, target.TargetStat.Size() )
        debug("havent checked error yet")
        if err != nil {
            err = fmt.Errorf("failed verifying target file: %v", err)
            return
        }
        debug("read extant data for target,size:", target.SourceInfo.Size())

        if ( compareByteSlices( extantData, rawData[ :len(extantData) ] ) ) {
            debug("tgt is an incomplete copy of src", target.TargetFile, target.SourceFile)
            //target file is bytewise identical to source, but incomplete
            target.Action = NeedsCopy
        } else {
            //target file has some byte value that the source does not
            debug("tgt contains data that src does not")
            target.Action = Conflict
            err = fmt.Errorf("conflict on copying file %s to %s", target.SourceFile, target.TargetFile)
            return
        }
    }

    if target.Action != NeedsCopy {
        err = fmt.Errorf("file %s somehow got to an uncreachable code path while processing copy tasks", target.SourceFile)
        return
    }
    err = writeData(rawData, target.TargetFile)
    if err != nil {
        return
    }
    debug("wrote file ", target.TargetFile)

    //force atime/mtime/ctime to be `mtime.Unix()` in syscall.Utime
    err = os.Chtimes(target.TargetFile, target.SourceInfo.ModTime(), target.SourceInfo.ModTime())
    if err != nil {
        err = fmt.Errorf("error setting atime/mtime for copied file (%s) to match source: %v", target.TargetFile, err)
        return
    }
    debug("set time")
    return
}

//This should create the dir paths for ArchivePath and SortPath
func (t *TargetFile) MakePaths() (error) {
    for _, entry := range t.Links {
        dirErr := os.MkdirAll(entry.Path, 0775)
        if dirErr != nil {
            return fmt.Errorf("Error creating directory: %s: %v", entry.Path)
        }
    }
    archDirErr := os.MkdirAll(t.ArchivePath, 0775)
    if archDirErr != nil {
        return fmt.Errorf("Error creating archive directory %s: %w", t.ArchivePath, archDirErr)
    }
    sortDirErr := os.MkdirAll(t.SortPath, 0775)
    if archDirErr != nil {
        return fmt.Errorf("Error creating sorting directory %s: %w", t.SortPath, sortDirErr)
    }
    return nil
}


func compareByteSlices(a, b []byte) (status bool) {
    if len(a) != len(b) {
        return
    }
    for k, v := range a {
        if v != b[k] {
            debug("differing value found at byte ", k)
            return
        }
    }
    status = true
    return
}

