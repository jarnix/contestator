package twitter

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

// Client is our wrapper for anaconda TwitterApi
type Client struct {
	api *anaconda.TwitterApi
}

// NewClient returns a Twitter api client
func NewClient(accessToken string, accessSecret string, apiKey string, apiSecret string) Client {
	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessSecret, apiKey, apiSecret)
	// api.SetDelay(10 * time.Second)
	var client = Client{api: api}
	return client
}

// SearchTweets search tweets containing a string
func (client Client) SearchTweets(searchQuery string) {

	v := url.Values{}
	v.Set("count", "30")

	searchResult, _ := client.api.GetSearch("#concours", v)
	for _, tweet := range searchResult.Statuses {
		fmt.Println("id", tweet.Id)
		fmt.Println("fullText", tweet.FullText)
	}
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
