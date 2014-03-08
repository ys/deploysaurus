package deploysaurus

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DB struct {
	*sql.DB
}

func SaveUser(user DbUser) (string, error) {
	if user.Id != "" {
		return UpdateUser(user)
	} else {
		return CreateUser(user)
	}
}

func UpdateUser(user DbUser) (string, error) {
	db, err := getDB()
	if err != nil {
		log.Println(err)
		return "", err
	}
	stmt, err := db.Prepare(UpdateStatement())
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var id string
	err = stmt.QueryRow(user.Id,
		user.Email,
		user.GitHubId,
		user.GitHubLogin,
		user.GitHubToken,
		user.HerokuId,
		user.HerokuToken,
		user.HerokuRefreshToken,
		user.HerokuExpiration).Scan(&id)
	stmt.Close()
	db.Close()
	return id, err
}

func CreateUser(user DbUser) (string, error) {
	db, err := getDB()
	if err != nil {
		log.Println(err)
		return "", err
	}
	stmt, err := db.Prepare(CreateStatement())
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var id string
	err = stmt.QueryRow(user.Email,
		user.GitHubId,
		user.GitHubLogin,
		user.GitHubToken,
		user.HerokuId,
		user.HerokuToken,
		user.HerokuRefreshToken,
		user.HerokuExpiration).Scan(&id)
	stmt.Close()
	db.Close()
	return id, err
}

func CreateStatement() string {
	return `INSERT INTO users (email,
                                   github_id,
                                   github_login,
                                   github_token,
                                   heroku_id,
                                   heroku_token,
                                   heroku_refresh_token,
                                   heroku_expiration)
                                 values ($1, $2, $3, $4, $5, $6, $7, $8)
                                 RETURNING id`
}

func UpdateStatement() string {
	return `UPDATE users SET   email=$2,
                                   github_id=$3,
                                   github_login=$4,
                                   github_token=$5,
                                   heroku_id=$6,
                                   heroku_token=$7,
                                   heroku_refresh_token=$8,
                                   heroku_expiration=$9
                                   WHERE id=$1
                                   RETURNING id`
}

func GetUser(id string) (DbUser, error) {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	var u DbUser
	err = db.QueryRow(`SELECT * FROM users WHERE id=$1`, id).Scan(&u.Id,
		&u.Email,
		&u.GitHubLogin,
		&u.GitHubId,
		&u.GitHubToken,
		&u.HerokuId,
		&u.HerokuToken,
		&u.HerokuRefreshToken,
		&u.HerokuExpiration)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return u, err
}

func GetUserFromProvider(provider string, id string) (DbUser, error) {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	var u DbUser
	err = db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE %s_id=$1", provider), id).Scan(&u.Id,
		&u.Email,
		&u.GitHubLogin,
		&u.GitHubId,
		&u.GitHubToken,
		&u.HerokuId,
		&u.HerokuToken,
		&u.HerokuRefreshToken,
		&u.HerokuExpiration)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return u, err
}

func getDB() (*DB, error) {
	url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
