package services

import (
	"flag"
	"forumAPIs/pkg/sqlite"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Server() {
	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  routes(),
	}

	sqlite.DataBase()
	defer sqlite.DB.Close()

	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	// fs := http.FileServer(http.Dir("./views/"))
	// mux.Handle("/", http.StripPrefix("", fs))
	mux.HandleFunc("/ws", webSocket)

	mux.HandleFunc("/", home)
	mux.HandleFunc("/post/", post) // --> /post/{id} /post/{id}/comment /post/{id}/like /post/{id}/dislike
	mux.HandleFunc("/post/all", posts)
	mux.HandleFunc("/post/create", createPost)
	mux.HandleFunc("/categories", categories)
	mux.HandleFunc("/category/", postsByCategoryId)
	mux.HandleFunc("/signup", signUp)
	mux.HandleFunc("/signin", signIn)
	mux.HandleFunc("/signout", signOut)
	mux.HandleFunc("/me", getMyProfile)
	mux.HandleFunc("/user/", getOtherUserProfile)
	mux.HandleFunc("/chat/", chatHandler) // /chat   /chat/{userID}    /chat/{userId}/message
	return mux
}
