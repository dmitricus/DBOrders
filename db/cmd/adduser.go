package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"../../db"
	"../../model"
	"golang.org/x/crypto/bcrypt"
)

/*
t := time.Now()
formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
*/

//argsWithProg := os.Args

type Config struct {
	ListenSpec string

	Db db.Config
}

func processFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", "localhost:3000", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "host=localhost port=5432 user=postgres password=31yu*#km dbname=gowebapp sslmode=disable", "DB Connect String")

	flag.Parse()
	return cfg
}

func Run(cfg *Config) (*model.Model, error) {
	log.Printf("Starting, HTTP on: %s\n", cfg.ListenSpec)
	// Инициализация соединение с БД
	db, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return nil, err
	}
	// Создание модели БД
	m := model.New(db)

	return m, err
}

func main() {
	cfg := processFlags()
	m, err := Run(cfg)

	username := flag.String("username", "admin", "a username...")
	password := flag.String("password", "12345", "a password...")
	email := flag.String("email", "admin@uszn.avo.ru", "a email...")
	//username := argsWithProg[:1]
	//password := argsWithProg[:2]
	//email := argsWithProg[:3]
	//is_admin := argsWithProg[:4]

	t := time.Now()

	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Println("Need username and password")
		return
	}

	u := model.User{}

	u, err = m.GetUserByUsername(*username)
	if err == nil {
		fmt.Println("Updating user and promoting to admin")
	} else {
		u.Username = *username
		u.Created = t
		u.Email = *email
		u.IsAdmin = true

		fmt.Println("Creating user and promoting to admin")
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}
		u.Password = string(hash)
		fmt.Printf("%+v\n", u)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
		err = m.CreateUser(u)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	u.Password = string(hash)
	fmt.Printf("%+v\n", u)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	err = m.UpdateUser(u)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

}
