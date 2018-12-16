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
// attendre au moins 2/3 heures entre chaque lancement et 1/2 minutes entre chaque action
// les faire en journée uniquement (à heure régulière)

// anti anti bot

// ## poster des status à la con via markov préfixés par "#gamingsousLSD"
// ## tweeter des emojis en messages codés avec des emojis aléatoires
// ## retweeter des conneries depuis ce site http://twog.fr/
// traduire /r/showerthoughts via google translate api
// retweeter des comptes auxquels je suis abonné
// retweeter des célébrités du milieu

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	twitterClient := twitter.NewClient(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"), os.Getenv("TWITTER_API"), os.Getenv("TWITTER_SECRET"))

	/*
		b := []string{"users", "search", "statuses"}
		rateLimits, _ := twitterClient.GetAPI().GetRateLimits(b)
		fmt.Println(rateLimits)
	*/

	todo := flag.String("todo", "", "action to launch (downloadforindex, ...)")
	flag.Parse()
	log.SetPrefix(*todo + " ")
	switch *todo {
	case "downloadforindex":
		markov.DownloadForMarkov(1)
	case "downloadforretweet":
		retweet.DownloadForRetweet(1)
	case "tweetretweet":
		retweet.PostRandomRetweet(&twitterClient)
	case "tweetmarkov":
		stupidText := markov.GenerateText(2, 1, 10)
		twitterClient.TweetSomething(stupidText)
	case "tweetemojis":
		stupidText := emoji.GenerateText()
		twitterClient.TweetSomething(stupidText)
	case "contestplay":
		contest.GetContestTweets(&twitterClient, "popular")
		contest.GetContestTweets(&twitterClient, "latest")
	case "":
		flag.PrintDefaults()
	}
}
