package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize database connection
func InitializeDB() (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite3", "./storage/dev.sqlite")
	if err != nil {
		panic(err)
	}
	// Execute PRAGMA statement to set the journal mode
	_, err = DB.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		panic(err)
	}
	// Execute PRAGMA statement to set busy timeout duration
	_, err = DB.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		panic(err)
	}

	exists, err := tableExists(DB, "users")
	if err != nil {
		return nil, err
	}

	// Check if the "users" table exists, and if not, create it
	if exists {
		log.Println("Table 'users' exists")
	} else {
		log.Println("Table 'users' does not exist, creating...")
		err = createUserTable(DB)
		if err != nil {
			return nil, err
		}
	}

	return DB, nil
}

func tableExists(db *sql.DB, tableName string) (bool, error) {
	query :=
		"SELECT name FROM sqlite_master WHERE type='table' AND name=?;"

	row := db.QueryRow(query, tableName)

	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil // Table doesn't exist
	} else if err != nil {
		return false, err // An error occurred
	}

	return true, nil // Table exists
}

func createUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY,
		access_token TEXT,
		expires INTEGER,
		scope TEXT,
		refresh_token TEXT
	);`

	_, err := db.Exec(query)
	return err
}

func AddRowUser(db *sql.DB, user_id int, accessToken string, expires int64, scope string, refreshToken string) error {
	query := `
	INSERT INTO users (user_id, access_token, expires, scope, refresh_token)
	VALUES (?, ?, ?, ?, ?);`

	_, err := db.Exec(query, user_id, accessToken, expires, scope, refreshToken)
	return err
}

// May want to add scope update as well
func UpdateRowUser(db *sql.DB, user_id int, newAccessToken string, newExpiry int64, scope string, newRefreshToken string) error {
	query := `
	UPDATE users
	SET access_token = ?, expires = ?, scope=?, refresh_token = ?
	WHERE user_id = ?`

	_, err := db.Exec(query, newAccessToken, newExpiry, scope, newRefreshToken, user_id)
	return err
}

func GetUserByID(db *sql.DB, userID int) (*User, bool, error) {
	query := `
	SELECT access_token, expires, scope, refresh_token
	FROM users
	WHERE user_id = ?`

	var user User
	err := db.QueryRow(query, userID).Scan(&user.AccessToken, &user.Expires, &user.Scope, &user.RefreshToken)
	if err == sql.ErrNoRows {
		return nil, false, nil // No user found
	} else if err != nil {
		return nil, false, err // Error occurred
	}
	user.UserID = userID
	return &user, true, nil // User found
}

type User struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}
