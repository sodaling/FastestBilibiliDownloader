package main

import (
	"log"
	"os"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/parser"
	"simple-golang-crawler/persist"
	"simple-golang-crawler/scheduler"
)

func main() {
	itemChan, err := persist.FileItemSaver(".")

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	queueScheduler := scheduler.NewConcurrentScheduler()
	conEngine := engine.NewConcurrentEngine(10, queueScheduler, itemChan)

	req := engine.NewRequest("https://api.bilibili.com/x/space/arc/search?mid=585267&ps=30&tid=0&pn=2&keyword=&order=pubdate&jsonp=jsonp", parser.UpSpaceParseFun, fetcher.DefaultFetcher)

	conEngine.Run(req)

}
