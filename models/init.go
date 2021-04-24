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
	DBLogger := logger.WithFields(logrus.Fields{"action": "CreateDBTable"})

	if err := db.createUserInfoTable(); err != nil {
		log.Fatalln("create user information table:", err)
	}

	if err := db.createAuthInfoTable(); err != nil {
		log.Fatalln("create authtorization information table:", err)
	}

	if err := db.createRegionTable(); err != nil {
		DBLogger.Error("Invalid region table")
		os.Exit(1)
	}

	if err := db.createProviderTable(); err != nil {
		DBLogger.Error("Invalid provider table: ", err)
		os.Exit(1)
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
		userInfo INTEGER REFERENCES UserInformation (id) ON DELETE CASCADE
	);`)

	return err
}

func (db *DBController) createProviderTable() error {
	err := db.createTable(
		`CREATE TABLE IF NOT EXISTS Provider (
			vendor_code SERIAL PRIMARY KEY UNIQUE,
			name VARCHAR(200) NOT NULL,
			unp CHAR(9) NOT NULL,
			region_code INTEGER,
			terms_of_payment VARCHAR(100),
			FOREIGN KEY (region_code) REFERENCES Region (id)
	);`)

	return err
}

func (db *DBController) createRegionTable() error {
	err := db.createTable(
		`CREATE TABLE IF NOT EXISTS Region (
			id SERIAL PRIMARY KEY UNIQUE,
			country VARCHAR(100) NOT NULL,
			city VARCHAR(100) NOT NULL,
			address VARCHAR(100) NOT NULL,
			phone_number CHAR(14) CHECK(char_length(phone_number) = 14),
			email VARCHAR(100)
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
