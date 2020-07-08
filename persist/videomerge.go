package persist

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
	"strings"
	"sync"

	"github.com/go-cmd/cmd"
)

var _videoPageMap = make(map[int64]map[int64]int64)
var _contactFileName = "contact.txt"
var _videoOutputNameExt = ".mp4"

func VideoItemProcessor(wgOutside *sync.WaitGroup) (chan *engine.Item, error) {
	out := make(chan *engine.Item)
	go func() {
		defer wgOutside.Done()
		var wgInside sync.WaitGroup
		for item := range out {

			switch x := item.Payload.(type) {
			case *model.VideoAid:
				_videoPageMap[x.Aid] = make(map[int64]int64)
			case *model.VideoCid:
				_videoPageMap[x.ParAid.Aid][x.Page] = x.AllOrder

			case *model.Video:
				_videoPageMap[x.ParCid.ParAid.Aid][x.ParCid.Page] -= 1
				if _videoPageMap[x.ParCid.ParAid.Aid][x.ParCid.Page] == 0 {
					delete(_videoPageMap[x.ParCid.ParAid.Aid], x.ParCid.Page)
				}
				if len(_videoPageMap[x.ParCid.ParAid.Aid]) == 0 {
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

func mergeVideo(video *model.Video, wg *sync.WaitGroup) {
	defer wg.Done()
	var _videoOutputName = video.ParCid.ParAid.Title + _videoOutputNameExt
	aidDirPath := tool.GetAidFileDownloadDir(video.ParCid.ParAid.Aid, video.ParCid.ParAid.Title)
	contactTxtPath := filepath.Join(aidDirPath, _contactFileName)
	videoOutputPath := filepath.Join(tool.GetMp4Dir(), _videoOutputName)

	// merge cid
	for i := int64(1); i <= video.ParCid.ParAid.GetPage(); i++ {
		err := createMergeCidInfoTxt(aidDirPath, video.ParCid.Page, video.ParCid.AllOrder)
		if err != nil {
			log.Printf("some thing wrong while merging video %d", video.ParCid.ParAid.Aid)
			return
		}
		log.Println(video.ParCid.ParAid.Title, " download completed.Start merging videos now.")
		cidFilename := fmt.Sprintf("%d.flv", video.ParCid.Page)
		cidOutput := filepath.Join(aidDirPath, cidFilename)
		command := []string{"ffmpeg", "-f", "concat", "-safe", "0", "-i", contactTxtPath, "-c", "copy", cidOutput}
		findCmd := cmd.NewCmd(command[0], command[1:]...)
		<-findCmd.Start()
	}

	err := createMergeAidInfoTxt(aidDirPath, video.ParCid.ParAid.GetPage())
	if err != nil {
		log.Printf("some thing wrong while merging video %d", video.ParCid.ParAid.Aid)
		return
	}
	command := []string{"ffmpeg", "-f", "concat", "-safe", "0", "-i", contactTxtPath, "-c", "copy", videoOutputPath}
	findCmd := cmd.NewCmd(command[0], command[1:]...)
	<-findCmd.Start()
	log.Println("Video ", video.ParCid.ParAid.Title, " merge is complete.")
	removeTempFile(aidDirPath, _videoOutputName)
}

func createMergeAidInfoTxt(aidPath string, aidPage int64) error {
	videoCidPathTemp := "file '" + filepath.Join(aidPath, "%d.flv") + "'\n"
	txtPath := filepath.Join(aidPath, _contactFileName)

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

func createMergeCidInfoTxt(aidPath string, cidPage int64, cidAllOrder int64) error {
	videoCidPathTemp := "file '" + filepath.Join(aidPath, "%d_%d.flv") + "'\n"
	txtPath := filepath.Join(aidPath, _contactFileName)

	file, err := os.Create(txtPath)
	if err != nil {
		return err
	}
	defer file.Close()
	strBuilder := strings.Builder{}
	for i := int64(1); i <= cidAllOrder; i++ {
		strBuilder.WriteString(fmt.Sprintf(videoCidPathTemp, cidPage, i))
	}
	_, err = fmt.Fprintln(file, strBuilder.String())
	return err
}

func removeTempFile(dir, excludeFile string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == excludeFile {
			continue
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
