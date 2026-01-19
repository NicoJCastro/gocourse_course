package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NicoJCastro/gocourse_course/internal/course"
	"github.com/NicoJCastro/gocourse_course/pkg/bootstrap"
	"github.com/NicoJCastro/gocourse_course/pkg/handler"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load(".env")
	//logger
	logger := bootstrap.InitLogger()

	//db
	db, err := bootstrap.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	pagLimitDef := os.Getenv("PAGINATION_LIMIT_DEFAUL")
	if pagLimitDef == "" {
		logger.Fatal("PAGINATION_LIMIT_DEFAUL is not set")
	}

	ctx := context.Background()

	courseRepo := course.NewRepo(db, logger)
	courseService := course.NewService(logger, courseRepo)
	courseEndpoints := course.MakeEndpoint(courseService, course.Config{LimPageDef: pagLimitDef})

	h := handler.NewCourseHTTPServer(ctx, courseEndpoints)

	port := os.Getenv("PORT")
	adress := "localhost:" + port

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         adress,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		logger.Println("listen in ", adress)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		logger.Println("error: ", err)
		os.Exit(1)
	}

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Cache-Control, X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
