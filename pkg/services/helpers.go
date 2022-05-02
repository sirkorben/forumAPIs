package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"forumAPIs/pkg/models"
	"io"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

type ErrorMsg struct {
	ErrorDescription string `json:"error_description"`
	ErrorType        string `json:"error_type"`
}

func (mr *ErrorMsg) Error() string {
	return mr.ErrorDescription
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errDescription := "Content Type is not application/json"
		errType := "WRONG_CONTENCT_TYPE"
		return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&dst)
	if err != nil {
		// big error handling done here
		// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {

		case errors.Is(err, io.EOF):
			errDescription := "Request body must not be empty."
			errType := "REQUEST_BODY_EMPTY"
			return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			errDescription := msg
			errType := "INVALID_VALUE_FOR_FIELD"
			return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			errDescription := "Request body contains unknown field " + fieldName
			errType := "UNKNOWN_FIELD"
			return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}

		default:
			return err
		}

	}
	return nil
}

func errorResponse(w http.ResponseWriter, errMessage ErrorMsg, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(errMessage)
	w.Write(jsonResp)
}

func validateUserData(w http.ResponseWriter, user *models.User) bool {
	var errMsg ErrorMsg
	space := regexp.MustCompile(`\s+`)
	user.FirstName = space.ReplaceAllString(strings.TrimSpace(user.FirstName), " ")
	user.LastName = space.ReplaceAllString(strings.TrimSpace(user.LastName), " ")
	user.Username = space.ReplaceAllString(user.Username, "")
	user.Email = space.ReplaceAllString(strings.TrimSpace(user.Email), "")

	if user.FirstName == "" {
		errMsg.ErrorDescription = "Firstname is missing."
		errMsg.ErrorType = "FIRSTNAME_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	if user.LastName == "" {
		errMsg.ErrorDescription = "Lastname is missing."
		errMsg.ErrorType = "LASTNAME_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	if user.Username == "" {
		errMsg.ErrorDescription = "Username is missing."
		errMsg.ErrorType = "USERNAME_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	if user.Email == "" {
		errMsg.ErrorDescription = "Email is missing."
		errMsg.ErrorType = "EMAIL_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if user.Gender == "" {
		errMsg.ErrorDescription = "Gender is missing."
		errMsg.ErrorType = "GENDER_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if user.Password == "" {
		errMsg.ErrorDescription = "Password is missing."
		errMsg.ErrorType = "PASSWORD_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if len(user.Password) < 6 || len(user.Password) > 20 {
		errMsg.ErrorDescription = "Password is too short - 6 chars min."
		errMsg.ErrorType = "PASSWORD_TOO_SHORT"
		errorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}

	if user.Age <= 0 || user.Age > 120 {
		errMsg.ErrorDescription = "Age is not valid"
		errMsg.ErrorType = "AGE_NOT_VALID"
		errorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}

	_, errMail := mail.ParseAddress(user.Email)
	if errMail != nil {
		errMsg.ErrorDescription = "Email is not valid"
		errMsg.ErrorType = "EMAIL_INVALID"
		errorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}
	return true
}

func validateUserPostData(w http.ResponseWriter, post models.Post) bool {
	title := strings.TrimSpace(post.Title)
	content := strings.TrimSpace(post.Content)

	if title == "" || content == "" {
		var errMsg ErrorMsg
		if title == "" {
			errMsg.ErrorDescription = "Title field is empty"
			errMsg.ErrorType = "TITLE_FIELD_EMPTY"

		} else if content == "" {
			errMsg.ErrorDescription = "Content field is empty"
			errMsg.ErrorType = "CONTENT_FIELD_EMPTY"

		}
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
}

func validateUserMessage(w http.ResponseWriter, message *models.Message) bool {
	space := regexp.MustCompile(`\s+`)
	message.Content = space.ReplaceAllString(strings.TrimSpace(message.Content), " ")
	fmt.Println(message.Content)
	if message.Content == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Content field is empty"
		errMsg.ErrorType = "CONTENT_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
}
