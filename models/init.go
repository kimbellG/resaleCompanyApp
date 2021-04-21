package models

import (
	"log"
	"os"

	"database/sql"

	neasted "github.com/antonfisher/nested-logrus-formatter"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type DBController struct {
	*sql.DB
}

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&neasted.Formatter{
		HideKeys: true,
	})
	logger.SetLevel(logrus.DebugLevel)
	initEnv()

}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func CreateNewDBConnection() *DBController {
	lib_db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db := &DBController{lib_db}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	db.createTables()

	return db
}

func (db *DBController) createTables() {
	if err := db.createUserInfoTable(); err != nil {
		log.Fatalln("create user information table:", err)
	}

	if err := db.createAuthInfoTable(); err != nil {
		log.Fatalln("create authtorization information table:", err)
	}
}

func (db *DBController) createUserInfoTable() error {
	err := db.createTable(
		`CREATE TABLE IF NOT EXISTS UserInformation (
			id 		SERIAL PRIMARY KEY,
			surname VARCHAR(50) NOT NULL,
			name 	VARCHAR(50) NOT NULL
		);`)

	return err
}

func (db *DBController) createAuthInfoTable() error {
	err := db.createTable(
		`CREATE TABLE IF NOT EXISTS AuthtorizationInformation (
 		id 		 SERIAL PRIMARY KEY UNIQUE,
 		login	 VARCHAR(30), 
 		password VARCHAR(30),
 		status	 INTEGER,
		userInfo INTEGER REFERENCES UserInformation (id)
	);`)

	return err
}

func (db *DBController) createTable(query string) error {
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}
