package postgres

import (
	"context"
	"cw/dbutil"
	"cw/dbutil/condition"
	"cw/models"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type OutputUser struct {
	Id       int
	Login    string
	Password string
	Status   bool
	Access   string
	Name     string
}

type AdminRepo struct {
	db                *sql.DB
	dbWriteController *dbutil.DBController
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{
		db:                db,
		dbWriteController: dbutil.NewAddController(db, "userInformation"),
	}
}

func (a *AdminRepo) OnOffUser(ctx context.Context, username string, status bool) error {

	cond := createUsernameCondition()
	if err := a.dbWriteController.Update([]string{status_a}, cond, status, username); err != nil {
		return fmt.Errorf("update status in db: %v", err)
	}

	return nil
}

func createUsernameCondition() *condition.Condition {
	cond := condition.NewCondition()
	cond.AddCondition(condition.NOTHING, login_a, condition.EQ)

	return cond
}

func (a *AdminRepo) SetAccessProfile(ctx context.Context, username string, access string) error {
	cond := createUsernameCondition()

	if err := a.dbWriteController.Update([]string{access_a}, cond, access, username); err != nil {
		return fmt.Errorf("update in db: %v", err)
	}

	return nil
}

func (a *AdminRepo) GetUser(ctx context.Context, username string) (*models.User, error) {
	cond := createUsernameCondition()

	resultRows, err := a.dbWriteController.Select("*", cond, username)
	if err != nil {
		return nil, fmt.Errorf("select user: %v", err)
	}
	defer resultRows.Close()

	result := &OutputUser{}
	for resultRows.Next() {
		if err := resultRows.Scan(&result.Id, &result.Login, &result.Password, &result.Status, &result.Access, &result.Name); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}
	}

	return outputToModels(result), nil
}

func (a *AdminRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	resultRows, err := a.dbWriteController.SelectAllTable("login, password, status, access, name")
	if err != nil {
		return nil, fmt.Errorf("select from db: %v", err)
	}

	result := []*models.User{}
	for resultRows.Next() {
		tmp := &models.User{}
		if err := resultRows.Scan(&tmp.Login, &tmp.Password, &tmp.Status, &tmp.Access, &tmp.Name); err != nil {
			return nil, fmt.Errorf("scan user: %v", err)
		}

		result = append(result, tmp)
	}

	return result, nil
}

func outputToModels(out *OutputUser) *models.User {
	return &models.User{
		Login:    out.Login,
		Password: out.Password,
		Status:   out.Status,
		Access:   out.Access,
		Name:     out.Name,
	}
}

func (a *AdminRepo) UpdateUser(ctx context.Context, username string, key, value string) error {
	cond := createUsernameCondition()

	if err := a.dbWriteController.Update([]string{key}, cond, value, username); err != nil {
		return err
	}

	return nil
}

func (a *AdminRepo) DeleteUser(ctx context.Context, username string) error {
	cond := createUsernameCondition()

	if err := a.dbWriteController.Delete(cond, username); err != nil {
		return err
	}

	return nil
}
