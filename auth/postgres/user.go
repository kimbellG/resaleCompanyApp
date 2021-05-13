package postgres

import (
	"context"
	"cw/logger"
	"cw/models"
	"errors"
	"fmt"
	"os"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
)

type User struct {
	ID       int
	Login    string
	Password string
	Status   bool
	Access   string
	Name     string
}

func ModelsToPostgres(user *models.User) *User {
	return &User{
		Login:    user.Login,
		Password: user.Password,
		Status:   user.Status,
		Access:   user.Access,
		Name:     user.Login,
	}
}

func PostgresToModels(user *User) *models.User {
	return &models.User{
		Login:    user.Login,
		Password: user.Password,
		Status:   user.Status,
		Access:   user.Access,
		Name:     user.Name,
	}
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	err := createTable(db,
		`CREATE TABLE IF NOT EXISTS userInformation (
	id 		 SERIAL PRIMARY KEY UNIQUE,
 	login	 VARCHAR(30) NOT NULL UNIQUE, 
 	password VARCHAR(30) NOT NULL,
 	status	 BOOLEAN,
	access	 VARCHAR(100) NOT NULL,
	name	 VARCHAR(200) NOT NULL
);`)

	if err != nil {
		DBlog := logger.NewLoggerWithFields(map[string]interface{}{"database": "Postgresql", "action": "CreateDBTable"})
		DBlog.Errorf("create user table: %v", err)
		os.Exit(1)
	}

	return &UserRepository{
		db: db,
	}
}

func createTable(db *sql.DB, query string) error {
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	dbuser := ModelsToPostgres(user)
	if err := r.addAuthRecord(dbuser); err != nil {
		return err
	}

	return nil

}

func (r *UserRepository) addAuthRecord(user *User) error {
	stmt, err := r.db.Prepare("INSERT INTO userInformation (login, password, status, access, name) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		logger.AssertMessage(map[string]interface{}{"object": "postgres", "action": "insert user"}, fmt.Sprintf("stmt is invalid: %v", err))
	}

	if _, err := stmt.Exec(user.Login, user.Password, user.Status, user.Access, user.Name); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	dbUser := r.requestUserToDB(username, password)
	if dbUser.Login == "" {
		return nil, errors.New("Authorization failed. User doesn't exist.")
	}

	return PostgresToModels(dbUser), nil
}

func (r *UserRepository) requestUserToDB(username, password string) *User {
	fields := map[string]interface{}{"object": "postgers", "action": "select user"}

	result := &User{}
	stmt, err := r.db.Prepare("SELECT login, password, status, access, name FROM userInformation WHERE login = $1 AND password = $2")
	if err != nil {
		logger.AssertMessage(fields, fmt.Sprintf("stmt is invalid: %v", err))
	}

	userAuthInfo, err := stmt.Query(username, password)
	if err != nil {
		logger.AssertMessage(fields, fmt.Sprintf("query is invalid: %v", err))
	}

	defer userAuthInfo.Close()

	for userAuthInfo.Next() {
		if err := userAuthInfo.Scan(&result.Login, &result.Password, &result.Status, &result.Access, &result.Name); err != nil {
			logger.AssertMessage(fields, fmt.Sprintf("scan is invalid: %v", err))
		}
	}

	return result
}

func (r *UserRepository) GetNameByLogin(login string) (string, error) {
	stmt, err := r.db.Prepare("SELECT name FROM userInformation WHERE login = $1")
	if err != nil {
		return "", fmt.Errorf("prepare stmt: %v", err)
	}

	name := ""
	if err := stmt.QueryRow(login).Scan(&name); err != nil {
		return "", fmt.Errorf("scan name: %v", err)
	}

	return name, nil
}

func (r *UserRepository) GetIdByLogin(login string) (int, error) {
	stmt, err := r.db.Prepare("SELECT id FROM userInformation WHERE login = $1")
	if err != nil {
		return -1, fmt.Errorf("prepare stmt: %v", err)
	}

	id := 0
	if err := stmt.QueryRow(login).Scan(&id); err != nil {
		return -1, fmt.Errorf("scan id: %v", err)
	}

	return id, nil
}
