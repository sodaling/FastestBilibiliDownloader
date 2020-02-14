package main

import (
	"fmt"
	"log"
	"os"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/parser"
	"simple-golang-crawler/persist"
	"simple-golang-crawler/scheduler"
)

func main() {
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
		fmt.Println("Wrong type you enter")
		os.Exit(1)
	}
	itemChan, err := persist.FileItemSaver(".")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	queueScheduler := scheduler.NewConcurrentScheduler()
	conEngine := engine.NewConcurrentEngine(10, queueScheduler, itemChan)
	conEngine.Run(req)
	fmt.Println("All work has done")
}
