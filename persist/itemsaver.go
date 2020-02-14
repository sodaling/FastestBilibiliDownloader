package persist

import (
	"fmt"
	"log"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
	"sync"
)

func VideoItemProcessor(wgOutside *sync.WaitGroup) (chan *engine.Item, error) {
	out := make(chan *engine.Item)
	go func() {
		defer wgOutside.Done()
		var wgInside sync.WaitGroup
		itemCount := 0
		for item := range out {
			log.Printf("Item Saver:got item "+
				"#%d: %v", itemCount, item)
			itemCount++
			err := save(item, &wgInside)
			if err != nil {
				log.Printf("Item Saver: error "+
					"saving item %v:%v", item, err)
			}
		}
		wgInside.Wait()
	}()
	return out, nil
}

func save(item *engine.Item, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()
	switch x := item.Payload.(type) {
	case *model.VideoCidInfo:
		fmt.Println("cid:", *x)
	case *model.VideoAidInfo:
		fmt.Println("aid:", *x)
	default:
		panic(fmt.Sprintf("unexpected type %T: %v", x, x))
	}
	return nil
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
