package main

import (
	_ "database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
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

	dbHost := os.Getenv("DB_HOST")
	fmt.Println(dbHost)

	api := anaconda.NewTwitterApiWithCredentials(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"), os.Getenv("TWITTER_API"), os.Getenv("TWITTER_SECRET"))

	v := url.Values{}
	v.Set("count", "30")

	searchResult, _ := api.GetSearch("#concours", v)
	for _, tweet := range searchResult.Statuses {
		fmt.Println(tweet.Text)
	}
}
