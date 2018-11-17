package connector

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

var dbConnections map[string]*sql.DB
var dbOnce map[string]bool
var dbOnceMutex sync.Mutex

// Db returns a connection
func Db(name string) *sql.DB {

	if name == "" {
		log.Println("DB's name cannot be empty")
		return nil
	}

	if len(dbOnce) == 0 {
		dbOnce = make(map[string]bool, 15)
		dbConnections = make(map[string]*sql.DB, 15)
	}

	dbOnceMutex.Lock()
	if !dbOnce[name] {
		dbOnce[name] = true
		connectionStringArr := []string{
			os.Getenv("DB_" + name + "_USER"),
			":",
			os.Getenv("DB_" + name + "_PASS"),
			"@(",
			os.Getenv("DB_" + name + "_HOST"),
			":",
			os.Getenv("DB_" + name + "_PORT"),
			")/",
			os.Getenv("DB_" + name + "_BASE"),
			"?parseTime=true",
		}
		dbConnections[name], _ = sql.Open("mysql", strings.Join(connectionStringArr, ""))
		err := dbConnections[name].Ping()
		if err != nil {
			fmt.Println(err.Error())
		}
		dbConnections[name].SetConnMaxLifetime(time.Second)
		dbOnceMutex.Unlock()
	} else {
		dbOnceMutex.Unlock()
	}

	return dbConnections[name]

}
