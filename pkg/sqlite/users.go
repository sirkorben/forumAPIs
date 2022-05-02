package sqlite

import (
	"database/sql"
	"errors"
	"forumAPIs/pkg/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func InsertUser(user models.User) error {
	firstName := user.FirstName
	lastName := user.LastName
	age := user.Age
	gender := user.Gender
	userName := user.Username
	email := user.Email
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	row := DB.QueryRow("select id from users where username = ?", userName)
	err = row.Scan()
	if !errors.Is(err, sql.ErrNoRows) {
		return models.ErrDuplicateUsername
	}

	row = DB.QueryRow("select id from users where email = ?", email)
	err = row.Scan()
	if !errors.Is(err, sql.ErrNoRows) {
		return models.ErrDuplicateEmail
	}

	_, err = DB.Exec("insert into users (firstname, lastname, age, gender, username, email, password, creation_date) values (?,?,?,?,?,?,?, strftime('%s','now'))",
		firstName, lastName, age, gender, userName, email, string(hashedPassword))
	if err != nil {
		log.Println("sqlite.InsertUser()", err)
		return err
	}
	return nil
}

func Authenticate(credName, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := DB.QueryRow("select id, password from users where email = ? or username = ?", credName, credName)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func GetUserProfile(id int) (*models.User, error) {
	row := DB.QueryRow("select firstname, lastname, age, gender, username, email, creation_date from users where id = ?", id)
	u := &models.User{}
	err := row.Scan(&u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Username, &u.Email, &u.CreationDate)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func GetUsernameById(id int) (*models.User, error) {
	row := DB.QueryRow("select id, username from users where id = ?", id)
	u := &models.User{}
	err := row.Scan(&u.Id, &u.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}
