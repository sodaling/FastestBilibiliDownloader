package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/parser"
	"simple-golang-crawler/persist"
	"simple-golang-crawler/scheduler"
	"strconv"
	"sync"
)

func main() {
	itemProcessFun := persist.GetItemProcessFun()
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	itemChan, err := itemProcessFun(&wg)
	if err != nil {
		panic(err)
	}

	var urlInput string

	var idType = "else"
	var aid int64
	var upid int64
	var bvid string

	var params []string

	var req *engine.Request

	fmt.Println("欢迎使用B站视频下载器 v1.0.1")
	fmt.Println("项目地址：  https://github.com/laorange/FastestBilibiliDownloader")
	fmt.Println("原项目地址：https://github.com/sodaling/FastestBilibiliDownloader")
	fmt.Println("\n\n支持以下几种格式的输入：")
	fmt.Println("·  https://www.bilibili.com/video/旧版的av号/ | av号 是以`av`开头的一串数字")
	fmt.Println("·  https://www.bilibili.com/video/新版的BV号/ | BV号 是以`BV`开头的一串字符")
	fmt.Println("·  https://space.bilibili.com/UP主的ID/       | UP主的ID 是一串数字")
	fmt.Print("\n\n请输入想要下载的视频网址/up主个人主页网址: ")
	fmt.Scan(&urlInput)

	// bvid
	bvidRegexp := regexp.MustCompile(`/?(BV\w+)[/?]?`)
	params = bvidRegexp.FindStringSubmatch(urlInput)
	if params != nil {
		idType = "bvid"
		bvid = params[1]
	}

	// aid
	aidRegexp := regexp.MustCompile(`/?(av\d+)/?`)
	params = aidRegexp.FindStringSubmatch(urlInput)
	if params != nil {
		idType = "aid"
		aid, _ = strconv.ParseInt(params[1], 10, 64)
	}

	// upid
	upidRegexp := regexp.MustCompile(`space.bilibili.com/(\d+)/?`)
	params = upidRegexp.FindStringSubmatch(urlInput)
	if params != nil {
		idType = "upid"
		upid, _ = strconv.ParseInt(params[1], 10, 64)
	}

	if idType == "aid" {
		req = parser.GetRequestByAid(aid)
	} else if idType == "bvid" {
		aid = parser.Bv2av(bvid)
		req = parser.GetRequestByAid(aid)
	} else if idType == "upid" {
		req = parser.GetRequestByUpId(upid)
	} else {
		req = nil
		log.Fatalln("您输入的网址无法解析，请查证后重试")
		os.Exit(1)
	}

	queueScheduler := scheduler.NewConcurrentScheduler()
	conEngine := engine.NewConcurrentEngine(30, queueScheduler, itemChan)
	log.Println("开始下载...")
	conEngine.Run(req)
	wg.Wait()
	log.Print("所有视频均已下载完成。按 Ctrl+C 来退出程序。")
	var eof string
	fmt.Scan(&eof)
}
