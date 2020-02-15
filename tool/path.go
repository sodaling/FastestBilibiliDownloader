package tool

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func GetAidFileDownloadDir(aid int64, title string) string {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fullDirPath := path.Join(curDir, "download", fmt.Sprintf("%d:%s", aid, title))
	err = os.MkdirAll(fullDirPath, 0777)
	if err != nil {
		panic(err)
	}
	return fullDirPath
}

func FileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func CheckFfmegStatus() bool {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return false
	} else {
		return true
	}
}
