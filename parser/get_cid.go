package parser

import (
	"crypto/md5"
	"fmt"
	"github.com/tidwall/gjson"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/model"
	"simple-golang-crawler/tool"
	"strconv"
)

var videoDownloadApi = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var entropy = "rbMCKn@KuamXWlPMoJGsKcbiJKUfkPF_8dABscJntvqhRSETg"
var paramsTemp = "appkey=%s&cid=%s&otype=json&qn=%s&quality=%s&type="
var playApiTemp = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var quailty = "80"

func GenGetAidChildrenParseFun(videoAid *model.VideoAidInfo) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		var retParseResult engine.ParseResult
		data := gjson.GetBytes(contents, "data").Array()
		appkey, sec := tool.GetAppkey(entropy)

		for _, i := range data {
			cid := i.Get("cid").Int()
			page := i.Get("page").Int()
			videoCid := model.NewVideoCidInfo(cid, videoAid, page)
			videoAid.AddCid(videoCid)
			cidStr := strconv.FormatInt(videoCid.Cid, 10)

			params := fmt.Sprintf(paramsTemp, appkey, cidStr, quailty, quailty)
			chksum := fmt.Sprintf("%x", md5.Sum([]byte(params+sec)))
			urlApi := fmt.Sprintf(playApiTemp, params, chksum)
			req := engine.NewRequest(urlApi, GenVideoDownloadParseFun(videoCid), fetcher.DefaultFetcher)
			retParseResult.Requests = append(retParseResult.Requests, req)
		}

		return retParseResult
	}
}

func GetRequestByAid(aid int64) (*engine.Request, *model.VideoAidInfo) {
	reqUrl := fmt.Sprintf(getCidUrlTemp, aid)
	videoAid := model.NewVideoAidInfo(aid, fmt.Sprintf("%d", aid))
	reqParseFunction := GenGetAidChildrenParseFun(videoAid)
	req := engine.NewRequest(reqUrl, reqParseFunction, fetcher.DefaultFetcher)
	return req, videoAid
}
