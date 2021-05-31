package postgres

import (
	"context"
	"cw/dbutil"
	"cw/models"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type RangPostgres struct {
	db             *sql.DB
	alternativeAdd *dbutil.DBController
}

func NewRangPostgres(lib_db *sql.DB) *RangPostgres {
	err := dbutil.Create(lib_db,
		`CREATE TABLE IF NOT EXISTS Problems (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE,
		description VARCHAR(1000)
	);`)

	if err != nil {
		panic(fmt.Errorf("problem table: %v", err))
	}

	err = dbutil.Create(lib_db,
		`CREATE TABLE IF NOT EXISTS Alternatives (
			id SERIAL PRIMARY KEY,
			problemId INT REFERENCES Problems(id) ON DELETE CASCADE,
			name VARCHAR(100),
			description VARCHAR(1000),
			weight DOUBLE PRECISION DEFAULT 0
		);`)

	if err != nil {
		panic(fmt.Errorf("alternative table: %v", err))
	}

	err = dbutil.Create(lib_db,
		`CREATE TABLE IF NOT EXISTS Marks (
			id SERIAL PRIMARY KEY,
			problemId INT REFERENCES Problems (id) ON DELETE CASCADE,
			alternativeId INT REFERENCES Alternatives (id) ON DELETE CASCADE,
			expertID INT REFERENCES userInformation(id) ON DELETE CASCADE,
			mark DOUBLE PRECISION DEFAULT 0,
			weight DOUBLE PRECISION DEFAULT 0
		);`)

	if err != nil {
		panic(fmt.Errorf("marks table: %v", err))
	}

	err = dbutil.Create(lib_db, "CREATE UNIQUE INDEX IF NOT EXISTS unq_mark ON Marks(problemId, alternativeId, expertId)")
	if err != nil {
		panic(fmt.Errorf("create index: %v", err))
	}

	return &RangPostgres{
		db:             lib_db,
		alternativeAdd: dbutil.NewAddController(lib_db, "Alternatives"),
	}
}

func (r *RangPostgres) Add(ctx context.Context, problem *models.ProblemInput) error {
	prodlemID := int(0)
	if err := r.queryRow("INSERT INTO Problems (name, description) VALUES ($1, $2) RETURNING id", &prodlemID, problem.Name, problem.Description); err != nil {
		return fmt.Errorf("insert problem info: %v", err)
	}

	for _, alternative := range problem.Alternatives {
		if err := r.exec("INSERT INTO Alternatives (name, problemId, description) VALUES ($1, $2, $3)", alternative.Name, prodlemID, alternative.Description); err != nil {
			return fmt.Errorf("insert alternative %v: %v", alternative.Name, err)
		}

	}

	return nil
}

func (r *RangPostgres) exec(query string, argv ...interface{}) error {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	if _, err := stmt.Exec(argv...); err != nil {
		return fmt.Errorf("exec stmt: %v", err)
	}

	return nil
}

func (r *RangPostgres) queryRow(query string, result interface{}, arg ...interface{}) error {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	if err := stmt.QueryRow(arg...).Scan(result); err != nil {
		return fmt.Errorf("query row: %v", err)
	}

	return nil
}

func (r *RangPostgres) AddAlternativeMark(ctx context.Context, problem *models.Problem) error {
	for _, alternative := range problem.Alternatives {
		if err := r.exec("UPDATE Alternatives SET weight=$1 WHERE id = $2", alternative.Weight, alternative.Id); err != nil {
			return fmt.Errorf("update alternative weight in %v: %v", alternative.Name, err)
		}
		if err := r.insertMarks(problem.Id, alternative); err != nil {
			return fmt.Errorf("insert marks: %v", err)
		}
	}

	return nil
}

func (r *RangPostgres) insertMarks(problemId int, alternative *models.Alternative) error {
	for key, value := range alternative.Marks {
		weight, ok := alternative.SpecificWeights[key]
		if !ok {
			weight = 0
		}
		if err := r.exec(`INSERT INTO Marks (problemId, alternativeId, expertID, mark, weight) VALUES ($1, $2, $3, $4, $5)
						  ON CONFLICT (problemId, alternativeId, expertId) DO UPDATE SET mark=$4, weight=$5`,
			problemId, alternative.Id, key, value, weight); err != nil {
			return err
		}
	}

	return nil
}

func (r *RangPostgres) AddAlternative(ctx context.Context, problemId int, alternative *models.AlternativeInput) error {
	if err := r.alternativeAdd.Add("name, problemId, description", alternative.Name, problemId, alternative.Description); err != nil {
		return fmt.Errorf("add alternative: %v", err)
	}

	return nil
}

func (r *RangPostgres) Gets(ctx context.Context) ([]*models.Problem, error) {
	problem := *new([]*models.Problem)
	stmt, err := r.db.Prepare("SELECT * FROM Problems")
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query stmt: %v", err)
	}

	for rows.Next() {
		tmp := &models.Problem{}
		if err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Description); err != nil {
			return nil, fmt.Errorf("scan stmt: %v", err)
		}

		tmp.Alternatives, err = r.getAlternative(tmp.Id)
		if err != nil {
			return nil, fmt.Errorf("get alternatives: %v", err)
		}

		problem = append(problem, tmp)
	}

	return problem, nil
}

func (r *RangPostgres) getAlternatives(problems []*models.Problem) ([]*models.Problem, error) {
	var err error
	for _, problem := range problems {
		problem.Alternatives, err = r.getAlternative(problem.Id)
		if err != nil {
			return nil, fmt.Errorf("get alternative: %v", err)
		}
	}

	return problems, nil
}

func (r *RangPostgres) getAlternative(problemID int) ([]*models.Alternative, error) {
	stmt, err := r.db.Prepare("SELECT id, name, description, weight FROM Alternatives WHERE ProblemId = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	rows, err := stmt.Query(problemID)
	if err != nil {
		return nil, fmt.Errorf("query stmt: %v", err)
	}

	result := make([]*models.Alternative, 0)
	for rows.Next() {
		tmp := new(models.Alternative)
		if err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Description, &tmp.Weight); err != nil {
			return nil, fmt.Errorf("scan element: %v", err)
		}

		tmp.Marks, tmp.SpecificWeights, err = r.getMarks(problemID, tmp.Id)
		if err != nil {
			return nil, fmt.Errorf("get marks: %v", err)
		}

		result = append(result, tmp)
	}

	return result, nil
}

func (r *RangPostgres) getMarks(problemId, alternativeId int) (map[int]float32, map[int]float32, error) {
	stmt, err := r.db.Prepare("SELECT expertID, mark, weight FROM Marks WHERE problemId = $1 AND alternativeId = $2")
	if err != nil {
		return nil, nil, fmt.Errorf("prepare stmt: %v", err)
	}

	rows, err := stmt.Query(problemId, alternativeId)
	if err != nil {
		return nil, nil, fmt.Errorf("query stmt: %v", err)
	}

	marks := make(map[int]float32)
	specificWeight := make(map[int]float32)

	for rows.Next() {
		expertID := 0
		mark := float32(0)
		weight := float32(0)
		if err := rows.Scan(&expertID, &mark, &weight); err != nil {
			return nil, nil, fmt.Errorf("scan rows: %v", err)
		}
		marks[expertID] = mark
		specificWeight[expertID] = weight
	}

	return marks, specificWeight, nil
}

func (r *RangPostgres) GetProblemReport(ctx context.Context, id int) (*models.Problem, error) {
	stmt, err := r.db.Prepare("SELECT * FROM Problems WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	result := &models.Problem{}
	if err := stmt.QueryRow(id).Scan(&result.Id, &result.Name, &result.Description); err != nil {
		return nil, fmt.Errorf("scan element: %v", err)
	}

	result.Alternatives, err = r.getAlternative(result.Id)
	if err != nil {
		return nil, fmt.Errorf("get alternatives: %v", err)
	}

	return result, nil
}
