package contest

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jarnix/contestator/twitter"
)

// DataFolder is where already parsed tweets are stored
const DataFolder = "data/contest"

// GetContestTweets returns the latest tweets containing interesting words
func GetContestTweets(twitterClient *twitter.Client) {
	queries := strings.Split(os.Getenv("CONTEST_PLAY_QUERIES"), ",")

	for _, query := range queries {
		tweets := twitterClient.SearchTweets(query)
		for _, tweet := range tweets {
			filepath := DataFolder + "/" + strconv.FormatInt(tweet.ID, 10) + ".txt"
			fmt.Println(filepath)
			/*
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
			*/
		}

	}

}
