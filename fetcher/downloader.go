package fetcher

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
	"time"
)

var _startUrlTem = "https://api.bilibili.com/x/web-interface/view?aid=%d"

func GenVideoFetcher(video *model.Video) FetchFun {
	referer := fmt.Sprintf(_startUrlTem, video.ParCid.ParAid.Aid)
	for i := int64(1); i <= video.ParCid.Page; i++ {
		referer += fmt.Sprintf("/?p=%d", i)
	}

	return func(url string) (bytes []byte, err error) {
		<-_rateLimiter.C
		client := http.Client{CheckRedirect: genCheckRedirectfun(referer)}

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalln(url, err)
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
			log.Fatalf("下载 %d 时出错, 错误信息：%s", video.ParCid.Cid, err)
			return nil, err
		}

		if resp.StatusCode != http.StatusPartialContent {
			log.Fatalf("下载 %d 时出错, 错误码：%d", video.ParCid.Cid, resp.StatusCode)
			return nil, fmt.Errorf("错误码： %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		aidPath := tool.GetAidFileDownloadDir(video.ParCid.ParAid.Aid, video.ParCid.ParAid.Title)
		filename := fmt.Sprintf("%d_%d.flv", video.ParCid.Page, video.Order)
		file, err := os.Create(filepath.Join(aidPath, filename))
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		defer file.Close()

		log.Println("正在下载：" + video.ParCid.ParAid.Title + "\\" + filename)
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Printf("下载失败 aid: %d, cid: %d, title: %s, part: %s",
				video.ParCid.ParAid.Aid, video.ParCid.Cid, video.ParCid.ParAid.Title, video.ParCid.Part)
			log.Println("错误信息：", err)

			// request again
			go requestLater(file, resp, video)
			return nil, err
		}
		log.Println("下载完成：" + video.ParCid.ParAid.Title + "\\" + filename)

		return nil, nil
	}
}

func genCheckRedirectfun(referer string) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Referer", referer)
		return nil
	}
}

func requestLater(file *os.File, resp *http.Response, video *model.Video) error {

	log.Println("连接失败，30秒后重试 (Unable to open the file due to the remote host, request in 30 seconds)")
	time.Sleep(time.Second * 30)

	_, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("下载失败 aid: %d, cid: %d, title: %s, part: %s again",
			video.ParCid.ParAid.Aid, video.ParCid.Cid, video.ParCid.ParAid.Title, video.ParCid.Part)
	}
	return err
}
