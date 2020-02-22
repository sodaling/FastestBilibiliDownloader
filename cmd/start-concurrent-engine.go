package main

import (
	"fmt"
	"log"
	"os"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/parser"
	"simple-golang-crawler/persist"
	"simple-golang-crawler/scheduler"
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

	var idType string
	var id int64
	var req *engine.Request
	fmt.Println("Please enter your id type(`aid` or `upid`)")
	fmt.Scan(&idType)
	fmt.Println("Please enter your id")
	fmt.Scan(&id)

	if idType == "aid" {
		req = parser.GetRequestByAid(id)
	} else if idType == "upid" {
		req = parser.GetRequestByUpId(id)
	} else {
		log.Fatalln("Wrong type you enter")
		os.Exit(1)
	}

	queueScheduler := scheduler.NewConcurrentScheduler()
	conEngine := engine.NewConcurrentEngine(10, queueScheduler, itemChan)
	log.Println("Start working.")
	conEngine.Run(req)
	wg.Wait()
	log.Println("All work has done")
}
