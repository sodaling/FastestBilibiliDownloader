package parser

import (
	"fmt"
	"github.com/tidwall/gjson"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
)

var count int

var startUrlTem = "https://api.bilibili.com/x/web-interface/view?aid=%d/?p=%d"

func GenVideoParseFun(video *model.VideoInfo) engine.ParseFunc {
	return func(contents []byte, url string) engine.ParseResult {
		//referer := fmt.Sprintf(startUrlTem, video.Aid, video.Page)

		value := gjson.GetBytes(contents, "durl")
		fmt.Println(value.String())
		return engine.ParseResult{}
	}
}
