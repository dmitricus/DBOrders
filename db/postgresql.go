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
	   password TEXT NOT NULL);

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
		"SELECT id, username, password FROM users",
	); err != nil {
		return err
	}
	if p.sqlInsertUser, err = p.dbConn.PrepareNamed(
		"INSERT INTO users (username, password) VALUES (:username, :password) " +
			"RETURNING id, username, password",
	); err != nil {
		return err
	}
	if p.sqlSelectUser, err = p.dbConn.Prepare(
		"SELECT id, username, password FROM users WHERE id = $1",
	); err != nil {
		return err
	}

	return nil
}

func (p *pgDb) SelectUsers() ([]*model.User, error) {
	user := make([]*model.User, 0)
	if err := p.sqlSelectUsers.Select(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (p *pgDb) SelectUser(aut map[string]interface{}) ([]*model.User, bool) {
	user := model.User{}
	row := p.dbConn.QueryRow("SELECT id, username, password FROM users WHERE username = $1", aut["autLogin"])
	fmt.Printf("SelectUser() login: %v\n", aut["autLogin"])
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	fmt.Printf("SelectUser() password: %v\n", user.Password)
	if err == sql.ErrNoRows {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	//if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {

	if user.Password != aut["autPass"] {
		return nil, false
	}
	return nil, true
}

func (p *pgDb) SelectOrders() ([]model.Order, error) {
	rows, err := p.dbConn.Query("select * from orders")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	orders := []model.Order{}

	for rows.Next() {
		o := model.Order{}
		err := rows.Scan(&o.ID, &o.DocType, &o.KindOfDoc, &o.DocLabel, &o.RegDate, &o.RegNumber,
			&o.Description, &o.Author, &o.FileOriginal, &o.FileCopy, &o.Current, &o.OldOrderID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		orders = append(orders, o)
	}
	return orders, err
}

func (p *pgDb) DeleteOrder(id string) error {
	_, err := p.dbConn.Exec("delete from orders where id = ?", id)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

// возвращаем пользователю страницу для редактирования объекта
func (p *pgDb) EditOrders(id string) (model.Order, error) {

	row := p.dbConn.QueryRow("select * from orders where id = ?", id)
	o := model.Order{}
	err := row.Scan(&o.ID, &o.DocType, &o.KindOfDoc, &o.DocLabel, &o.RegDate, &o.RegNumber,
		&o.Description, &o.Author, &o.FileOriginal, &o.FileCopy, &o.Current, &o.OldOrderID)

	if err != nil {
		log.Println(err)
		return o, err
	} else {
		return o, err
	}
}

// получаем измененные данные и сохраняем их в БД
func (p *pgDb) UpdateOrders(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID string) error {
	_, err := p.dbConn.Exec("update orders set id = ?, doc_type = ?, kind_of_doc = ?, doc_label = ?, reg_date = ?, reg_number = ?, description = ?, author = ?, file_original = ?, file_copy = ?, current = ?, old_order_id = ? where id = ?",
		ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID)

	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (p *pgDb) CreateOrders(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID string) error {
	_, err := p.dbConn.Exec("insert into orders (id, doc_type, kind_of_doc, doc_label, reg_date, reg_number, description, author, file_original, file_copy, current, old_order_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}
