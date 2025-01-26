package main
const blkidCache = "/run/blkid/blkid.tab"
const fsTypeKey = "type"
const mode = "ro"

const fsLabelKey = "label"
const partLabel = "nikon"

const tempDir = "temp"
const tempPrefix = "camera_sd"

const dateDirFormat = "20060102"
const sortDir = "mnt/media/photos/sort/"
const archiveDir = "mnt/media/photos/.copy_target/"
var archiveUid int
var archiveGid int

func fsTypes() ([]string) {
    return []string{"exfat", "fat32"}
}
