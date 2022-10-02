package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/SultanKs4/walfie-scrap/internal"
	"github.com/SultanKs4/walfie-scrap/types"
)

func main() {
	var wg sync.WaitGroup

	var dataCh = make(chan types.ResponseGetUrl)
	errCh := make(chan error)
	// quitCh := make(chan struct{})
	for i := 1; ; i++ {
		wg.Add(1)
		fmt.Printf("get data from page: %v\n", i)
		go internal.GetLink(i, &wg, dataCh, errCh)

		if err := <-errCh; err != nil {
			if err.Error() == "last page" {
				break
			}
			log.Fatalf("get link: %v", err)
		}
	}

	pack := 1
	i := 1
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("start get url image and save")
	for k := range dataCh {
		wg.Add(1)
		if i%2 == 0 {
			pack++
		}
		folder := fmt.Sprintf("pack %s", strconv.Itoa(pack))
		path := filepath.Join(wd, "../", "walfie gifs", folder)
		os.MkdirAll(path, os.ModePerm)
		go internal.ScrapImgLink(k.Html, &wg, path)
		i++
	}

	wg.Wait()
	fmt.Println("job done")
}
