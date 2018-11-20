package contest

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jarnix/contestator/twitter"
)

// DataFolder is where already parsed tweets are stored
const DataFolder = "data/contest"
const filesToKeep = 300

// GetContestTweets returns the latest tweets containing interesting words
func GetContestTweets(twitterClient *twitter.Client) {
	queries := strings.Split(os.Getenv("CONTEST_PLAY_QUERIES"), ",")

	// var contestTweets []twitter.Tweet

	for _, query := range queries {
		tweets := twitterClient.SearchTweets(query)
		for _, tweet := range tweets {

			// ignore all the tweets with 0 retweets that are probably false positive
			if tweet.RetweetCount > 0 {
				filepath := DataFolder + "/" + strconv.FormatInt(tweet.ID, 10) + ".txt"

				if _, err := os.Stat(filepath); os.IsNotExist(err) {
					/*
						f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0600)
						if err != nil {
							log.Fatal(err)
						}
					*/

					fmt.Println(tweet.FullText)

					var usersToFollow = make(map[string]bool)

					usersToFollow[tweet.User.ScreenName] = true

					// analyze the tweet
					// look for @ mentions => follow them too
					var regexMentions = regexp.MustCompile(`(?m)@(\w{1,15})`)
					tweetMentions := regexMentions.FindAllStringSubmatch(tweet.FullText, -1)
					for _, match := range tweetMentions {
						if _, ok := usersToFollow[match[1]]; !ok {
							usersToFollow[match[1]] = true
						}
					}
					for userToFollow := range usersToFollow {
						log.Println("followed: ", userToFollow)
						twitterClient.GetAPI().FollowUser(userToFollow)
					}

					// like the post

					// retweet

					// mention a friend in the reply
					// add all the hashtags from the tweet

					/*
						defer f.Close()
						if _, err2 := f.WriteString(tweet.FullText); err2 != nil {
							log.Fatal(err2)
						}
					*/

					panic("haha")
				} else {
					log.Println("this tweet was already downloaded")
				}

				panic("haha")
			}
		}

	}

}
