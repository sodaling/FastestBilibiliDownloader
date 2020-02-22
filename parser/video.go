package parser

import (
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/model"

	"github.com/tidwall/gjson"
)

func GenVideoDownloadParseFun(videoCid *model.VideoCid) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		retParseResult := engine.ParseResult{}

		durlSlice := gjson.GetBytes(contents, "durl").Array()
		for _, i := range durlSlice {
			videoUrl := i.Get("url").String()
			req := engine.NewRequest(videoUrl, recordCidParseFun(videoCid), fetcher.GenVideoFetcher(videoCid))
			retParseResult.Requests = append(retParseResult.Requests, req)
		}
		return retParseResult
	}
}

func recordCidParseFun(cidVideo *model.VideoCid) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		var retResult engine.ParseResult

		item := engine.NewItem(cidVideo)
		retResult.Items = append(retResult.Items, item)
		return retResult
	}
}
