package main

import (
	_ "database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/dghubble/go-twitter/twitter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// aller chercher 4 tweets via search avec #concours, rt, retweet, follow (...)
// voir dans le tweet les comptes à suivre, s'il y a des indications particulières
// stocker les comptes followés dans la base pour les défollower 1 mois après
// attendre au moins 6 heures entre chaque action, les faire en journée uniquement (à heure régulière)

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

}
