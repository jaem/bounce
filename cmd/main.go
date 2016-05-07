package main

import (
	"fmt"
	"net/http"

	"github.com/jaem/bounce"
	"github.com/jaem/bounce/providers/local"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jaem/nimble"
	"github.com/jaem/bounce/storage/jwt"

)

func main() {
	//trial.TestRouters()
	theMain()
}

func theMain() {

	bou := bounce.New(jwt.NewIdManager())
	bou.Register("local", local.NewProvider())
	bou.Register("local2", local.NewProvider())

	nim := nimble.Default()
	nim.UseHandlerFunc(bou.IdentifyRequest)
	//nim.UseFunc(middlewareA)
	//nim.UseFunc(middlewareB)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/", helloHandlerFunc)
	router.HandleFunc("/hello", helloHandlerFunc).Methods("GET")

	router.HandleFunc("/auth/login", authHandlerFunc).Methods("GET")
	router.HandleFunc("/auth/login_post", bou.Authenticate("local")).Methods("GET")

	userRoutes := mux.NewRouter()
	userRoutes.HandleFunc("/user/{userid}/profile", profileHandlerFunc)
	router.PathPrefix("/user").Handler(nimble.New().
		UseHandlerFunc(bounce.IsLoggedIn).
		Use(userRoutes),
	)

	// router goes last
	nim.Use(router)
	nim.Run(":3000")
}

func profileHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Visiting my logged in profile page!")
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from profile, value is " + value.(string))
	}
	// using mux
	fmt.Println("from profile, the userid is " + mux.Vars(r)["userid"])
}

func helloHandlerFunc(w http.ResponseWriter, r *http.Request) {

	val := r.Header.Get("Cookie")
	fmt.Println("... val = "+ val)

	fmt.Fprintf(w, "Hello world!")
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from helloHandlerFunc, value is " + value.(string))
	}
	if value, ok := context.GetOk(r, "valueA"); ok {
		fmt.Println("from helloHandlerFunc, valueA is " + value.(string))
	}
	if value, ok := context.GetOk(r, "valueB"); ok {
		fmt.Println("from helloHandlerFunc, valueB is " + value.(string))
	}
}

func middlewareA(w http.ResponseWriter, r *http.Request) {
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from middlewareA, value is " + value.(string))
	} else {
		fmt.Println("from middlewareA, value is nil")
	}
	context.Set(r, "value", "A")
	context.Set(r, "valueA", "A")
}

func middlewareB(w http.ResponseWriter, r *http.Request) {
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from middlewareB, value is " + value.(string))
	}
	context.Set(r, "value", "B")
	context.Set(r, "valueB", "B")
}

func authHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login page: key in user name and password !")
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from authHandlerFunc, value is " + value.(string))
	}
	// using mux
	fmt.Println("from authHandlerFunc, the userid is " + mux.Vars(r)["userid"])
}

func middlewareAuth(w http.ResponseWriter, r *http.Request) {
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from middlewareAuth, value is " + value.(string))
	}
	context.Set(r, "value", "AUTH")
}