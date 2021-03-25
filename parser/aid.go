package parser

import (
	"fmt"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
	"math"
	"github.com/tidwall/gjson"
)

var _getAidUrlTemp = "https://api.bilibili.com/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"
var _getCidUrlTemp = "https://api.bilibili.com/x/web-interface/view?aid=%d"
//var _getCidUrlTemp = "https://api.bilibili.com/x/player/pagelist?aid=%d"

var table string = "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"
var s = [6]int{11, 10, 3, 8, 4, 6}
var xor = 177451812
var add = 8728348608
var tr map[string]int


func UpSpaceParseFun(contents []byte, url string) engine.ParseResult {

	var retParseResult engine.ParseResult
	value := gjson.GetManyBytes(contents, "data.list.vlist", "data.page")

	var upid int64
	retParseResult.Requests, upid = getAidDetailReqList(value[0])
	retParseResult.Requests = append(retParseResult.Requests, getNewBilibiliUpSpaceReqList(value[1], upid)...)

	return retParseResult

}

func getAidDetailReqList(pageInfo gjson.Result) ([]*engine.Request, int64) {

	var retRequests []*engine.Request
	var upid int64
	for _, i := range pageInfo.Array() {
		aid := i.Get("aid").Int()
		upid = i.Get("mid").Int()
		title := i.Get("title").String()
		title = tool.TitleEdit(title)  // remove special characters
		reqUrl := fmt.Sprintf(_getCidUrlTemp, aid)
		videoAid := model.NewVideoAidInfo(aid, title)
		reqParseFunction := GenGetAidChildrenParseFun(videoAid) //子视频
		req := engine.NewRequest(reqUrl, reqParseFunction, fetcher.DefaultFetcher)
		retRequests = append(retRequests, req)
	}
	return retRequests, upid
}

// 访问up主的时候 需要翻页
func getNewBilibiliUpSpaceReqList(pageInfo gjson.Result, upid int64) []*engine.Request {

	var retRequests []*engine.Request

	count := pageInfo.Get("count").Int()
	pn := pageInfo.Get("pn").Int()
	ps := pageInfo.Get("ps").Int()
	var extraPage int64
	if count%ps > 0 {
		extraPage = 1
	}
	totalPage := count/ps + extraPage
	for i := int64(1); i <= totalPage; i++ {
		if i == pn {
			continue
		}
		reqUrl := fmt.Sprintf(_getAidUrlTemp, upid, i)
		req := engine.NewRequest(reqUrl, UpSpaceParseFun, fetcher.DefaultFetcher)
		retRequests = append(retRequests, req)
	}
	return retRequests
}

func GetRequestByUpId(upid int64) *engine.Request {

	reqUrl := fmt.Sprintf(_getAidUrlTemp, upid, 1)
	return engine.NewRequest(reqUrl, UpSpaceParseFun, fetcher.DefaultFetcher)
}

// source code: https://blog.csdn.net/dotastar00/article/details/108805779
func Bv2av(x string) int64 {
    tr = make(map[string]int)
    for i:=0; i<58; i++ {
        tr[string(table[i])] = i
    }
    r := 0
    for i:=0; i<6; i++ {
        r += tr[string(x[s[i]])] * int(math.Pow(float64(58), float64(i)))
    }
    return int64((r - add) ^ xor)
}
