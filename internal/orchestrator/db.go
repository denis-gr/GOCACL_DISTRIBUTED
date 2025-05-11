package orchestrator

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	dbConnection *sql.DB
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createTables(dbConnection *sql.DB) error {
	expressionTableQuery := `
    CREATE TABLE IF NOT EXISTS expressions (
        id TEXT PRIMARY KEY,
        expression TEXT NOT NULL,
        status TEXT NOT NULL,
        result REAL,
		creator_id TEXT NOT NULL,
		FOREIGN KEY (creator_id) REFERENCES users(id)
    );`

	userTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL
    );`

	_, err := dbConnection.Exec(userTableQuery)
	if err != nil {
		return err
	}
	_, err = dbConnection.Exec(expressionTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return &DB{dbConnection: db}, nil
}

func (db *DB) Close() error {
	if db.dbConnection != nil {
		return db.dbConnection.Close()
	}
	return nil
}

func (db *DB) CreateUser(form UserCreateForm) (*UserPublic, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	idStr := id.String()

	if len(form.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	if len(form.Password) > 70 {
		return nil, errors.New("password must be at most 70 characters long")
	}

	if len(form.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters long")
	}

	if len(form.Username) > 20 {
		return nil, errors.New("username must be at most 20 characters long")
	}

	passwordHash, err := hashPassword(form.Password)
	if err != nil {
		return nil, err
	}

	var existingUserCount int
	err = db.dbConnection.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", form.Username).Scan(&existingUserCount)
	if err != nil {
		return nil, err
	}
	if existingUserCount > 0 {
		return nil, errors.New("username already exists")
	}

	_, err = db.dbConnection.Exec("INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)",
		idStr, form.Username, passwordHash)
	if err != nil {
		return nil, err
	}

	user := &UserPublic{
		ID:       idStr,
		Username: form.Username,
	}

	return user, nil
}

func (db *DB) GetUserByUsername(username string) (UserPublic, error) {
	row := db.dbConnection.QueryRow("SELECT id, username FROM users WHERE username = ?", username)

	var user UserPublic
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserPublic{}, nil
		}
		return UserPublic{}, err
	}

	return user, nil
}

func (db *DB) GetUserAll() ([]UserPublic, error) {
	rows, err := db.dbConnection.Query("SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserPublic
	for rows.Next() {
		var user UserPublic
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) CheckUserPassword(username, password string) (bool, error) {
	row := db.dbConnection.QueryRow("SELECT password_hash FROM users WHERE username = ?", username)

	password_hash := ""
	err := row.Scan(&password_hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return checkPasswordHash(password, password_hash), nil
}

func (db *DB) CreateExpression(creatorID string, form CalculateRequest) (ExpressionDB, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return ExpressionDB{}, err
	}
	idStr := id.String()

	_, err = db.dbConnection.Exec("INSERT INTO expressions (id, expression, status, result, creator_id) VALUES (?, ?, ?, ?, ?)",
		idStr, form.Expression, "running", 0, creatorID)
	if err != nil {
		return ExpressionDB{}, err
	}

	expression := ExpressionDB{
		ID:         idStr,
		Expression: form.Expression,
		Status:     "running",
		Result:     0,
		CreatorId:  creatorID,
	}
	return expression, nil
}

func (db *DB) GetExpressionByID(id string) (ExpressionDB, error) {
	row := db.dbConnection.QueryRow("SELECT id, expression, status, result, creator_id FROM expressions WHERE id = ?", id)

	var expression ExpressionDB
	err := row.Scan(&expression.ID, &expression.Expression, &expression.Status, &expression.Result, &expression.CreatorId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ExpressionDB{}, nil
		}
		return ExpressionDB{}, err
	}

	return expression, nil
}

func (db *DB) SetResultExpression(id, status string, result float64) error {
	_, err := db.dbConnection.Exec("UPDATE expressions SET status = ?, result = ? WHERE id = ?", status, result, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetAllExpressions() ([]ExpressionDB, error) {
	rows, err := db.dbConnection.Query("SELECT id, expression, status, result, creator_id FROM expressions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []ExpressionDB
	for rows.Next() {
		var expression ExpressionDB
		err := rows.Scan(&expression.ID, &expression.Expression, &expression.Status, &expression.Result, &expression.CreatorId)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}

	return expressions, nil
}

func (db *DB) GetAllExpressionsByUserID(userID string) ([]ExpressionDB, error) {
	rows, err := db.dbConnection.Query("SELECT id, expression, status, result, creator_id FROM expressions WHERE creator_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []ExpressionDB
	for rows.Next() {
		var expression ExpressionDB
		err := rows.Scan(&expression.ID, &expression.Expression, &expression.Status, &expression.Result, &expression.CreatorId)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}

	return expressions, nil
}
