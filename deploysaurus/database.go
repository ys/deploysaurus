package deploysaurus

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DB interface {
	SaveUser(user DbUser) (string, error)
	CreateUser(user DbUser) (string, error)
	UpdateUser(user DbUser) (string, error)
	GetUser(id string) (DbUser, error)
	GetUserFromProvider(provider string, id string) (DbUser, error)
	GetUsersCount() (int, error)
}

type DeployDB struct {
	*sql.DB
}

func (db DeployDB) SaveUser(user DbUser) (string, error) {
	if user.Id != "" {
		return db.UpdateUser(user)
	} else {
		return db.CreateUser(user)
	}
}

func (db DeployDB) CreateUser(user DbUser) (string, error) {
	stmt, err := db.Prepare(createStatement())
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
	return id, err
}

func (db DeployDB) UpdateUser(user DbUser) (string, error) {
	stmt, err := db.Prepare(updateStatement())
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
	return id, err
}

func (db DeployDB) GetUser(id string) (DbUser, error) {
	var u DbUser
	err := db.QueryRow(`SELECT * FROM users WHERE id=$1`, id).Scan(&u.Id,
		&u.Email,
		&u.GitHubLogin,
		&u.GitHubId,
		&u.GitHubToken,
		&u.HerokuId,
		&u.HerokuToken,
		&u.HerokuRefreshToken,
		&u.HerokuExpiration)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return u, err
}

func (db DeployDB) GetUserFromProvider(provider string, id string) (DbUser, error) {
	var u DbUser
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE %s_id=$1", provider), id).Scan(&u.Id,
		&u.Email,
		&u.GitHubLogin,
		&u.GitHubId,
		&u.GitHubToken,
		&u.HerokuId,
		&u.HerokuToken,
		&u.HerokuRefreshToken,
		&u.HerokuExpiration)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return u, err
}

func (db DeployDB) GetUsersCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count, err

}

var GetDB = func() (DB, error) {
	url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return DeployDB{db}, nil
}

func createStatement() string {
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

func updateStatement() string {
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
