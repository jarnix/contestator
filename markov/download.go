package markov

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jarnix/contestator/twitter"
	"github.com/jarnix/contestator/webreader"
)

// DataFolder is the root folder containing .txt files for generating text
const DataFolder = "data/markov"
const filesToKeep = 30

// DownloadForMarkov downloads some text of an article for generating markov chains
func DownloadForMarkov(idx int) {
	indexURLBase := os.Getenv("URL_MARKOV_INDEX_" + strconv.Itoa(idx))

	// get the number of files for this idx and keep only the most toKeep recent ones
	files, err := ioutil.ReadDir(DataFolder + "/" + strconv.Itoa(idx))
	if err != nil {
		log.Fatal("cannot read data folder")
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() < files[j].ModTime().Unix()
	})
	// keep only "toKeep" files
	var filesToDelete []os.FileInfo
	if len(files) > filesToKeep {
		filesToDelete = files[filesToKeep:]
		files = files[0:filesToKeep]
	}
	// remove older files
	for _, file := range filesToDelete {
		os.Remove(DataFolder + "/" + strconv.Itoa(idx) + "/" + file.Name())
	}

	// crawl the index on multiple pages and fetch the content for building our markov index
	pageNumber := 5
	for pageNumber >= 2 {
		indexURL := strings.Replace(indexURLBase, "##page##", strconv.Itoa(pageNumber), 1)
		webpageWithLinks := webreader.GetWebPage(indexURL)

		links := webreader.GetLinksFromQuery(webpageWithLinks, os.Getenv("URL_MARKOV_QUERY_CRAWL_"+strconv.Itoa(idx)), indexURL)

		for _, link := range links {

			hasher := md5.New()
			hasher.Write([]byte(link.Href))
			filepath := DataFolder + "/" + strconv.Itoa(idx) + "/" + hex.EncodeToString(hasher.Sum(nil)) + ".txt"
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					log.Fatal(err)
				}

				webpageArticle := webreader.GetWebPage(link.Href)
				contentFromArticle, _ := webreader.GetArticleContentAndImage(link.Href, webpageArticle)

				defer f.Close()

				if _, err2 := f.WriteString(contentFromArticle + "\n"); err2 != nil {
					log.Fatal(err2)
				}

				log.Println(link.Href + " was downloaded")

				// fmt.Println(titleFromQuery, contentFromArticle, imageFromArticle)

			} else {
				log.Println(link.Href + " was already downloaded")
			}

			twitter.RandomSleep()

		}

		pageNumber--
		log.Println("downloading, pageNumber", pageNumber)
	}

}
