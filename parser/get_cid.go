package parser

import (
	"crypto/md5"
	"fmt"
	"github.com/tidwall/gjson"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/model"
	"strconv"
	"strings"
)

var videoDownloadApi = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var entropy = "rbMCKn@KuamXWlPMoJGsKcbiJKUfkPF_8dABscJntvqhRSETg"
var paramsTemp = "appkey=%s&cid=%s&otype=json&qn=%s&quality=%s&type="
var playApiTemp = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var quailty = "80"

func GenGetAidInfoParseFun(aid int64, title string) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		var retParseResult engine.ParseResult
		data := gjson.GetBytes(contents, "data").Array()
		appkey, sec := getAppkey(entropy)

		for _, i := range data {
			videoInfo := &model.VideoInfo{}
			videoInfo.Aid = aid
			videoInfo.Title = title
			videoInfo.Cid = i.Get("cid").Int()
			videoInfo.Page = i.Get("page").Int()
			cidStr := strconv.FormatInt(videoInfo.Cid, 10)

			params := fmt.Sprintf(paramsTemp, appkey, cidStr, quailty, quailty)
			chksum := fmt.Sprintf("%x", md5.Sum([]byte(params+sec)))
			urlApi := fmt.Sprintf(playApiTemp, params, chksum)
			req := engine.NewRequest(urlApi, GenVideoDownloadParseFun(videoInfo), fetcher.DefaultFetcher)
			retParseResult.Requests = append(retParseResult.Requests, req)

			//filename := fmt.Sprintf("%d:%s", aid, title)
			//err := os.MkdirAll(path.Join("./download/", filename), 0777)
			//if err != nil {
			//	log.Println(err)
			//}
		}

		return retParseResult
	}
}

func getAppkey(entropy string) (appkey, sec string) {
	revEntropy := reverseRunes([]rune(entropy))
	for i := range revEntropy {
		revEntropy[i] = revEntropy[i] + 2
	}
	ret := strings.Split(string(revEntropy), ":")

	return ret[0], ret[1]
}

func reverseRunes(runes []rune) []rune {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return runes
}
