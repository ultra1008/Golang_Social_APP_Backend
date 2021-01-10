package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"

	"github.com/niklod/highload-social-network/user"
	"github.com/niklod/highload-social-network/user/city"
	"github.com/niklod/highload-social-network/user/interest"

	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", cfg.DB.ConnectionString())

	db, err := sql.Open("mysql", cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepo := user.NewRepository(db)
	cityRepo := city.NewRepository(db)
	interestRepo := interest.NewRepository(db)

	// Services
	cityService := city.NewService(cityRepo)
	interestService := interest.NewService(interestRepo)
	userService := user.NewService(userRepo, cityService, interestService)

	ss := sessions.NewCookieStore([]byte(cfg.SecretKey))
	gob.Register(user.User{})

	// Handlers
	userHandler := user.NewHandler(userService, cityService, ss, interestService)

	srv := server.NewHTTPServer(cfg.Server)
	srv.BaseRouterGroup.Use(userHandler.AuthMiddleware)

	// Регистрациия
	srv.BaseRouterGroup.GET("/registrate", userHandler.HandleUserRegistrate)
	srv.BaseRouterGroup.POST("/registrate", userHandler.HandleUserRegistrateSubmit)

	// Вход Выход
	srv.BaseRouterGroup.GET("/login", userHandler.HandleUserLogin)
	srv.BaseRouterGroup.POST("/login", userHandler.HandleUserLoginSubmit)
	srv.BaseRouterGroup.GET("/logout", userHandler.HandleUserLogout)

	// User detail page
	srv.BaseRouterGroup.GET("/user/:login", userHandler.HandleUserDetail)

	// Добавление Удаление из друзей
	srv.BaseRouterGroup.POST("/user/:login/add_friend", userHandler.HandleAddFriend)
	srv.BaseRouterGroup.POST("/user/:login/delete_friend", userHandler.HandleDeleteFriend)

	// Список пользователей
	srv.BaseRouterGroup.GET("/users", userHandler.HandleUsersList)

	// Static
	srv.BaseRouterGroup.Static("/public/", "./static")

	srv.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Printf("received signal %s, stopping program...", sig)
	srv.Shutdown()
	signal.Stop(sigCh)
	log.Println("program stopped")
}
