package models

import (
	"log"

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

// func (db *DBController) createTables() {
// 	DBLogger := logger.WithFields(logrus.Fields{"action": "CreateDBTable"})

// 	if err := db.createUserInfoTable(); err != nil {
// 		log.Fatalln("create user information table:", err)
// 	}

// 	if err := db.createAuthInfoTable(); err != nil {
// 		log.Fatalln("create authtorization information table:", err)
// 	}

// 	if err := db.createProviderTable(); err != nil {
// 		DBLogger.Error("Invalid provider table: ", err)
// 		os.Exit(1)
// 	}
// }

// func (db *DBController) createUserInfoTable() error {
// 	err := db.createTable(
// 		`CREATE TABLE IF NOT EXISTS UserInformation (
// 			id 		SERIAL PRIMARY KEY,
// 			surname VARCHAR(50) NOT NULL,
// 			name 	VARCHAR(50) NOT NULL
// 		);`)

// 	return err
// }

// func (db *DBController) createProviderTable() error {
// 	err := db.createTable(
// 		`CREATE TABLE IF NOT EXISTS Provider (
// 			vendor_code		 SERIAL PRIMARY KEY UNIQUE,
// 			name			 VARCHAR(200) NOT NULL,
// 			unp				 VARCHAR(10) NOT NULL CHECK(char_length(unp) = 9),
// 			terms_of_payment VARCHAR(100),
// 			address			 VARCHAR(200) NOT NULL,
// 			phone_number	 CHAR(14) CHECK(char_length(phone_number) = 13),
// 			email			 VARCHAR(100),
// 			web_site		 VARCHAR(100)
// 	);`)

// 	return err
// }
