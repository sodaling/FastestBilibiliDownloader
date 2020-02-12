package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
)

func VideoFetcher(referer string) FetchFun {
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

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
		}

		bodyReader := bufio.NewReader(resp.Body)

		e := determineEncoding(bodyReader)
		utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
		defer resp.Body.Close()
		return ioutil.ReadAll(utf8Reader)
	}
}
