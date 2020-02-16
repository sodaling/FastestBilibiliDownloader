package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var _rateLimiter = time.NewTicker(100 * time.Microsecond)

type FetchFun func(url string) ([]byte, error)

func DefaultFetcher(url string) ([]byte, error) {
	<-_rateLimiter.C
	client := http.DefaultClient
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("fetch err while request :%s,and the err is %s", url, err)
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("fetch err while request :%s,and the err is %s", url, err)
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

func determineEncoding(reader *bufio.Reader) encoding.Encoding {
	bytes, err := reader.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
