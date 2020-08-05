package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"nathanielwheeler.com/controllers"
	"nathanielwheeler.com/middleware"
	"nathanielwheeler.com/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const port = ":3000"

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	// Start up database connection
	dbEnv := getDBEnv()
	psqlConnectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable",
		dbEnv.host, dbEnv.port, dbEnv.user, dbEnv.password, dbEnv.name,
	)

	// Initialize services
	services, err := models.NewServices(psqlConnectionStr)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// Router Initilization
	r := mux.NewRouter()

	// Initialize controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	postsC := controllers.NewPosts(services.Posts, services.Images, r)

	// Middleware
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	// Statics Routes
	r.Handle("/",
		staticC.Home).
		Methods("GET")
	r.Handle("/resume",
		staticC.Resume).
		Methods("GET")

	// User Routes
	r.HandleFunc("/register",
		usersC.Registration).
		Methods("GET")
	r.HandleFunc("/register",
		usersC.Register).
		Methods("POST")
	r.Handle("/login",
		usersC.LoginView).
		Methods("GET")
	r.HandleFunc("/login",
		usersC.Login).
		Methods("POST")
	r.HandleFunc("/cookietest",
		usersC.CookieTest).
		Methods("GET")

	// Post Routes
	r.HandleFunc("/posts",
		requireUserMw.ApplyFn(postsC.Create)).
		Methods("POST")
	// FIXME: I can't figure out why Index will render for GET "/posts/index" but not GET "/posts"
	r.HandleFunc("/posts/index",
		postsC.Index).
		Methods("GET").
		Name(controllers.IndexPosts)
	r.Handle("/posts/new",
		requireUserMw.Apply(postsC.New)).
		Methods("GET")
	r.HandleFunc("/posts/{year:20[0-9]{2}}/{title}",
		postsC.Show).
		Methods("GET").
		Name(controllers.ShowPost)
	r.HandleFunc("/posts/{year:20[0-9]{2}}/{title}/edit",
		requireUserMw.ApplyFn(postsC.Edit)).
		Methods("GET").
		Name(controllers.EditPost)
	r.HandleFunc("/posts/{year:20[0-9]{2}}/{title}/update",
		requireUserMw.ApplyFn(postsC.Update)).
		Methods("POST")
	r.HandleFunc("/posts/{year:20[0-9]{2}}/{title}/upload",
		requireUserMw.ApplyFn(postsC.Upload)).
		Methods("POST")
	r.HandleFunc("/posts/{year:20[0-9]{2}}/{title}/delete",
		requireUserMw.ApplyFn(postsC.Delete)).
		Methods("POST")

	// Start that server!
	fmt.Println("Now listening on", port)
	http.ListenAndServe(port, userMw.Apply(r))
}

// #region DB HELPERS

type dbEnv struct {
	host, user, password, port, name string
}

func getDBEnv() dbEnv {
	return dbEnv{
		host:     checkDBEnv("host"),
		user:     checkDBEnv("user"),
		password: checkDBEnv("password"),
		port:     checkDBEnv("port"),
		name:     checkDBEnv("name"),
	}
}

func checkDBEnv(str string) string {
	str, exists := os.LookupEnv("DB_" + strings.ToUpper(str))
	if !exists {
		panic(".env is missing environment variable: '" + str + "'")
	}
	return str
}

// #endregion
