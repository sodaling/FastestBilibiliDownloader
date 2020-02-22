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

var _startUrlTem = "https://api.bilibili.com/x/web-interface/view?aid=%d"

func GenVideoFetcher(videoCid *model.VideoCid) FetchFun {
	referer := fmt.Sprintf(_startUrlTem, videoCid.ParAid.Aid)
	for i := int64(1); i <= videoCid.Page; i++ {
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
			log.Fatalf("Fail to download the video %d,err is %s", videoCid.Cid, err)
			return nil, err
		}

		if resp.StatusCode != http.StatusPartialContent {
			log.Fatalf("Fail to download the video %d,status code is %d", videoCid.Cid, resp.StatusCode)
			return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		aidPath := tool.GetAidFileDownloadDir(videoCid.ParAid.Aid, videoCid.ParAid.Title)
		filename := fmt.Sprintf("%d.flv", videoCid.Page)
		file, err := os.Create(path.Join(aidPath, filename))
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		defer file.Close()

		log.Println(videoCid.ParAid.Title + ":" + filename + " is downloading.")
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Printf("Failed to download video %d", videoCid.Cid)
			return nil, err
		}
		log.Println(videoCid.ParAid.Title + ":" + filename + " has finished.")

		return nil, nil
	}
}

func genCheckRedirectfun(referer string) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Referer", referer)
		return nil
	}
}
