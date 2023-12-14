package main

import (
	"log"
	"net/http"
	"time"

	"github.com/curso2-2/internal/course"
	"github.com/curso2-2/internal/enrollment"
	"github.com/curso2-2/internal/user"
	"github.com/curso2-2/pkg/bootstrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()

	if err != nil {
		l.Fatal(err)
	}

	userRepo := user.NewRepo(l, db)
	userServ := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userServ)

	courseRepo := course.NewRepo(l, db)
	courseSrv := course.NewService(l, courseRepo)
	courseEnd := course.MakeEndpoints(courseSrv)

	enrollRepo := enrollment.NewRepo(db, l)
	enrollSrv := enrollment.NewService(l, userServ, courseSrv, enrollRepo)
	enrollEnd := enrollment.MakeEndpoints(enrollSrv)

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	router.HandleFunc("/courses", courseEnd.Create).Methods("POST")
	router.HandleFunc("/courses/{id}", courseEnd.Get).Methods("GET")
	router.HandleFunc("/courses", courseEnd.GetAll).Methods("GET")
	router.HandleFunc("/courses/{id}", courseEnd.Update).Methods("PATCH")
	router.HandleFunc("/courses/{id}", courseEnd.Delete).Methods("DELETE")

	router.HandleFunc("/enrollment", enrollEnd.Create).Methods("POST")

	srv := &http.Server{
		Handler:      http.TimeoutHandler(router, time.Second*5, "Tiempo culminado"),
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err)
	}

}
