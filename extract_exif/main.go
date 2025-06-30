package main

import (
	"fmt"
	exif "github.com/dsoprea/go-exif"
)

const testFile = "/mnt/media/photos/.archive/20250606/DSC_1377.NEF"

func main() {
	data, _ := exif.SearchFileAndExtractExif(testFile)
	im := exif.NewIfdMapping()

	_ = exif.LoadStandardIfds(im)

	ti := exif.NewTagIndex()

	eh, index, _ := exif.Collect(im, ti, data)
	fmt.Println(eh)

	res, _ := index.RootIfd.FindTagWithName("PreviewImageLength")
	ite := res[0]
	val, _ := index.RootIfd.TagValue(ite)
	fmt.Println(val.(string))


}
