package persist

import (
	"log"
	"simple-golang-crawler/engine"
)

func FileItemSaver(savePath string) (chan engine.Item, error) {
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			//log.Printf("Item Saver:got item "+
			//	"#%d: %v", itemCount, item)
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

func save(item engine.Item) error {
	//video := item.Payload.(*model.Video)
	//fmt.Println((*video).Aid)
	return nil
}
