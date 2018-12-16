package retweet

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jarnix/contestator/twitter"
	"github.com/jarnix/contestator/webreader"
)

// DataFolder is the root folder containing .txt files for generating text
const DataFolder = "data/retweet"
const filesToKeep = 20

// DownloadForRetweet gets some tweet to retweet
func DownloadForRetweet(idx int) {
	indexURL := os.Getenv("URL_RETWEET_INDEX_" + strconv.Itoa(idx))

	// get the number of files for this idx and keep only the most toKeep recent ones
	files, err := ioutil.ReadDir(DataFolder + "/" + strconv.Itoa(idx))
	if err != nil {
		os.MkdirAll(DataFolder+"/"+strconv.Itoa(idx), os.ModePerm)
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

	// crawl the index on multiple pages and fetch the content for building our retweet index
	webpageWithLinks := webreader.GetWebPage(indexURL)

	baseURL := webreader.GetBaseURL(indexURL)

	links := webreader.GetLinksFromQuery(webpageWithLinks, os.Getenv("URL_RETWEET_QUERY_CRAWL_"+strconv.Itoa(idx)), baseURL)

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

			// extract all the tweets for the downloaded webpage and put them in the file
			tweets := webreader.GetLinksFromQuery(webpageArticle, os.Getenv("URL_RETWEET_QUERY_TWEET_"+strconv.Itoa(idx)), baseURL)

			var regexTweetID, _ = regexp.Compile(`/status/(\d+)`)

			var tweetsIDs []string

			for _, tweet := range tweets {
				if regexTweetID.MatchString(tweet.Href) {
					tweetID := regexTweetID.FindStringSubmatch(tweet.Href)[1]
					tweetsIDs = append(tweetsIDs, tweetID)
				}
			}

			defer f.Close()

			if _, err2 := f.WriteString(strings.Join(tweetsIDs, "\n")); err2 != nil {
				log.Fatal(err2)
			}

			log.Println(link.Href + " was downloaded")

		} else {
			log.Println(link.Href + " was already downloaded")
		}

		time.Sleep(500 * time.Millisecond)

	}

}

// PostRandomRetweet gets a random tweet from the index
func PostRandomRetweet(twitterClient *twitter.Client) {
	rand.Seed(time.Now().UnixNano())

	// contains all the files to parse
	var allTextFilesInsideRoot []string

	// recursively read the folder
	err := filepath.Walk(DataFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".txt") {
				allTextFilesInsideRoot = append(allTextFilesInsideRoot, path)
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	randomFileNumber := rand.Intn(len(allTextFilesInsideRoot) - 1)

	file, _ := os.Open(allTextFilesInsideRoot[randomFileNumber])
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	randomTweetID, _ := strconv.ParseInt(lines[rand.Intn(len(lines)-1)], 10, 64)

	twitter.RandomSleep()

	_, err = twitterClient.GetAPI().Retweet(randomTweetID, true)
	if err != nil {
		log.Fatal(err)
	}

}
