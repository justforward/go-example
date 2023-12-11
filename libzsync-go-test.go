package main

import (
	"fmt"
	"github.com/AppImageCrafters/libzsync-go"
	"os"
)

func main() {

	sync, _ := zsync.NewZSync("https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage.zsync")
	sync.RemoteFileUrl = "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage"

	output, _ := os.Create("/tmp/appimagetool-new-x86_64.AppImage")
	err := sync.Sync("/tmp/appimagetool-x86_64.AppImage", output)
	if err != nil {
		fmt.Println(err)
	}
}
