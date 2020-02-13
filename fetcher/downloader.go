package fetcher

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
)

var startUrlTem = "https://api.bilibili.com/x/web-interface/view?aid=%d/?p=%d"

func GenVideoFetcher(videoInfo *model.VideoInfo) FetchFun {
	referer := fmt.Sprintf(startUrlTem, videoInfo.Aid, videoInfo.Page)

	return func(url string) (bytes []byte, err error) {
		<-rateLimiter
		client := http.DefaultClient
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(url)
			return nil, err
		}
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:56.0) Gecko/20100101 Firefox/56.0")
		request.Header.Set("Accept", "*/*")
		request.Header.Set("Accept-Language", "en-US,en;q=0.5")
		request.Header.Set("Accept-Encoding", "gzip, deflate, br")
		request.Header.Set("Range", "bytes=0-")
		request.Header.Set("Referer", referer)
		request.Header.Set("Origin", "https://www.bilibili.com")
		request.Header.Set("Connection", "keep-alive")

		resp, err := client.Do(request)
		if err != nil {
			fmt.Println(url)
			return nil, err
		}

		if resp.StatusCode != http.StatusPartialContent {
			fmt.Println(resp.StatusCode)
			return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		aidPath := tool.GetAidFileDownloadDir(videoInfo.Aid,videoInfo.Title)
		filename := fmt.Sprintf("%d.flv", videoInfo.Cid)
		file, err := os.Create(path.Join(aidPath, filename))
		if err != nil {
			os.Exit(1)
		}
		defer file.Close()

		log.Println(filename + " is downloading.")
		io.Copy(file, resp.Body)
		log.Println(filename + " has finished.")

		return nil, nil
	}
}
