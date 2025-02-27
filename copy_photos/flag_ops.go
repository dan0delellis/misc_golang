package main
import (
    flags "github.com/jessevdk/go-flags"
)

//a A B c C D e E F g G H i I j K L m M n N o O P q Q R S T u U V w W x X y Y z Z
type Opts struct {
    TempDir     string   `short:"t" long:"tempdir" description:"root dir, relative to CWD, to create mountpoint for the discovered block device" default:"temp"`
    TempPrefix  string   `short:"p" long:"tempname" description:"name of directory to create as a temp mountpoint for discovered block device" default:"camera_sd"`

    BlkidCache  string   `short:"b" long:"blkidtab" description:"path, relative to CWD, for the contents of blkid" default:"/run/blkid/blkid.tab"`
    FsTypes     []string `short:"T" long:"fstypes" description:"filesystem type to examine; call multiple times for multiple types" default:"exfat" default:"fat32"`
    FsLabels    []string `short:"l" long:"fslabels" description:"filesystem labels to examine, not case-sensitive" default:"nikon"`

    RootPath    string   `short:"r" long:"rootdir" description:"root of directory path, relative to CWD, to copy files into." default:"photos"`
    TargetDirs  []string `short:"d" long:"dirname" description:"path, relative to 'rootdir', to create date-formatted directory names. call multiple times to create hardlink clones of the initial directory" default:"."`
    DirFormat   string   `short:"f" long:"dirformat" description:"go-date formatted string for target directory names to move files to" default:"20060102"`

    MinSize     int      `short:"s" long:"minsize" description:"minimum filesize, in bytes, to be included in the file copy list" default:"1000"`
    Verbose     bool     `short:"v" long:"verbose" description:"print ops while processing"`
}

func parseFlags() (o Opts, err error) {
    _, err = flags.Parse(&o)
    return
}
