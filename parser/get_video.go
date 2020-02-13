package parser

import (
	"github.com/tidwall/gjson"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/model"
)

func GenVideoDownloadParseFun(video *model.VideoInfo) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		retParseResult := engine.ParseResult{}

		durlArray := gjson.GetBytes(contents, "durl").Array()
		for _, i := range durlArray {
			videoUrl := i.Get("url").String()
			req := engine.NewRequest(videoUrl, NilParseFun, fetcher.GenVideoFetcher(video))
			retParseResult.Requests = append(retParseResult.Requests, req)
		}
		return retParseResult
	}
}

func NilParseFun(contents []byte, url string) engine.ParseResult {
	return engine.ParseResult{}
}
