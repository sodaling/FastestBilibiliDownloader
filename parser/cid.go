package parser

import (
	"crypto/md5"
	"fmt"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
	"strconv"

	"github.com/tidwall/gjson"
)

var _entropy = "rbMCKn@KuamXWlPMoJGsKcbiJKUfkPF_8dABscJntvqhRSETg"
var _paramsTemp = "appkey=%s&cid=%s&otype=json&qn=%s&quality=%s&type="
var _playApiTemp = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var _quality = "80"

func GenGetAidChildrenParseFun(videoAid *model.VideoAid) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		var retParseResult engine.ParseResult
		data := gjson.GetBytes(contents, "data").Array()
		appKey, sec := tool.GetAppKey(_entropy)

		var videoTotalPage int64
		for _, i := range data {
			cid := i.Get("cid").Int()
			page := i.Get("page").Int()
			videoCid := model.NewVideoCidInfo(cid, videoAid, page)
			videoTotalPage += 1
			cidStr := strconv.FormatInt(videoCid.Cid, 10)

			params := fmt.Sprintf(_paramsTemp, appKey, cidStr, _quality, _quality)
			chksum := fmt.Sprintf("%x", md5.Sum([]byte(params+sec)))
			urlApi := fmt.Sprintf(_playApiTemp, params, chksum)
			req := engine.NewRequest(urlApi, GenVideoDownloadParseFun(videoCid), fetcher.DefaultFetcher)
			retParseResult.Requests = append(retParseResult.Requests, req)
		}

		videoAid.SetPage(videoTotalPage)
		item := engine.NewItem(videoAid)
		retParseResult.Items = append(retParseResult.Items, item)

		return retParseResult
	}
}

func GetRequestByAid(aid int64) *engine.Request {
	reqUrl := fmt.Sprintf(_getCidUrlTemp, aid)
	videoAid := model.NewVideoAidInfo(aid, fmt.Sprintf("%d", aid))
	reqParseFunction := GenGetAidChildrenParseFun(videoAid)
	req := engine.NewRequest(reqUrl, reqParseFunction, fetcher.DefaultFetcher)
	return req
}
