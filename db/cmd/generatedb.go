package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"../../db"
	"../../model"
	"../../util"
	"golang.org/x/crypto/bcrypt"
)

//argsWithProg := os.Args

type Config struct {
	ListenSpec string

	Db db.Config
}

func processFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", "localhost:3000", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "host=localhost port=5433 user=postgres password=31yu*#km dbname=gowebapp sslmode=disable", "DB Connect String")

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

func adduser() {
	cfg := processFlags()
	m, err := Run(cfg)

	username := flag.String("username", "user", "a username...")
	password := flag.String("password", "12345", "a password...")
	email := flag.String("email", "user@uszn.avo.ru", "a email...")
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
		u.IsAdmin = false

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

func generateOrders() {
	cfg := processFlags()
	m, err := Run(cfg)

	docType := []string{0: "приказ", 1: "распоряжение"}
	kindOfDoc := []string{0: "По основной (профильной) деятельности", 1: "По личному составу", 2: "По административно-хозяйственным вопросам"}
	docLabel := []string{0: "ПД", 1: "ДСП", 2: "Свободный доступ (для общего пользования)"}

	for i := 1; i < 200; i++ {
		o := model.Order{}

		o.DocType = util.RandString(docType)     // Тип документа (приказ, распоряжение)
		o.KindOfDoc = util.RandString(kindOfDoc) // Вид документа (личный состав, основная деятельность)
		o.DocLabel = util.RandString(docLabel)   // Пометка секретности (персональные данные, ДСП)
		o.RegDate = util.RanDate()               // Дата регистрации
		o.RegNumber = "578"                      // Регистрационный номер
		o.Description = "О работе..."            // Описание
		o.Author = "Д.А. Бородулин"              // Автор
		o.FileOriginal = "ссылка"                // Оригинальный файл
		o.FileCopy = "ссылка"                    // Копия файла
		o.Current = true                         // Флаг действия документа

		err = m.CreateOrder(o)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	}

}

func main() {
	adduser()
	//generateOrders()
}
