package contest

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jarnix/contestator/twitter"
)

// DataFolder is where already parsed tweets are stored
const DataFolder = "data/contest"
const filesToKeep = 300

// GetContestTweets returns the latest tweets containing interesting words
func GetContestTweets(twitterClient *twitter.Client, resultType string) {
	rand.Seed(time.Now().UnixNano())

	queries := strings.Split(os.Getenv("CONTEST_PLAY_QUERIES"), ",")

	for _, query := range queries {

		fmt.Println("query", query)

		v := url.Values{}
		v.Set("lang", os.Getenv("CONTEST_LANG"))
		v.Set("count", "100")
		v.Set("result_type", resultType)

		tweets := twitterClient.SearchTweets(query, v)

		log.Println("query", query)
		log.Println("results", len(tweets))

		for _, tweet := range tweets {

			fmt.Println(strings.Repeat("*", 80))

			// ignore all the tweets with 0 retweets that are probably false positive
			if tweet.RetweetCount > 0 {
				filepath := DataFolder + "/" + strconv.FormatInt(tweet.ID, 10) + ".txt"

				var err error

				if _, err = os.Stat(filepath); os.IsNotExist(err) {

					f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0600)
					if err != nil {
						log.Fatal(err)
					}

					log.Println("screenName", tweet.User.ScreenName)
					log.Println("fullText", tweet.FullText)

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
						_, err = twitterClient.GetAPI().FollowUser(userToFollow)
						if err != nil {
							log.Println(err)
						}
					}

					time.Sleep(5 * time.Second)

					// retweet
					_, err = twitterClient.GetAPI().Retweet(tweet.ID, true)
					if err != nil {
						log.Println(err)
					}

					time.Sleep(5 * time.Second)

					var hashtagsToAdd = make(map[string]bool)
					var regexHashtags = regexp.MustCompile(`(?m)#(\w{1,15})`)
					hashtags := regexHashtags.FindAllStringSubmatch(tweet.FullText, -1)
					for _, match := range hashtags {
						if _, ok := hashtagsToAdd[match[1]]; !ok {
							hashtagsToAdd[match[1]] = true
						}
					}
					hashtagString := ""
					for hashtag := range hashtagsToAdd {
						hashtagString += " #" + string(hashtag)
					}

					time.Sleep(5 * time.Second)

					// mention a friend in the reply
					// with all the hashtags from the tweet
					v := url.Values{}
					v.Set("in_reply_to_status_id", strconv.FormatInt(tweet.ID, 10))
					var intros = strings.Split(os.Getenv("CONTEST_PLAY_INTROS"), ",")
					reply := "@" + tweet.User.ScreenName + " " + intros[rand.Intn(len(intros)-1)] + " " + os.Getenv("CONTEST_MENTION")
					reply += hashtagString

					_, err = twitterClient.GetAPI().PostTweet(reply, v)
					if err != nil {
						log.Fatal(err)
					}

					defer f.Close()
					if _, err2 := f.WriteString(tweet.FullText); err2 != nil {
						log.Fatal(err2)
					}

					time.Sleep(5 * time.Second)

				} else {
					log.Println("this tweet was already downloaded")
				}

			}
		}

	}

}
