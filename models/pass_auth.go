package models

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type PasswordAutheficationInfo struct {
	Login, Password string
}

func (db *DBController) GetUserInfo(logpass *PasswordAutheficationInfo) (*AuthorizationUserInformation, error) {
	DBLog := logger.WithFields(log.Fields{"action": "PasswordAuthtorization", "func": "GetUserFromDB"})

	result := &AuthorizationUserInformation{}
	stmt, err := db.Prepare("SELECT login, password, status, userinfo FROM authtorizationinformation WHERE login = $1 AND password = $2")
	if err != nil {
		DBLog.Error("Invalid stmt to db")
		os.Exit(1)
	}

	userAuthInfo, err := stmt.Query(logpass.Login, logpass.Password)
	if err != nil {
		DBLog.Error("Invalid query to database")
		os.Exit(1)
	}

	defer userAuthInfo.Close()

	for userAuthInfo.Next() {
		if err := userAuthInfo.Scan(&result.Login, &result.Password, &result.Status, &result.UserInfoId); err != nil {
			DBLog.Error(fmt.Sprintf("Error with scan from query: %v"), err)
			os.Exit(1)
		}
	}

	if result.Login == "" {
		return nil, errors.New("Authorization failed. User doesn't exist.")
	}

	return result, nil
}
