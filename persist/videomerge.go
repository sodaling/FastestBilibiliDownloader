package persist

import (
	"fmt"
	"log"
	"os"
	"path"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
	"strings"
	"sync"

	"github.com/go-cmd/cmd"
)

var _videoPageMap = make(map[int64]int64)
var _contactFileName = "contact.txt"
var _videoOutputName = "output.mp4"

func VideoItemProcessor(wgOutside *sync.WaitGroup) (chan *engine.Item, error) {
	out := make(chan *engine.Item)
	go func() {
		defer wgOutside.Done()
		var wgInside sync.WaitGroup
		for item := range out {

			switch x := item.Payload.(type) {
			case *model.VideoAid:
				_videoPageMap[x.Aid] = x.GetPage()
			case *model.VideoCid:
				_videoPageMap[x.ParAid.Aid] -= 1
				if _videoPageMap[x.ParAid.Aid] == 0 {
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

func mergeVideo(videoCiD *model.VideoCid, wg *sync.WaitGroup) {
	defer wg.Done()
	aidDirPath := tool.GetAidFileDownloadDir(videoCiD.ParAid.Aid, videoCiD.ParAid.Title)
	contactTxtPath := path.Join(aidDirPath, _contactFileName)
	videoOutputPath := path.Join(aidDirPath, _videoOutputName)

	err := createMergeInfoTxt(aidDirPath, videoCiD.ParAid.GetPage())
	if err != nil {
		log.Printf("some thing wrong while merging video %d", videoCiD.ParAid.Aid)
		return
	}
	log.Println(videoCiD.ParAid.Title, " download completed.Start merging videos now.")
	command := []string{"ffmpeg", "-f", "concat", "-safe", "0", "-i", contactTxtPath, "-c", "copy", videoOutputPath}
	findCmd := cmd.NewCmd(command[0], command[1:]...)
	<-findCmd.Start()
	log.Println("Video ", videoCiD.ParAid.Title, " merge is complete.")
}

func createMergeInfoTxt(aidPath string, aidPage int64) error {
	videoCidPathTemp := "file '" + path.Join(aidPath, "%d.flv") + "'\n"
	txtPath := path.Join(aidPath, _contactFileName)

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
