package sqlite

import (
	"errors"
	"forumAPIs/pkg/models"
	"net/http"
	"time"
)

func InsertSession(token string, uID int) error {
	_, err := DB.Exec("delete from sessions where user_id = ?", uID)
	if err != nil {
		return err
	}

	_, err = DB.Exec("insert into sessions (id, user_id, created_date) values (?,?, DATETIME('now'))", token, uID)
	if err != nil {
		return err
	}

	return nil
}

func CheckSession(r *http.Request) (*models.Session, error) {
	token, err := r.Cookie("session")
	if err != nil {
		// esli net cookie to vqhodim otsjuda bez kuki kak i zahodili?
		return nil, err
	}

	session := &models.Session{}
	row := DB.QueryRow("select id, user_id, created_date from sessions where id = ?", token.Value)
	session.User = &models.User{}
	createDate := ""
	err = row.Scan(&session.Id, &session.User.Id, &createDate)
	if err != nil {
		return nil, err
	}

	session.User, err = GetUserForPostInfo(session.User.Id)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02 15:04:05", createDate)
	if err != nil {
		return nil, err
	}

	if session.Id == "" {
		return nil, errors.New("token invalid or expired")
	}

	if date.AddDate(0, 0, 1).Before(time.Now()) {
		err := DeleteSession(session.Id)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("token invalid or expired")
	}

	return session, nil
}

func DeleteSession(token string) error {
	_, err := DB.Exec("delete from sessions where id = ?", token)
	if err != nil {
		return nil
	}

	return nil
}
