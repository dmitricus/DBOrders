package daemon

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"../db"
	"../model"
	"../ui"
)

type Config struct {
	ListenSpec string

	Db db.Config
	UI ui.Config
}

func Run(cfg *Config) error {
	log.Printf("Starting, HTTP on: %s\n", cfg.ListenSpec)
	// Инициализация соединение с БД
	db, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return err
	}
	// Создание модели БД
	m := model.New(db)

	l, err := net.Listen("tcp", cfg.ListenSpec)
	if err != nil {
		log.Printf("Error creating listener: %v\n", err)
		return err
	}
	// Запуск интерфейса пользователя
	ui.Start(cfg.UI, m, l)

	waitForSignal()

	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
