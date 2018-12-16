package twitter

import (
	"errors"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const minSleep = 60
const maxSleep = 900

// Client is our wrapper for anaconda TwitterApi
type Client struct {
	api *anaconda.TwitterApi
}

// Tweet is a tweet...
type Tweet struct {
	ID           int64
	TimeCreated  time.Time
	FullText     string
	User         anaconda.User
	RetweetCount int
}

// GetAPI returns the private api of our client
func (client Client) GetAPI() *anaconda.TwitterApi {
	return client.api
}

// NewClient returns a Twitter api client
func NewClient(accessToken string, accessSecret string, apiKey string, apiSecret string) Client {
	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessSecret, apiKey, apiSecret)
	// api.EnableThrottling(600*time.Second, 10)
	// api.SetDelay(10 * time.Second)
	var client = Client{api: api}
	return client
}

// RandomSleep sleeps for a random time
func RandomSleep() {
	rand.Seed(time.Now().UnixNano())
	randomTime := rand.Intn(maxSleep-minSleep) + minSleep
	log.Println("Sleeping for", randomTime, "seconds")
	time.Sleep(time.Duration(randomTime) * time.Second)
}

// SearchTweets search tweets containing a string
func (client Client) SearchTweets(searchQuery string, v url.Values) []Tweet {

	if v == nil {
		v = url.Values{}
		v.Set("count", "50")
		v.Set("result_type", "popular")
	}

	var resultingTweets []Tweet

	searchResult, _ := client.api.GetSearch(searchQuery, v)

	for _, tweet := range searchResult.Statuses {
		createdAtTime, _ := tweet.CreatedAtTime()
		appendedTweet := Tweet{
			ID:           tweet.Id,
			TimeCreated:  createdAtTime,
			FullText:     tweet.FullText,
			User:         tweet.User,
			RetweetCount: tweet.RetweetCount,
		}
		// if it's a retweet, we return the original tweet and not the retweet
		if tweet.RetweetedStatus != nil {
			retweetedTweet := *tweet.RetweetedStatus
			appendedTweet.ID = retweetedTweet.Id
			createdAtTime, _ := retweetedTweet.CreatedAtTime()
			appendedTweet.TimeCreated = createdAtTime
			appendedTweet.FullText = retweetedTweet.FullText
			appendedTweet.User = retweetedTweet.User
			appendedTweet.RetweetCount = tweet.RetweetCount
		}

		resultingTweets = append(resultingTweets, appendedTweet)
	}
	return resultingTweets
}

// TweetSomething posts a status with options
func (client Client) TweetSomething(status string) error {
	RandomSleep()
	if status == "" {
		return errors.New("status cannot be empty")
	}
	_, err := client.api.PostTweet(status, nil)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Retweet retweets a tweet by its id
func (client Client) Retweet(tweetID int64) error {
	RandomSleep()
	_, err := client.api.Retweet(tweetID, true)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
