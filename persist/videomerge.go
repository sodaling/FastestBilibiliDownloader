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
var _contactFile2Name = "contactCid.txt"
var _videoOutputNameExt = ".mp4"
var _x_map = make(map[int64]map[int64]*model.VideoCid)


func VideoItemProcessor(wgOutside *sync.WaitGroup) (chan *engine.Item, error) {
	out := make(chan *engine.Item)

	go func() {
		defer wgOutside.Done()
		var wgInside sync.WaitGroup
		for item := range out {

			switch x := item.Payload.(type) {
			case *model.VideoAid:
				_videoPageMap[x.Aid] = make(map[int64]int64)
				_x_map[x.Aid] = make(map[int64]*model.VideoCid)

			case *model.VideoCid:
				_videoPageMap[x.ParAid.Aid][x.Page] = x.AllOrder
				_x_map[x.ParAid.Aid][x.Page] = x  // save as video.ParCid with type of *model.VideoCid

			case *model.Video:
				_videoPageMap[x.ParCid.ParAid.Aid][x.ParCid.Page] -= 1
				if _videoPageMap[x.ParCid.ParAid.Aid][x.ParCid.Page] == 0 {
					delete(_videoPageMap[x.ParCid.ParAid.Aid], x.ParCid.Page)
				}
				//x_map[x.ParCid.Page] = append(x_map[x.ParCid.Page], x)

				if len(_videoPageMap[x.ParCid.ParAid.Aid]) == 0 { //当整个列表是空的时执行，即当最后一个文件下载完
					//fmt.Println("_x_map:  ", _x_map)
					wgInside.Add(1)
					go mergeVideo_mod(_x_map[x.ParCid.ParAid.Aid],&wgInside)
				}

			default:
				panic(fmt.Sprintf("Unexpected type %T: %v", x, x))
			}

		}
		wgInside.Wait()
	}()
	return out, nil
}

func mergeVideo_mod(x_map map[int64]*model.VideoCid, wg *sync.WaitGroup) {
//func mergeVideo_mod(video_array []*model.Video, wg *sync.WaitGroup) {

	defer wg.Done()
	videoTmpParCid := x_map[int64(1)]   //[0]  //从一个子视频中获取视频总名称和aid （assume: 子视频的cid不同但aid和标题是一致的）

	aidDirPath := tool.GetAidFileDownloadDir(videoTmpParCid.ParAid.Aid, videoTmpParCid.ParAid.Title)
	contactCidTxtPath := filepath.Join(aidDirPath, _contactFile2Name)
    mp4DirPath := tool.GetMp4Dir(videoTmpParCid.ParAid.Title)

    log.Println(videoTmpParCid.ParAid.Title, " download completed. Start to merge videos now.")
    // log.Println("len video_array, ", len(_x_map), ", page number, ", videoTmp.ParCid.ParAid.GetPage())
	for i := int64(1); i <= videoTmpParCid.ParAid.GetPage(); i++ {
	    videoParCid := x_map[i]

	    // merge small parts in each cid
        err := createMergeCidInfoTxt(aidDirPath, videoParCid.Page, videoParCid.AllOrder)
        if err != nil {
            log.Printf("Something wrong while merging video %d.", videoParCid.ParAid.Aid)
            return
        }
        cidFilename := fmt.Sprintf("%d.flv", videoParCid.Page)
        cidOutput := filepath.Join(aidDirPath, cidFilename)
        command := []string{"ffmpeg", "-f", "concat", "-safe", "0", "-i", contactCidTxtPath, "-c", "copy", cidOutput}
        //log.Println(command)
        findCmd := cmd.NewCmd(command[0], command[1:]...)
        <-findCmd.Start()

        //convert from flv to mp4
        mp4Filename := videoParCid.Part + ".mp4"
        mp4Output := filepath.Join(mp4DirPath, mp4Filename)
        log.Println(videoParCid.ParAid.Title+"/"+mp4Filename, " merge completed. Start to convert to mp4.")
        command_new := []string{"ffmpeg", "-i", cidOutput, mp4Output}
        //log.Println(command_new)
        findCmd_new := cmd.NewCmd(command_new[0], command_new[1:]...)
        <-findCmd_new.Start()
        log.Println("Video ", videoParCid.ParAid.Title+"/"+mp4Filename, " merge and conversion is complete.")

        // free the map
        delete(_x_map[videoParCid.ParAid.Aid], i)
	}

    // can comment out the line below for debugging
	removeTempFile(aidDirPath, _contactFile2Name)
}

func createMergeCidInfoTxt(aidPath string, cidPage int64, cidAllOrder int64) error {
	videoCidPathTemp := "file '" + filepath.Join(aidPath, "%d_%d.flv") + "'\n"
	txtPath := filepath.Join(aidPath, _contactFile2Name)

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
    log.Println("Merge is completed, start to remove all temporary files.")

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
