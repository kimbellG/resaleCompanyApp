package models

import "log"

type RegistrationInformation struct {
	AuthInfo     AuthorizationUserInformation `json:"auth_info"`
	PersonalInfo UserInfo                     `json:"personal_info"`
}

type AuthorizationUserInformation struct {
	Login      string
	Password   string
	Status     int
	UserInfoId int64 `json:"user_info_id"`
}

const (
	default_user = iota
	admin
)

type UserInfo struct {
	Surname string
	Name    string
}

func (db *DBController) RegisterUser(userInfo *RegistrationInformation) error {

	user_id, err := db.AddPersonalUserRecord(&userInfo.PersonalInfo)
	if err != nil {
		log.Println("add personal user info in db:", err)
		return err
	}

	userInfo.AuthInfo.UserInfoId = user_id

	if err := db.AddAuthRecord(&userInfo.AuthInfo); err != nil {
		log.Println("add auth info user db:", err)
		return err
	} else {
		log.Println("resigtration is successfull")
	}

	return nil
}

func (db *DBController) AddAuthRecord(authInfo *AuthorizationUserInformation) error {
	stmt, err := db.Prepare("INSERT INTO authtorizationinformation (login, password, status, userinfo) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(authInfo.Login, authInfo.Password, authInfo.Status, authInfo.UserInfoId); err != nil {
		return err
	}

	return nil
}

func (db *DBController) AddPersonalUserRecord(personalInfo *UserInfo) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO userinformation (surname, name) VALUES($1, $2) RETURNING id")
	if err != nil {
		return -1, err
	}

	rows_id, err := stmt.Query(personalInfo.Surname, personalInfo.Name)
	if err != nil {
		return -1, err
	}
	defer rows_id.Close()

	var id int64 = 0
	for rows_id.Next() {
		if err := rows_id.Scan(&id); err != nil {
			return -1, err
		}
	}

	return id, nil
}
