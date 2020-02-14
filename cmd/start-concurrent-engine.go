package main

import (
	"fmt"
	"log"
	"os"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
	"simple-golang-crawler/parser"
	"simple-golang-crawler/persist"
	"simple-golang-crawler/scheduler"
	"simple-golang-crawler/tool"
)

func main() {
	if !tool.CheckFfmegStatus() {
		fmt.Println("Can't locate your ffmeg.The video your download can't be merged")
	}
	itemChan, err := persist.VideoItemProcessor()
	var idType string
	var id int64
	var req *engine.Request
	fmt.Println("Please enter your id type(`aid` or `upid`)")
	fmt.Scan(&idType)
	fmt.Println("Please enter your id")
	fmt.Scan(&id)

	if idType == "aid" {
		var videoAid *model.VideoAidInfo
		req, videoAid = parser.GetRequestByAid(id)
		item := engine.NewItem(videoAid)
		go func() { itemChan <- item }()
	} else if idType == "upid" {
		req = parser.GetRequestByUpId(id)
	} else {
		fmt.Println("Wrong type you enter")
		os.Exit(1)
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	queueScheduler := scheduler.NewConcurrentScheduler()
	conEngine := engine.NewConcurrentEngine(10, queueScheduler, itemChan)
	conEngine.Run(req)
	fmt.Println("All work has done")
}
