package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"simple-golang-crawler/engine"
	"simple-golang-crawler/parser"
	"simple-golang-crawler/persist"
	"simple-golang-crawler/scheduler"

	"github.com/alexflint/go-arg"
)

var cmdArgs struct {
	IdType string `arg:"-t,--type" help:"id type,support:\n\taid:\t\t下载单个视频\n\tupid\t\t下载指定up主的视频"`
	Id     int64  `arg:"-i,--id" help:"视频或up主的id"`
}

func main() {
	var err error
	var idType string
	var id int64
	arg.MustParse(&cmdArgs)
	if cmdArgs.IdType == "" {
		fmt.Println("Please enter your id type(`aid` or `upid`)")
		fmt.Scan(&idType)
		fmt.Println("Please enter your id")
		fmt.Scan(&id)
	} else {
		idType = cmdArgs.IdType
		id = cmdArgs.Id
	}

	var req *engine.Request
	if idType == "aid" {
		req = parser.GetRequestByAid(id)
	} else if idType == "upid" {
		req = parser.GetRequestByUpId(id)
	} else {
		log.Fatalln("Wrong type you enter")
		os.Exit(1)
	}

	itemProcessFun := persist.GetItemProcessFun()
	var wg sync.WaitGroup
	wg.Add(1)
	itemChan, err := itemProcessFun(&wg)
	if err != nil {
		panic(err)
	}

	queueScheduler := scheduler.NewConcurrentScheduler()
	conEngine := engine.NewConcurrentEngine(30, queueScheduler, itemChan)
	log.Println("Start working.")
	conEngine.Run(req)
	wg.Wait()
	log.Println("All work has done")
}
