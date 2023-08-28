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

func AddSpotifyToken(db *sql.DB, user_id int, accessToken string, expires int64, scope string, refreshToken string) error {
	query := `
	INSERT INTO users (user_id, access_token, expires, scope, refresh_token)
	VALUES (?, ?, ?, ?, ?);`

	_, err := db.Exec(query, user_id, accessToken, expires, scope, refreshToken)
	return err
}

type SpotifyToken struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func GetSpotifyToken(db *sql.DB, user_id int) (SpotifyToken, error) {
	var token SpotifyToken

	query := `SELECT user_id, access_token, expires, scope, refresh_token
	          FROM users
	          WHERE user_id = ?`

	row := db.QueryRow(query, user_id)
	err := row.Scan(&token.UserID, &token.AccessToken, &token.Expires, &token.Scope, &token.RefreshToken)
	if err != nil {
		return SpotifyToken{}, err
	}
	return token, nil
}

// May want to add scope update as well
func UpdateSpotifyToken(db *sql.DB, user_id int, newAccessToken string, newExpiry int64, newRefreshToken string) error {
	query := `
	UPDATE users
	SET access_token = ?, expires = ?, refresh_token = ?
	WHERE user_id = ?`

	_, err := db.Exec(query, newAccessToken, newExpiry, newRefreshToken, user_id)
	return err
}
