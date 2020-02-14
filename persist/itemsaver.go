package persist

import (
	"fmt"
	"log"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/model"
)

func VideoItemProcessor() (chan *engine.Item, error) {
	out := make(chan *engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver:got item "+
				"#%d: %v", itemCount, item)
			itemCount++
			err := save(item)
			if err != nil {
				log.Printf("Item Saver: error "+
					"saving item %v:%v", item, err)
			}
		}
	}()
	return out, nil
}

func save(item *engine.Item) error {
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
