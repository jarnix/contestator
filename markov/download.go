package markov

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jarnix/contestator/webreader"
)

// DownloadForIndex downloads some text of an article for generating markov chains
func DownloadForIndex(idx int) {
	indexURLBase := os.Getenv("URL_INDEX_MARKOV_" + strconv.Itoa(idx))

	pageNumber := 2
	for pageNumber < 10 {
		indexURL := strings.Replace(indexURLBase, "##page##", strconv.Itoa(pageNumber), 1)
		webpageWithLinks := webreader.GetWebPage(indexURL)
		baseURL := webreader.GetBaseURL(os.Getenv("URL_BASE_CRAWL_" + strconv.Itoa(idx)))

		links := webreader.GetLinksFromQuery(webpageWithLinks, os.Getenv("URL_QUERY_CRAWL_"+strconv.Itoa(idx)), baseURL)

		for _, link := range links {

			hasher := md5.New()
			hasher.Write([]byte(link.Href))
			filepath := "data/" + strconv.Itoa(idx) + "-" + hex.EncodeToString(hasher.Sum(nil)) + ".txt"
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					log.Fatal(err)
				}

				webpageArticle := webreader.GetWebPage(link.Href)
				titleFromQuery := webreader.GetTextFromQuery(webpageArticle, os.Getenv("URL_QUERY_TITLE_"+strconv.Itoa(idx)))
				contentFromArticle, imageFromArticle := webreader.GetArticleContentAndImage(link.Href, webpageArticle)

				defer f.Close()

				if _, err2 := f.WriteString(contentFromArticle + "\n"); err2 != nil {
					log.Fatal(err2)
				}

				log.Println(link.Href + " was downloaded")

				fmt.Println(titleFromQuery, contentFromArticle, imageFromArticle)

			} else {
				log.Println(link.Href + " was already downloaded")
			}

			time.Sleep(500 * time.Millisecond)

		}

		pageNumber++
		log.Println("downloading, pageNumber", pageNumber)
	}

}
