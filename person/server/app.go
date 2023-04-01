package server

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"person/common"
	"person/db"
	"person/person"
	"time"
	userHttp "person/person/delivery/http"
	userPostgres "person/person/repository/postgres"
	userUceCase "person/person/usecase"
)

type App struct {
	httpServer *http.Server
	userUC user.UseCase
}

func NewApp() *App {
	userRepo := userPostgres.NewUserRepository(db.GetDB())

	return &App{
		userUC: userUceCase.NewUserUseCase(userRepo),
	}
}

func (a *App) Run(port string) error {
	router := mux.NewRouter()

	userHttp.RegisterHTTPEndpoints(router, a.userUC)

	router.Use(common.CORSMiddlware)
	router.Use(mux.CORSMethodMiddleware(router))
	
	a.httpServer = &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
