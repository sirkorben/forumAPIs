package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"forumAPIs/pkg/models"
	"io"
	"net/http"
	"net/mail"
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

func validateUserData(w http.ResponseWriter, user models.User) bool {

	firstName := strings.TrimSpace(user.FirstName)
	lastName := strings.TrimSpace(user.LastName)
	age := user.Age
	gender := user.Gender
	userName := strings.TrimSpace(user.Username)
	email := strings.TrimSpace(user.Email)
	password := user.Password

	if firstName == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Firstname is missing."
		errMsg.ErrorType = "FIRSTNAME_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	if lastName == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Lastname is missing."
		errMsg.ErrorType = "LASTNAME_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	if userName == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Username is missing."
		errMsg.ErrorType = "USERNAME_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	if email == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Email is missing."
		errMsg.ErrorType = "EMAIL_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if gender == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Gender is missing."
		errMsg.ErrorType = "GENDER_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if password == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Password is missing."
		errMsg.ErrorType = "PASSWORD_FIELD_EMPTY"
		errorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if len(password) < 6 || len(password) > 20 {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Password is too short - 6 chars min."
		errMsg.ErrorType = "PASSWORD_TOO_SHORT"
		errorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}

	if age <= 0 || age > 120 {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Age is not valid"
		errMsg.ErrorType = "AGE_NOT_VALID"
		errorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}

	_, errMail := mail.ParseAddress(email)
	if errMail != nil {
		var errMsg ErrorMsg
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
