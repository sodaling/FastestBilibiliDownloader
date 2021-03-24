package main

import (
	//"fmt"
	"log"
	"os"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/parser"
	"simple-golang-crawler/persist"
	"simple-golang-crawler/scheduler"
	"sync"
	"flag"
	"strconv"
)

func main() {

	var arg_idType *string = flag.String("t", "", "id type (i.e. aid, bvid or upid)")
	var arg_id *string = flag.String("v", "", "id (直接输入id不需要加双引号)")
	var arg_worker *int = flag.Int("w", 30, "number of workers for this id, depends on the videos to download")
	flag.Parse()
	//flag.PrintDefaults()

	// 如果没有输入任何值
	if *arg_idType == "" {
	    log.Fatalln("No argument entered, using -h to find what is required")
		os.Exit(1)
	}

	itemProcessFun := persist.GetItemProcessFun()
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	itemChan, err := itemProcessFun(&wg)
	if err != nil {
		panic(err)
	}

	var req *engine.Request
	var idType string = *arg_idType
	var id string = *arg_id
	var num_worker int = *arg_worker

	if idType == "aid" {
        int_id,_ := strconv.ParseInt(id, 10, 64)
		req = parser.GetRequestByAid(int_id)
	} else if idType == "bvid" {
	    aid := parser.Bv2av(id)
	    req = parser.GetRequestByAid(aid)
	} else if idType == "upid" {
	    int_id,_ := strconv.ParseInt(id, 10, 64)
		req = parser.GetRequestByUpId(int_id)  
	} else {
		log.Fatalln("Wrong type you enter")
		os.Exit(1)
	}

	queueScheduler := scheduler.NewConcurrentScheduler()
	conEngine := engine.NewConcurrentEngine(num_worker, queueScheduler, itemChan)
	log.Println("Start working.")
	conEngine.Run(req)
	wg.Wait()
	log.Println("All work has done")
}
