package main

import (
	_ "database/sql"
	"flag"
	"log"
	"os"

	"github.com/jarnix/contestator/contest"
	"github.com/jarnix/contestator/emoji"
	"github.com/jarnix/contestator/markov"
	"github.com/jarnix/contestator/retweet"
	"github.com/jarnix/contestator/twitter"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// aller chercher 4 tweets via search avec #concours, rt, retweet, follow (...)
// voir dans le tweet les comptes à suivre, s'il y a des indications particulières
// stocker les comptes followés dans la base pour les défollower 1 mois après
// attendre au moins 2/3 heures entre chaque lancement et 1/2 minutes entre chaque action
// les faire en journée uniquement (à heure régulière)

// anti anti bot

// ## poster des status à la con via markov préfixés par "#gamingsousLSD"
// traduire /r/showerthoughts via google translate api
// retweeter des comptes auxquels je suis abonné
// retweeter des célébrités du milieu
// ## tweeter des emojis en messages codés avec des emojis aléatoires
// retweeter des conneries depuis ce site http://twog.fr/

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	/*
		// fetch some websites for tweeting random shit
		webpageWithLinks := webreader.GetWebPage(os.Getenv("URL_BASE_CRAWL_1"))
		baseURL := webreader.GetBaseURL(os.Getenv("URL_BASE_CRAWL_1"))
		links := webreader.GetLinksFromQuery(webpageWithLinks, os.Getenv("URL_QUERY_CRAWL_1"), baseURL)
		for _, link := range links {
			webpageArticle := webreader.GetWebPage(link.Href)
			textFromQuery := webreader.GetTextFromQuery(webpageArticle, os.Getenv("URL_QUERY_TITLE_1"))
			fmt.Println("title", textFromQuery)
			contentFromArticle, imageFromArticle := webreader.GetArticleContentAndImage(link.Href, webpageArticle)
			fmt.Println(contentFromArticle, imageFromArticle)
		}
	*/

	twitterClient := twitter.NewClient(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"), os.Getenv("TWITTER_API"), os.Getenv("TWITTER_SECRET"))

	todo := flag.String("todo", "", "action to launch (downloadforindex, ...)")
	flag.Parse()
	log.SetPrefix(*todo + " ")
	switch *todo {
	case "downloadforindex":
		markov.DownloadForMarkov(1)
	case "downloadforretweet":
		retweet.DownloadForRetweet(1)
	case "tweetretweet":
		randomtweetID := retweet.GetRandomTweet()
		twitterClient.Retweet(randomtweetID)
	case "tweetmarkov":
		stupidText := markov.GenerateText(2, 1, 10)
		twitterClient.TweetSomething(stupidText)
	case "tweetemojis":
		stupidText := emoji.GenerateText()
		twitterClient.TweetSomething(stupidText)
	case "contestplay":
		contest.GetContestTweets(&twitterClient)
	case "":
		flag.PrintDefaults()
	}
}
