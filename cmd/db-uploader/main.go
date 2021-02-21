package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/user"
	"github.com/niklod/highload-social-network/user/city"
	"github.com/niklod/highload-social-network/user/interest"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.DB.ConnectionString())

	master, err := createDB(cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	slave, err := createDB(cfg.Slave1.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	userRepo := user.NewRepository(master, slave)
	cityRepo := city.NewRepository(master)
	interestRepo := interest.NewRepository(master)

	citySvc := city.NewService(cityRepo)
	interestSvc := interest.NewService(interestRepo)
	userSvc := user.NewService(userRepo, citySvc, interestSvc)

	ctx, cancel := context.WithCancel(context.Background())

	go startUploadUsers(ctx, userSvc)

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, syscall.SIGINT, syscall.SIGTERM)

	<-doneCh

	fmt.Println("Завершение работы")
	cancel()

}

func startUploadUsers(ctx context.Context, svc *user.Service) {
	var counter int
	userPostfix := 1

	for {

		select {
		case <-ctx.Done():
			fmt.Printf("Загружено %d строк", counter)
			return
		default:
			u := &user.User{
				FirstName: fmt.Sprintf("TestFirstName %d", userPostfix),
				Lastname:  fmt.Sprintf("TestLastName %d", userPostfix),
				Age:       12,
				Sex:       "Мужчина",
				City: city.City{
					Name: "Москва",
				},
				Login:    fmt.Sprintf("testuser%d", userPostfix),
				Password: "qwerty",
			}

			_, err := svc.Create(u)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("User %s created\n", u.FirstName)
		}

		userPostfix++
		time.Sleep(1 * time.Second)
	}
}

func createDB(conn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
