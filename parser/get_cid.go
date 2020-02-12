package parser

import (
	"fmt"
	"simple-golang-crawler/engine"
	"strings"
)

var videoDownloadApi = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var entropy = "rbMCKn@KuamXWlPMoJGsKcbiJKUfkPF_8dABscJntvqhRSETg"
var paramsTemp = "appkey=%s&cid=%s&otype=json&qn=%s&quality=%s&type="
var playApiTemp = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var quailty = "80"

func GenGetAidInfoParseFun(aid int64) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		var retParseResult engine.ParseResult
		fmt.Println(aid)
		//data := gjson.GetBytes(contents, "data").Array()
		//appkey, sec := getAppkey(entropy)
		//
		//for _, i := range data {
		//	var req engine.
		//	videoInfo := &VideoInfo{}
		//	videoInfo.Aid = aid
		//	videoInfo.Cid = i.Get("cid").Int()
		//	cidStr := strconv.FormatInt(videoInfo.Cid, 10)
		//
		//	params := fmt.Sprintf(paramsTemp, appkey, cidStr, quailty, quailty)
		//	chksum := fmt.Sprintf("%x", md5.Sum([]byte(params+cidStr)))
		//	urlApi := fmt.Sprintf(playApiTemp, params, chksum)
		//}

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
