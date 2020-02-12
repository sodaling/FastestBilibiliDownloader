package parser

import (
	"fmt"
	"github.com/tidwall/gjson"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
)

var getAidUrl = "https://api.bilibili.com/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"
var getCidUrl = "https://api.bilibili.com/x/player/pagelist?aid=%d"

func UpSpaceParseFun(contents []byte, url string) engine.ParseResult {
	var retParseResult engine.ParseResult
	value := gjson.GetManyBytes(contents, "data.list.vlist", "data.page")

	retParseResult.Requests = getAidDetailReqList(value[0])
	retParseResult.Requests = append(retParseResult.Requests, getNewBilibiliUpSpaceReqList(value[1])...)

	return retParseResult

}
func getAidDetailReqList(pageInfo gjson.Result) []*engine.Request {
	var retRequests []*engine.Request
	for _, i := range pageInfo.Array() {
		aid := i.Get("aid").Int()
		reqUrl := fmt.Sprintf(getCidUrl, aid)
		reqParseFunction := GenGetAidInfoParseFun(aid)
		req := engine.NewRequest(reqUrl, reqParseFunction, fetcher.DefaultFetch)
		retRequests = append(retRequests, req)
	}
	return retRequests
}

func getNewBilibiliUpSpaceReqList(pageInfo gjson.Result) []*engine.Request {
	var retRequests []*engine.Request

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
		reqUrl := fmt.Sprintf(getAidUrl, 585267, i)
		reqParseFunction := UpSpaceParseFun
		req := engine.NewRequest(reqUrl, reqParseFunction, fetcher.DefaultFetch)
		retRequests = append(retRequests, req)
	}
	return retRequests
}
