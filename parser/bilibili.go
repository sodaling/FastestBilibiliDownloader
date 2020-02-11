package parser

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
)

var upSpaceVideoUrlTemp = "https://api.bilibili.com/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"

func BilibiliParseFun(content []byte, url string) engine.ParseResult {
	var retParseResult engine.ParseResult
	value := gjson.GetManyBytes(content, "data.list.vlist", "data.page")

	for _, i := range value[0].Array() {
		videoModel := &model.Video{}
		videoByte := []byte(i.String())
		err := json.Unmarshal(videoByte, videoModel)
		if err != nil {
			continue
		}

		var videoItem engine.Item
		videoItem.Url = url
		videoItem.Payload = videoModel
		retParseResult.Items = append(retParseResult.Items, videoItem)
	}
	retParseResult.Requests = getNewBilibiliUpSpaceReq(value[1])

	return retParseResult

}

func getNewBilibiliUpSpaceReq(pageInfo gjson.Result) []engine.Request {
	var retRequests []engine.Request

	count := pageInfo.Get("count").Int()
	pn := pageInfo.Get("pn").Int()
	ps := pageInfo.Get("ps").Int()
	var extraPage int64
	if count%ps > 0 {
		extraPage = 1
	}
	totalPage := count/ps + extraPage
	for i := int64(1); i < totalPage; i++ {
		if i == pn {
			continue
		}
		var req engine.Request
		req.Url = fmt.Sprintf(upSpaceVideoUrlTemp, 585267, i)
		req.ParseFunction = BilibiliParseFun
		retRequests = append(retRequests, req)
	}
	return retRequests
}
