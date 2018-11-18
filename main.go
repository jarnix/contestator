package main

import (
	_ "database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/jarnix/contestator/markov"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// aller chercher 4 tweets via search avec #concours, rt, retweet, follow (...)
// voir dans le tweet les comptes à suivre, s'il y a des indications particulières
// stocker les comptes followés dans la base pour les défollower 1 mois après
// attendre au moins 2/3 heures entre chaque lancement et 1/2 minutes entre chaque action
// les faire en journée uniquement (à heure régulière)

// anti anti bot
// poster des trucs en allant chercher les news ailleurs pour faire genre le compte est normal
// poster des status à la con
// retweeter des célébrités du milieu

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

	todo := flag.String("todo", "", "action to launch (downloadforindex, ...)")
	flag.Parse()
	log.SetPrefix(*todo + " ")
	switch *todo {
	case "downloadforindex":
		markov.DownloadForIndex(1)
	case "markovgenerate":
		fmt.Println(markov.GenerateText(markov.DataFolder, 2, 2, 20))
	case "":
		flag.PrintDefaults()
	}

	/*
		twitterClient := twitter.NewClient(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"), os.Getenv("TWITTER_API"), os.Getenv("TWITTER_SECRET"))
		twitterClient.SearchTweets("#concours")
	*/

}
