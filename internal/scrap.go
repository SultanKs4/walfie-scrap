package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/SultanKs4/walfie-scrap/types"
)

func sanitizeHtml(html string) (sanitizeHtml string) {
	// sanitize html
	replacer := strings.NewReplacer("\n", "", "\t", "", `\"`, `"`, "#038;", "", "<!-- #post-## -->", "", "<!-- .entry-content -->", "", "<!-- .entry-header -->", "", "<!-- .entry-meta -->", "")
	sanitizeHtml = replacer.Replace(html)
	return
}

func GetLink(page int, wg *sync.WaitGroup, dataCh chan types.ResponseGetUrl, errCh chan error) {
	defer wg.Done()
	res, err := http.PostForm("https://walfiegif.wordpress.com/?infinity=scrolling", url.Values{
		"page":  []string{strconv.Itoa(page)},
		"order": []string{"DESC"},
	})
	if err != nil {
		errCh <- err
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errCh <- err
		return
	}
	var bodyResponse types.ResponseGetUrl
	if err := json.NewDecoder(res.Body).Decode(&bodyResponse); err != nil {
		errCh <- err
		return
	}

	if bodyResponse.Lastbatch {
		errCh <- errors.New("last page")
		dataCh <- bodyResponse
		close(dataCh)
		return
	}
	errCh <- nil
	dataCh <- bodyResponse
}

func ScrapImgLink(html string, wg *sync.WaitGroup, path string) {
	defer wg.Done()
	// Sanitize HTML string from response body
	sanitizeHtml := sanitizeHtml(html)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sanitizeHtml))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".entry-image img").Each(func(i int, s *goquery.Selection) {
		url, exist := s.Attr("src")
		if !exist {
			log.Fatal("attribute src not found")
		}
		title, _ := s.Attr("title")
		// url = strings.Replace(url, "?w=560&h=9999", "?w=512", -1)
		wg.Add(1)
		go saveImage(url, title, path, wg)
	})
}

func saveImage(url, title, path string, wg *sync.WaitGroup) {
	defer wg.Done()

	output, err := os.Create(fmt.Sprintf("%s/%s", path, fmt.Sprintf("%v.gif", title)))
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	_, err = io.Copy(output, res.Body)
	if err != nil {
		log.Fatal(err)
	}
}
