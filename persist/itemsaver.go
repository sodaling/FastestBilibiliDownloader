package persist

import (
	"fmt"
	"log"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
	"sync"
)

var videPageMap = make(map[int64]int64)

func VideoItemProcessor(wgOutside *sync.WaitGroup) (chan *engine.Item, error) {
	out := make(chan *engine.Item)
	go func() {
		defer wgOutside.Done()
		var wgInside sync.WaitGroup
		for item := range out {

			switch x := item.Payload.(type) {
			case *model.VideoAidInfo:
				fmt.Println("aid:", x.Aid)
				videPageMap[x.Aid] = x.GetPage()
				fmt.Println(videPageMap[x.Aid])
			case *model.VideoCidInfo:
				fmt.Println("cid:", x.Cid)
				videPageMap[x.ParAid.Aid] -= 1
				fmt.Println(videPageMap[x.ParAid.Aid])
				if videPageMap[x.ParAid.Aid] == 0 {
					wgInside.Add(1)
					go mergeVideo(x, &wgInside)
				}
			default:
				panic(fmt.Sprintf("unexpected type %T: %v", x, x))
			}

		}
		wgInside.Wait()
	}()
	return out, nil
}

func mergeVideo(videoCiD *model.VideoCidInfo, wg *sync.WaitGroup) {
	defer wg.Done()
}

func VideoItemCleaner() (chan *engine.Item, error) {
	out := make(chan *engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver:got item "+
				"#%d: %v", itemCount, item)
			itemCount++
		}
	}()
	return out, nil
}
