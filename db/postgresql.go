package db

import (
	"database/sql"
	"fmt"
	"log"

	"../model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config is ...
type Config struct {
	ConnectString string
}

// InitDb is ...
func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("postgres", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}

		if err := p.createTablesIfNotExist(); err != nil {
			return nil, err
		}
		if err := p.prepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

type pgDb struct {
	dbConn         *sqlx.DB
	prefix         string
	sqlSelectUsers *sqlx.Stmt
	sqlInsertUser  *sqlx.NamedStmt
	sqlSelectUser  *sql.Stmt
}

func (p *pgDb) createTablesIfNotExist() error {
	create_sql := `
	-- users

       CREATE TABLE IF NOT EXISTS users (
       id SERIAL NOT NULL PRIMARY KEY,
       username TEXT NOT NULL,
	   password TEXT NOT NULL,
	   created  DATE NOT NULL,
	   email TEXT NOT NULL,
	   is_admin BOOLEAN NOT NULL DEFAULT false);

	-- orders

	   CREATE TABLE IF NOT EXISTS orders (
		id SERIAL NOT NULL PRIMARY KEY,
		doc_type TEXT NOT NULL,
		kind_of_doc TEXT NOT NULL,
		doc_label TEXT NOT NULL,
		reg_date DATE NOT NULL,
		reg_number TEXT NOT NULL,
		description TEXT NOT NULL,
		author TEXT NOT NULL,
		file_original TEXT NOT NULL,
		file_copy TEXT NOT NULL,
		current BOOLEAN NOT NULL DEFAULT false,
		old_order_id SERIAL NOT NULL);

    `
	if rows, err := p.dbConn.Query(create_sql); err != nil {
		return err
	} else {
		rows.Close()
	}
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {

	if p.sqlSelectUsers, err = p.dbConn.Preparex(
		"SELECT id, username, password, created, email, is_admin FROM users",
	); err != nil {
		return err
	}
	if p.sqlInsertUser, err = p.dbConn.PrepareNamed(
		"INSERT INTO users (username, password, created, email, is_admin) VALUES (:username, :password, :created, :email, :is_admin) " +
			"RETURNING id, username, password, created, email, is_admin",
	); err != nil {
		return err
	}
	if p.sqlSelectUser, err = p.dbConn.Prepare(
		"SELECT id, username, password, created, email, is_admin FROM users WHERE id = $1",
	); err != nil {
		return err
	}

	return nil
}

func (p *pgDb) GetUsers() ([]model.User, error) {
	users := []model.User{}
	rows, err := p.dbConn.Query("SELECT id, username, password, created, email, is_admin FROM users")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (p *pgDb) GetUser(userID int64) (model.User, error) {
	row := p.dbConn.QueryRow("SELECT id, username, password, created, email, is_admin FROM users WHERE id = $1", userID)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin)
	if err == sql.ErrNoRows {
		return user, err
	} else if err != nil {
		return user, err
	}
	return user, err
}

func (p *pgDb) CreateUser(user model.User) error {
	_, err := p.dbConn.Exec("INSERT INTO users (username, password, created, email, is_admin) VALUES ($1, $2, $3, $4, $5)",
		&user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin)

	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (p *pgDb) UpdateUser(user model.User) error {
	_, err := p.dbConn.Exec("UPDATE users set username = $1, password = $2, created = $3, email = $4, is_admin = $5 WHERE id = $6",
		&user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin, &user.ID)

	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (p *pgDb) DeleteUser(id int64) error {
	_, err := p.dbConn.Exec("DELETE from users where id = $1", id)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (p *pgDb) GetUserByUsername(username string) (model.User, error) {
	row := p.dbConn.QueryRow("SELECT id, username, password, created, email, is_admin FROM users WHERE username = $1", username)
	user := model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin)
	if err == sql.ErrNoRows {
		return user, err
	} else if err != nil {
		return user, err
	}
	return user, err
}

func (p *pgDb) GetOrder(id int64) (model.Order, error) {
	row := p.dbConn.QueryRow("SELECT id, username, password, created, email, is_admin FROM orders WHERE id = $1", id)

	order := model.Order{}
	err := row.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
		&order.Description, &order.Author, &order.FileOriginal, &order.FileCopy, &order.Current, &order.OldOrderID)
	if err == sql.ErrNoRows {
		return order, err
	} else if err != nil {
		return order, err
	}
	return order, err
}

func (p *pgDb) GetOrders() ([]model.Order, error) {
	rows, err := p.dbConn.Query("SELECT * from orders")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	orders := []model.Order{}

	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
			&order.Description, &order.Author, &order.FileOriginal, &order.FileCopy, &order.Current, &order.OldOrderID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, err
}

func (p *pgDb) DeleteOrder(id int64) error {
	_, err := p.dbConn.Exec("DELETE from orders where id = $1", id)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

// возвращаем пользователю страницу для редактирования объекта
func (p *pgDb) GetOrder(id int64) (model.Order, error) {

	row := p.dbConn.QueryRow("SELECT * from orders where id = $1", id)
	order := model.Order{}
	err := row.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
		&order.Description, &order.Author, &order.FileOriginal, &order.FileCopy, &order.Current, &order.OldOrderID)

	if err != nil {
		log.Println(err)
		return order, err
	}
	return order, err
}

// получаем измененные данные и сохраняем их в БД
func (p *pgDb) UpdateOrder(order model.Order) error {
	_, err := p.dbConn.Exec("UPDATE orders SET doc_type = $1, kind_of_doc = $2, doc_label = $3, reg_date = $4, reg_number = $5, description = $6, author = $7, file_original = $8, file_copy = $9, current = $10, old_order_id = $11 WHERE id = $12",
		&order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber, &order.Description, &order.Author, &order.FileOriginal, &order.FileCopy, &order.Current, &order.OldOrderID, &order.ID)

	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (p *pgDb) CreateOrder(order model.Order) error {
	_, err := p.dbConn.Exec("INSERT INTO orders (doc_type, kind_of_doc, doc_label, reg_date, reg_number, description, author, file_original, file_copy, current, old_order_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		&order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber, &order.Description, &order.Author, &order.FileOriginal, &order.FileCopy, &order.Current, &order.OldOrderID)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}
