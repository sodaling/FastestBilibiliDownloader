package persist

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"os"
	"path"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
	"strings"
	"sync"
)

var videPageMap = make(map[int64]int64)
var contactFileName = "contact.txt"
var videoOutputName = "output.mp4"

func VideoItemProcessor(wgOutside *sync.WaitGroup) (chan *engine.Item, error) {
	out := make(chan *engine.Item)
	go func() {
		defer wgOutside.Done()
		var wgInside sync.WaitGroup
		for item := range out {

			switch x := item.Payload.(type) {
			case *model.VideoAidInfo:
				videPageMap[x.Aid] = x.GetPage()
			case *model.VideoCidInfo:
				videPageMap[x.ParAid.Aid] -= 1
				if videPageMap[x.ParAid.Aid] == 0 {
					wgInside.Add(1)
					go mergeVideo(x, &wgInside)
				}
			default:
				panic(fmt.Sprintf("unexpected type %T: %v", x, x))
			}

		}
		wgInside.Wait()
	}()
	return out, nil
}

func mergeVideo(videoCiD *model.VideoCidInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	aidDirPath := tool.GetAidFileDownloadDir(videoCiD.ParAid.Aid, videoCiD.ParAid.Title)
	contactTxtPath := path.Join(aidDirPath, contactFileName)
	videoOutputPath := path.Join(aidDirPath, videoOutputName)

	createMergeInfoTxt(aidDirPath, videoCiD.ParAid.GetPage())
	fmt.Println(videoCiD.ParAid.Title, " download completed.Start merging videos now.")
	strs := []string{"ffmpeg", "-f", "concat", "-safe", "0", "-i", contactTxtPath, "-c", "copy", videoOutputPath}
	findCmd := cmd.NewCmd(strs[0], strs[1:]...)
	<-findCmd.Start()
	fmt.Println("Video ", videoCiD.ParAid.Title, " merge is complete.")
}

func createMergeInfoTxt(aidPath string, aidPage int64) error {
	videoCidPathTemp := "file '" + path.Join(aidPath, "%d.flv") + "'\n"
	txtPath := path.Join(aidPath, contactFileName)

	file, err := os.Create(txtPath)
	if err != nil {
		return err
	}
	defer file.Close()
	strBuilder := strings.Builder{}
	for i := int64(1); i <= aidPage; i++ {
		strBuilder.WriteString(fmt.Sprintf(videoCidPathTemp, i))
	}
	_, err = fmt.Fprintln(file, strBuilder.String())
	return err
}
