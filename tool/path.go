package tool

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetAidFileDownloadDir(aid int64, title string) string {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// remove special characters
	//title = TitleEdit(title)

	fullDirPath := filepath.Join(curDir, "download", fmt.Sprintf("%d_%s", aid, title))
	err = os.MkdirAll(fullDirPath, 0777)
	if err != nil {
		panic(err)
	}
	return fullDirPath
}

func GetMp4Dir(title string) string {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

    // remove special characters
	//title = TitleEdit(title)

	fullDirPath := filepath.Join(curDir, "output", title)
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

func TitleEdit(title string) string {  // will be used when save the title or the part
    // remove special symbol
	title = strings.Replace(title, ":" , "", -1)
	title = strings.Replace(title, "\\",  "", -1)
	title = strings.Replace(title, "/" , "", -1)
	title = strings.Replace(title, "*" , "", -1)
	title = strings.Replace(title, "?" , "", -1)
	title = strings.Replace(title, "\"" , "", -1)
	title = strings.Replace(title, "<" , "", -1)
	title = strings.Replace(title, ">" , "", -1)
	title = strings.Replace(title, "|" , "", -1)
	title = strings.Replace(title, "." , "", -1)

	return title
}