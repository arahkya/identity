package identity

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"

	"errors"
)

type ClientContext interface {
	CheckExistedUser(username string) bool
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Client struct {
	Context ClientContext
}

type CheckExistedUsername func(username string) bool

func CheckExistedUser(username string) bool {
	db := openDbConnection()

	prepareSelectCommand, _ := db.Prepare("SELECT count(*) FROM users WHERE username = ?")

	defer func(prepareSelectCommand *sql.Stmt) {
		err := prepareSelectCommand.Close()
		if err != nil {

		}
	}(prepareSelectCommand)

	exec := prepareSelectCommand.QueryRow(username)

	var existedUsernameCount int
	err := exec.Scan(&existedUsernameCount)
	if err != nil {
		return false
	}

	return existedUsernameCount > 0
}

func saveUser(user *UserRegisterRequest) (int, error) {
	db := openDbConnection()
	defer closeDbConnection(db)

	prepareInsCommand, _ := db.Prepare("insert into Users values(?,?,?)")

	defer func(prepareInsCommand *sql.Stmt) {
		err := prepareInsCommand.Close()
		if err != nil {

		}
	}(prepareInsCommand)

	exec, err := prepareInsCommand.Exec(user.Username, user.Password, user.Email)
	if err != nil {
		return 0, err
	}

	affected, _ := exec.RowsAffected()

	log.Println(fmt.Sprintf("Affected rows: %d", affected))

	return int(affected), nil
}

func openDbConnection() *sql.DB {
	db, dbOpenErr := sql.Open("sqlite3", "./identity.db")

	if dbOpenErr != nil {
		panic(dbOpenErr)
	}

	return db
}

func closeDbConnection(db *sql.DB) {
	err := db.Close()

	if err != nil {
		panic(err)
	}
}

func (c *Client) Register(request *UserRegisterRequest) error {

	existedUsername := c.Context.CheckExistedUser(request.Username)

	if existedUsername {
		return errors.New("username already existed")
	}

	_, err := saveUser(request)
	if err != nil {
		return err
	}

	return nil
}
