package main

import (
    "fmt"
    "io/fs"
    "os"
)

// A TargetFile struct stores the paths for the initial copy (archive) of a copied photo, and the non-hidden sortable file
// the sortable file will be a hardlink to the archive file.
type TargetFile struct {
    ArchivePath, ArchiveFile string
    SortPath, SortFile string
    Info    fs.FileInfo
}

func (t *TargetFile) Generate(f fs.DirEntry) ( e error ){
    t.Info, e = f.Info()
    if e != nil {
        return
    }
    dateDir := t.Info.ModTime().Local().Format(dateDirFormat) + "/"

    t.ArchivePath = fmt.Sprintf("%s/%s/%s", photosDir, archiveDir, dateDir)
    t.ArchiveFile = fmt.Sprintf("%s/%s", t.ArchivePath, f.Name())

    t.SortPath    = fmt.Sprintf("%s/%s/%s", photosDir, sortDir, dateDir)
    t.SortFile    = fmt.Sprintf("%s/%s", t.SortPath, f.Name())

    return
}

//This should create the dir paths for ArchivePath and SortPath
func (t *TargetFile) MakePaths() (err error) {
    archDirErr := os.MkdirAll(t.ArchivePath, 0775)
    if archDirErr != nil {
        err = fmt.Errorf("Error creating archive directory %s: %w", t.ArchivePath, archDirErr)
        return
    }
    sortDirErr := os.MkdirAll(t.SortPath, 0775)
    if archDirErr != nil {
        err = fmt.Errorf("Error creating sorting directory %s: %w", t.SortPath, sortDirErr)
        return
    }
    return
}
