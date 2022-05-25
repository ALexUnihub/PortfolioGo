package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"

	"log"
	"myRedditClone/pkg/handlers"
	"myRedditClone/pkg/middleware"
	"myRedditClone/pkg/posts"
	"myRedditClone/pkg/session"
	"myRedditClone/pkg/user"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// работа с DB

	// users DB
	dsn := "root:love@tcp(localhost:3306)/golang?"
	dsn += "charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Println(err.Error())
	}

	db.SetMaxOpenConns(10)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Println("BAD PING")
		panic(err)
	}

	// post DB
	sessMongodb, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Println("sess mongodb err")
		panic(err)
	}
	//
	sessionManager := session.NewSessionsManager(db)

	postRepo := posts.NewMemoryRepo(sessMongodb)
	userRepo := user.NewMemoryRepo(db)

	userHandler := handlers.UserHandler{
		UserRepo: userRepo,
		Sessions: sessionManager,
	}

	postsHandler := handlers.PostsHandler{
		PostsRepo: postRepo,
		Sessions:  sessionManager,
	}

	// posts.AddFirstPost(postRepo)

	// posts handler
	authPosts := mux.NewRouter()
	authPosts.HandleFunc("/api/posts", postsHandler.AddNewPost).Methods("POST")

	authPostsHandler := middleware.Auth(sessionManager, authPosts)

	// post handler
	authPost := mux.NewRouter()
	authPost.HandleFunc("/api/post/{POST_ID}", postsHandler.AddNewComment).Methods("POST")
	authPost.HandleFunc("/api/post/{POST_ID}", postsHandler.ShowPost).Methods("GET")
	authPost.HandleFunc("/api/post/{POST_ID}", postsHandler.DeletePost).Methods("DELETE")

	authPostHandler := middleware.Auth(sessionManager, authPost)

	// delete post comment handler & up/down/un -vote, need auth
	authSndPost := mux.NewRouter()
	authSndPost.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postsHandler.DeleteComment).Methods("DELETE")
	authSndPost.HandleFunc("/api/post/{POST_ID}/upvote", postsHandler.UpVote).Methods("GET")
	authSndPost.HandleFunc("/api/post/{POST_ID}/downvote", postsHandler.DownVote).Methods("GET")
	authSndPost.HandleFunc("/api/post/{POST_ID}/unvote", postsHandler.UnVote).Methods("GET")

	authSndPostHandler := middleware.Auth(sessionManager, authSndPost)

	//
	r := mux.NewRouter()
	r.Handle("/api/posts", authPostsHandler)
	r.Handle("/api/post/{POST_ID}", authPostHandler)
	r.Handle("/api/post/{POST_ID}/{COMMENT_ID}", authSndPostHandler)

	r.HandleFunc("/", userHandler.Index).Methods("GET")
	r.HandleFunc("/{LOGIN}", userHandler.Index).Methods("GET")
	r.HandleFunc("/a/{CATEGORY_NAME}", userHandler.Index).Methods("GET")
	r.HandleFunc("/a/{CATEGORY_NAME}/{POST_ID}", userHandler.Index).Methods("GET")
	r.HandleFunc("/u/{USER_LOGIN}", userHandler.Index).Methods("GET")

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("../../static/"))))

	r.HandleFunc("/api/posts/", postsHandler.List).Methods("GET")
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postsHandler.ShowCategory).Methods("GET")
	r.HandleFunc("/api/user/{USER_LOGIN}", postsHandler.ShowAllUserPosts).Methods("GET")

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	rout := middleware.AccessLog(r)
	rout = middleware.Panic(rout)

	addr := ":8080"
	log.Println("server starting on addr :8080")
	http.ListenAndServe(addr, rout)
}
