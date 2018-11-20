package twitter

import (
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

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

// GetApi returns the private api of our client
func (client Client) GetAPI() *anaconda.TwitterApi {
	return client.api
}

// NewClient returns a Twitter api client
func NewClient(accessToken string, accessSecret string, apiKey string, apiSecret string) Client {
	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessSecret, apiKey, apiSecret)
	api.EnableThrottling(10*time.Second, 60)
	api.SetDelay(5 * time.Second)
	var client = Client{api: api}
	return client
}

// SearchTweets search tweets containing a string
func (client Client) SearchTweets(searchQuery string) []Tweet {

	v := url.Values{}
	v.Set("count", "50")
	v.Set("result_type", "popular")

	var resultingTweets []Tweet

	searchResult, _ := client.api.GetSearch(searchQuery, v)
	for _, tweet := range searchResult.Statuses {
		createdAtTime, _ := tweet.CreatedAtTime()
		resultingTweets = append(resultingTweets, Tweet{
			ID:           tweet.Id,
			TimeCreated:  createdAtTime,
			FullText:     tweet.FullText,
			User:         tweet.User,
			RetweetCount: tweet.RetweetCount,
		})
	}
	return resultingTweets
}

// TweetSomething posts a status with options
func (client Client) TweetSomething(status string) error {
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
	_, err := client.api.Retweet(tweetID, true)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
