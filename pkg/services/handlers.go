package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"forumAPIs/pkg/models"
	"forumAPIs/pkg/sqlite"
	"log"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var internalErrorMsg = ErrorMsg{
	ErrorDescription: "Internal server error",
	ErrorType:        "INTERNAL_SERVER_ERROR",
}

var NotFoundErrorMsg = ErrorMsg{
	ErrorDescription: "Page not found",
	ErrorType:        "NOT_FOUND_ERROR",
}

var UnauthorizedErrorMsg = ErrorMsg{
	ErrorDescription: "Restricted, becouse of non authorization",
	ErrorType:        "UNAUTHORIZED_ERROR",
}
var EmptyCategoryErrorMsg = ErrorMsg{
	ErrorDescription: "Category is empty",
	ErrorType:        "EMPTY_CATEGORY",
}

var EmptyChatErrorMsg = ErrorMsg{
	ErrorDescription: "NO_STARTED_CHATS",
	ErrorType:        "EMPTY_CHATS",
}

func enableCors(w *http.ResponseWriter) {
	// https://www.stackhawk.com/blog/golang-cors-guide-what-it-is-and-how-to-enable-it/
	(*w).Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500") // address for live server
}

func home(w http.ResponseWriter, r *http.Request) {
	// will be used to serve opportunity to signin or signup
	enableCors(&w)
	if r.URL.Path != "/" {
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		var p models.Post
		err := decodeJSONBody(w, r, &p)
		if err != nil {
			var errMsg *ErrorMsg
			if errors.As(err, &errMsg) {
				errorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				log.Println(err.Error())
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		if validateUserPostData(w, p) {
			var categories []string
			for i := range p.Categories {
				if p.Categories[i].Name != "" {
					categories = append(categories, p.Categories[i].Name)
				} else {
					categories = append(categories, "Misc")
				}
			}
			fmt.Println(p.Title, p.Content, categories, s.User.Id)
			_, err := sqlite.InsertPost(p.Title, p.Content, categories, s.User.Id)
			if err != nil {
				log.Println(err.Error())
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
				return
			}
		}
	}
}

func categories(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	categories, err := sqlite.GetAllCategories()
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(categories)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func postsByCategoryId(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	catId, err := strconv.Atoi(r.URL.Path[10:])
	if err != nil || catId < 1 {
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	p, err := sqlite.GetPostsByCategory(catId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			//  choose http response status for no posts under category????
			errorResponse(w, EmptyCategoryErrorMsg, http.StatusBadRequest)
		} else {
			errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(p)
	if err != nil {
		// is internal error msg needed here?
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func post(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	user := &models.User{
		Id: -1,
	}
	s, authErr := sqlite.CheckSession(r)
	if authErr != nil {
		log.Println(authErr.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	} else {
		user = s.User

	}
	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // /post/1/ -> [post, 1]
	if len(url) == 1 {
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	postId, err := strconv.Atoi(url[1])
	if err != nil {
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		return
	}

	if len(url) == 2 { // /post/%id%/
		postById(w, r, postId)
		return
	}

	if len(url) == 3 { // /post/%id%/%something%
		if r.Method == http.MethodPost {
			var pr models.PostReaction
			err := decodeJSONBody(w, r, &pr)
			if err != nil {
				var errMsg *ErrorMsg
				if errors.As(err, &errMsg) {
					errorResponse(w, *errMsg, http.StatusBadRequest)

				} else {
					log.Println(err.Error())
					errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
				}
				return
			}
			if url[2] == "comment" && pr.PostComment != "" {
				sqlite.InsertComment(postId, pr.PostComment, user.Id)
				return
			}
			if url[2] == "like" && pr.PostLikeDislike == "like" {
				sqlite.ChangePostLike(postId, user.Id)
				return
			}
			if url[2] == "dislike" && pr.PostLikeDislike == "dislike" {
				sqlite.ChangePostDislike(postId, user.Id)
				return
			}
		}
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
	}
}

func posts(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	posts, err := sqlite.GetPosts()
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(posts)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func postById(w http.ResponseWriter, r *http.Request, postId int) {
	//auth check do not needed here, it`s being checked upper in post function
	enableCors(&w)
	if postId < 1 {
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	post, err := sqlite.GetPostById(postId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		} else {
			errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		}
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(post)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodPost {
		var u models.User
		err := decodeJSONBody(w, r, &u)
		if err != nil {
			var errMsg *ErrorMsg
			if errors.As(err, &errMsg) {
				errorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				log.Println(err.Error())
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		if validateUserData(w, &u) {
			err := sqlite.InsertUser(u)
			if err != nil {
				var errMsg ErrorMsg
				if errors.Is(err, models.ErrDuplicateUsername) {
					errMsg.ErrorDescription = "Username already taken."
					errMsg.ErrorType = "USERNAME_ALREADY_TAKEN"
					errorResponse(w, errMsg, http.StatusUnsupportedMediaType)
					return
				}
				if errors.Is(err, models.ErrTooManySpaces) {
					errMsg.ErrorDescription = "too many spaces"
					errMsg.ErrorType = "DOUBLE_SPACES_IF_FIELDS"
					errorResponse(w, errMsg, http.StatusUnsupportedMediaType)
					return
				}
				if errors.Is(err, models.ErrDuplicateEmail) {
					errMsg.ErrorDescription = "Email already taken."
					errMsg.ErrorType = "EMAIL_ALREADY_TAKEN"
					errorResponse(w, errMsg, http.StatusNotAcceptable)
					return
				}
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
				return
			} else {
				log.Println("User inserted - ", u.Username)
			}
		}
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodPost {
		var u models.User
		err := decodeJSONBody(w, r, &u)
		if err != nil {
			var errMsg *ErrorMsg
			if errors.As(err, &errMsg) {
				errorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				log.Println(err.Error())
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		var credential string
		if u.Email == "" {
			credential = u.Username
		} else {
			credential = u.Email
		}
		id, err := sqlite.Authenticate(credential, u.Password)
		if err != nil {
			var errMsg ErrorMsg
			if errors.Is(err, models.ErrInvalidCredentials) {
				errMsg.ErrorDescription = "Email/username and password don't match."
				errMsg.ErrorType = "CREDENTIALS_DONT_MATCH"
				errorResponse(w, errMsg, http.StatusBadRequest)
			} else {
				log.Println(err.Error())
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
			}
			return
		}

		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:   "session",
			Value:  sID.String(),
			MaxAge: 60 * 60 * 24,
		}
		http.SetCookie(w, c)

		err = sqlite.InsertSession(c.Value, id)
		if err != nil {
			log.Println(err.Error())
			errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func signOut(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	err = sqlite.DeleteSession(s.Id)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	})
}

func getMyProfile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)

		return
	}
	var u *models.User
	u, err = sqlite.GetUserProfile(s.User.Id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		} else {
			errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		}
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func getOtherUserProfile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	userId, err := strconv.Atoi(r.URL.Path[6:])
	if err != nil || userId < 1 {
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	var u *models.User
	u, err = sqlite.GetUserProfile(userId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			var errMsg ErrorMsg
			errMsg.ErrorDescription = "User not found"
			errMsg.ErrorType = "STATUS_BAD_REQUEST"
			errorResponse(w, errMsg, http.StatusBadRequest)
		} else {
			errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		}
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // /post/1/ -> [post, 1]
	if len(url) == 1 && url[0] == "chat" {
		getAllUserChats(w, r)
		return
	}
	if len(url) == 2 {
		chatPersonId, err := strconv.Atoi(url[1])
		if err != nil {
			errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
			return
		}
		// show chat window with userId
		// show all messages for now later will be only 10 ??
		getChatWithUser(w, r, s.User.Id, chatPersonId)
		return
	}

	if r.Method == http.MethodPost && url[2] == "message" {
		var message models.Message
		err := decodeJSONBody(w, r, &message)
		if err != nil {
			var errMsg *ErrorMsg
			if errors.As(err, &errMsg) {
				errorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				log.Println(err.Error())
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		if validateUserMessage(w, &message) {
			recevierId, err := strconv.Atoi(url[1])
			if err != nil {
				errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
				return
			}
			err = sqlite.SendMessage(message.Content, s.User.Id, recevierId)
			if err != nil {
				log.Println("sqlite.SendMessage()", err)
				errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
			}
		}
	}
}

func getAllUserChats(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		log.Println(err.Error())
		errorResponse(w, UnauthorizedErrorMsg, http.StatusUnauthorized)

		return
	}
	usersChattedWIth, err := sqlite.GetAllChats(s.User.Id)
	if err != nil {
		// if no chats started what to response? and choose what status?
		errorResponse(w, EmptyChatErrorMsg, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(usersChattedWIth)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func getChatWithUser(w http.ResponseWriter, r *http.Request, userId, chatPersonId int) {
	enableCors(&w)
	if chatPersonId < 1 {
		errorResponse(w, NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	messages, err := sqlite.GetChatWithUser(userId, chatPersonId)
	if err != nil {
		log.Println("sqlite.GetChatWithUser()", err)
		if errors.Is(err, models.ErrNoRecord) {
			var errMsg ErrorMsg
			errMsg.ErrorDescription = "Messages not found"
			errMsg.ErrorType = "STATUS_BAD_REQUEST"
			errorResponse(w, errMsg, http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(messages)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		errorResponse(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
