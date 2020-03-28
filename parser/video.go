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

		videoCid.AllOrder = int64(len(durlSlice))
		item := engine.NewItem(videoCid)
		retParseResult.Items = append(retParseResult.Items, item)

		for _, i := range durlSlice {
			video := &model.Video{Order: i.Get("order").Int(), ParCid: videoCid}
			videoUrl := i.Get("url").String()
			req := engine.NewRequest(videoUrl, recordCidParseFun(video), fetcher.GenVideoFetcher(video))
			retParseResult.Requests = append(retParseResult.Requests, req)
		}
		return retParseResult
	}
}

func recordCidParseFun(Video *model.Video) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		var retResult engine.ParseResult
		item := engine.NewItem(Video)
		retResult.Items = append(retResult.Items, item)
		return retResult
	}
}
