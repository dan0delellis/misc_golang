package main
import (
    flags "github.com/jessevdk/go-flags"
)

//a A B c C D e E F G H i I j K L m M n o O P q Q R S T U V w W x X y Y z Z
type Opts struct {
    TempDir     string   `short:"t" long:"tempdir" description:"root dir, relative to CWD, to create mountpoint for the discovered block device" default:"temp"`
    TempPrefix  string   `short:"p" long:"tempname" description:"name of directory to create as a temp mountpoint for discovered block device" default:"camera_sd"`

    DevIDs      []string `short:"D" long:"dev" description:"multiple instances allowed. full device path (ex: '/dev/sdb1') to scan. Overrides fstypes,fslabels options entirely. Will not override nikonfile or nikon flag"`
    NikonFile   bool     `short:"N" long:"nikon" description:"look for and validate a 512byte file named NIKON001.DSC in the root of the drive"`
    NikonFilePath string `short:"n" long:"nikonfile" description:"look for and validate a 512byte file with specified name in the root of the drive. Overrides default path assumed by nikon flag"`

    BlkidCache  string   `short:"b" long:"blkidtab" description:"path, relative to CWD, for the contents of blkid" default:"/run/blkid/blkid.tab"`
    FsTypes     []string `short:"T" long:"fstypes" description:"multiple instances allowed. filesystem type to examine" default:"exfat" default:"fat32"`
    FsLabels    []string `short:"l" long:"fslabels" description:"multiple instances allowed. filesystem labels to examine, not case-sensitive. Call multiple times for multiple values" default:"nikon"`

    RootPath    string   `short:"r" long:"rootdir" description:"root of directory path, relative to CWD, to copy files into." default:"photos"`
    TargetDirs  []string `short:"d" long:"dirname" description:"multiple instances allowed. path, relative to 'rootdir', to create date-formatted directory names. call multiple times to create hardlink clones of the initial directory" default:"."`
    DirFormat   string   `short:"f" long:"dirformat" description:"go-date formatted string for target directory names to move files to" default:"20060102"`

    UserID      int      `short:"u" long:"uid"  description:"numeric user ID to set ownership of files/dirs. setting both user and group IDs will skip checks for rootdir existing" default:"-1"`
    GroupID     int      `short:"g" long:"gid"  description:"numeric group ID to set ownership of files/dirs. setting both user and group IDs will skip checks for rootdir existing" default:"-1"`

    MinSize     int      `short:"s" long:"minsize" description:"minimum filesize, in bytes, to be included in the file copy list" default:"1000"`
    Verbose     bool     `short:"v" long:"verbose" description:"print ops while processing"`
}

func parseFlags() (o Opts, err error) {
    _, err = flags.Parse(&o)
    return
}
