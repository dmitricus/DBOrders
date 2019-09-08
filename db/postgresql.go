package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"../model"
	"../util"
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
	-- departaments

		CREATE TABLE IF NOT EXISTS departaments (
		 id SERIAL NOT NULL PRIMARY KEY,
		 title TEXT NOT NULL UNIQUE);
		 
	-- users

       CREATE TABLE IF NOT EXISTS users (
		id SERIAL NOT NULL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created  DATE NOT NULL,
		email TEXT NOT NULL,
		is_admin BOOLEAN NOT NULL DEFAULT false,
		departament_id SERIAL NOT NULL,
		FOREIGN KEY (departament_id) REFERENCES departaments (id) ON DELETE CASCADE);
	
	-- hbkind

		CREATE TABLE IF NOT EXISTS hbkind (
		 id SERIAL NOT NULL PRIMARY KEY,
		 name TEXT NOT NULL UNIQUE);
	
	-- hblabel

		 CREATE TABLE IF NOT EXISTS hblabel (
		  id SERIAL NOT NULL PRIMARY KEY,
		  name TEXT NOT NULL UNIQUE);

	-- hbtype

		  CREATE TABLE IF NOT EXISTS hbtype (
		   id SERIAL NOT NULL PRIMARY KEY,
		   name TEXT NOT NULL UNIQUE);

	-- orders

	   CREATE TABLE IF NOT EXISTS orders (
		id SERIAL NOT NULL PRIMARY KEY,
		doc_type_id SERIAL NOT NULL,
		kind_of_doc_id SERIAL NOT NULL,
		doc_label_id SERIAL NOT NULL,
		reg_date DATE NOT NULL,
		reg_number TEXT NOT NULL,
		description TEXT NOT NULL,
		user_id SERIAL NOT NULL,
		file_original TEXT NOT NULL,
		file_copy TEXT NOT NULL,
		current BOOLEAN NOT NULL DEFAULT false,
		FOREIGN KEY (doc_type_id) REFERENCES hbtype (id) ON DELETE CASCADE,
		FOREIGN KEY (kind_of_doc_id) REFERENCES hbkind (id) ON DELETE CASCADE,
		FOREIGN KEY (doc_label_id) REFERENCES hblabel (id) ON DELETE CASCADE);
    `
	if rows, err := p.dbConn.Query(create_sql); err != nil {
		log.Printf("error: %v", err)
		return err
	} else {
		rows.Close()
	}
	return nil
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Printf("error checkCount: %v", err)
		}
	}
	return count
}

func (p *pgDb) GetUsers() ([]model.User, error) {
	users := []model.User{}
	rows, err := p.dbConn.Query(`SELECT id, username, password, created, email, is_admin, 
	(SELECT title FROM departaments WHERE departaments.id = users.departament_id) AS title FROM users`)
	if err != nil {
		log.Printf("error GetUsers: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin, &user.Title)
		if err != nil {
			log.Printf("error GetUsers: %v", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (p *pgDb) GetUser(userID int64) (model.User, error) {
	row := p.dbConn.QueryRow(`SELECT id, username, created, email, is_admin, 
	(SELECT title FROM departaments WHERE departaments.id = users.departament_id) AS title FROM users WHERE id = $1`, userID)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Created, &user.Email, &user.IsAdmin, &user.Title)
	if err != nil {
		log.Printf("error GetUser: %v", err)
		return user, err
	}
	return user, err
}

func (p *pgDb) CreateUser(user model.User) error {
	_, err := p.dbConn.Exec(`INSERT INTO users (username, password, created, email, is_admin, departament_id) VALUES ($1, $2, $3, $4, $5, (SELECT id FROM departaments WHERE departaments.title = $6))`,
		&user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin, &user.Title)

	if err != nil {
		log.Printf("error CreateUser: %v", err)
		return err
	}
	return err
}

func (p *pgDb) UpdateUser(user model.User) error {
	_, err := p.dbConn.Exec(`UPDATE users set username = $1, password = $2, created = $3, email = $4, is_admin = $5, 
	departament_id = (SELECT id FROM departaments WHERE departaments.title = $6) WHERE id = $7`,
		&user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin, &user.Title, &user.ID)

	if err != nil {
		log.Printf("error UpdateUser: %v", err)
		return err
	}
	return err
}

func (p *pgDb) DeleteUser(id int64) error {
	_, err := p.dbConn.Exec("DELETE from users where id = $1", id)
	if err != nil {
		log.Printf("error DeleteUser: %v", err)
		return err
	}
	return err
}

func (p *pgDb) GetUserByUsername(username string) (model.User, error) {
	row := p.dbConn.QueryRow(`SELECT id, username, password, created, email, is_admin, 
	(SELECT title FROM departaments WHERE departaments.id = users.departament_id) AS title FROM users WHERE username = $1`, username)
	user := model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Created, &user.Email, &user.IsAdmin, &user.Title)
	if err != nil {
		log.Printf("error GetUserByUsername: %v", err)
		return user, err
	}
	return user, err
}

func (p *pgDb) GetOrders(limit, offset int) ([]model.Order, error) {
	rows, err := p.dbConn.Query(`SELECT id, 
	(SELECT name FROM hbtype WHERE hbtype.id = orders.doc_type_id) AS name, 
	(SELECT name FROM hbkind WHERE hbkind.id = orders.kind_of_doc_id) AS name, 
	(SELECT name FROM hblabel WHERE hblabel.id = orders.doc_label_id) AS name,
	reg_date, reg_number, description, 
	(SELECT username FROM users WHERE users.id = orders.user_id) AS username,
	file_original, file_copy, current FROM orders ORDER BY reg_date DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		log.Printf("error GetOrders: %v", err)
	}
	defer rows.Close()
	orders := []model.Order{}

	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
			&order.Description, &order.Username, &order.FileOriginal, &order.FileCopy, &order.Current)
		if err != nil {
			log.Printf("error GetOrders: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, err
}

// возвращаем пользователю страницу для редактирования объекта
func (p *pgDb) GetOrder(id int64) (model.Order, error) {

	row := p.dbConn.QueryRow(`SELECT id, 
	(SELECT name FROM hbtype WHERE hbtype.id = orders.doc_type_id) AS name, 
	(SELECT name FROM hbkind WHERE hbkind.id = orders.kind_of_doc_id) AS name, 
	(SELECT name FROM hblabel WHERE hblabel.id = orders.doc_label_id) AS name,
	reg_date, reg_number, description, 
	(SELECT username FROM users WHERE users.id = orders.user_id) AS username,
	file_original, file_copy, current FROM orders where id = $1`, id)
	order := model.Order{}
	err := row.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
		&order.Description, &order.Username, &order.FileOriginal, &order.FileCopy, &order.Current)

	if err != nil {
		log.Printf("error GetOrder: %v", err)
		return order, err
	}
	return order, err
}

func (p *pgDb) DeleteOrder(id int64) error {
	_, err := p.dbConn.Exec("DELETE from orders where id = $1", id)
	if err != nil {
		log.Printf("error DeleteOrder: %v", err)
		return err
	}
	return err
}

// получаем измененные данные и сохраняем их в БД
func (p *pgDb) UpdateOrder(order model.Order) error {
	_, err := p.dbConn.Exec(`UPDATE orders SET 
	(SELECT id FROM hbtype WHERE hbtype.name = $1), 
	(SELECT id FROM hbkind WHERE hbkind.name = $2), 
	(SELECT id FROM hblabel WHERE hblabel.name = $3)
	reg_date = $4, reg_number = $5, description = $6, user_id = (SELECT id FROM users WHERE users.username = $7), 
	file_original = $8, file_copy = $9, current = $10 WHERE id = $11`,
		&order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber, &order.Description,
		&order.Username, &order.FileOriginal, &order.FileCopy, &order.Current, &order.ID)

	if err != nil {
		log.Printf("error UpdateOrder: %v", err)
		return err
	}
	return err
}

func (p *pgDb) CreateOrder(order model.Order) error {
	_, err := p.dbConn.Exec(`INSERT INTO orders (doc_type_id, kind_of_doc_id, doc_label_id, reg_date, reg_number, description, user_id, 
			file_original, file_copy, current) VALUES (
			(SELECT id FROM hbtype WHERE hbtype.name = $1), 
			(SELECT id FROM hbkind WHERE hbkind.name = $2), 
			(SELECT id FROM hblabel WHERE hblabel.name = $3), 
			$4, $5, $6, (SELECT id FROM users WHERE users.username = $7), $8, $9, $10)`,
		&order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber, &order.Description, &order.Username, &order.FileOriginal, &order.FileCopy, &order.Current)
	if err != nil {
		log.Printf("error CreateOrder: %v", err)
		return err
	}
	return err
}

//WHERE DateField BETWEEN to_date('2010-01-01','YYYY-MM-DD') AND to_date('2010-01-02','YYYY-MM-DD')

// возвращаем количество приказов в промежутки дат
func (p *pgDb) GetDateOrders(startDate, endDate time.Time, limit, offset int) ([]model.Order, error) {
	rows, err := p.dbConn.Query(`SELECT id, 
	(SELECT name FROM hbtype WHERE hbtype.id = orders.doc_type_id) AS name, 
	(SELECT name FROM hbkind WHERE hbkind.id = orders.kind_of_doc_id) AS name, 
	(SELECT name FROM hblabel WHERE hblabel.id = orders.doc_label_id) AS name, 
	reg_date, reg_number, description, 
	(SELECT username FROM users WHERE orders.user_id = users.id) AS username,
	file_original, file_copy, current FROM orders WHERE reg_date BETWEEN $1 AND $2 ORDER BY reg_date DESC LIMIT $3 OFFSET $4`, util.FormatDate(startDate, "2006-01-02"), util.FormatDate(endDate, "2006-01-02"), limit, offset)

	orders := []model.Order{}
	if err != nil {
		log.Printf("error GetDateOrders: %v", err)
		return orders, err
	}
	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
			&order.Description, &order.Username, &order.FileOriginal, &order.FileCopy, &order.Current)
		if err != nil {
			log.Printf("error GetDateOrders: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, err
}

// возвращаем количество приказов в промежутки дат
func (p *pgDb) GetDateUserByUsername(startDate, endDate time.Time, username string, limit, offset int) ([]model.Order, error) {
	rows, err := p.dbConn.Query(`SELECT id, 
	(SELECT name FROM hbtype WHERE hbtype.id = orders.doc_type_id) AS name, 
	(SELECT name FROM hbkind WHERE hbkind.id = orders.kind_of_doc_id) AS name, 
	(SELECT name FROM hblabel WHERE hblabel.id = orders.doc_label_id) AS name, 
	reg_date, reg_number, description, 
	(SELECT username FROM users WHERE orders.user_id = users.id) AS username,
	file_original, file_copy, current FROM orders WHERE reg_date BETWEEN $1 AND $2 
	AND user_id = (SELECT id FROM users WHERE users.username = $3) ORDER BY reg_date DESC LIMIT $4 OFFSET $5`,
		util.FormatDate(startDate, "2006-01-02"), util.FormatDate(endDate, "2006-01-02"), username, limit, offset)

	orders := []model.Order{}
	if err != nil {
		log.Printf("error GetDateOrders: %v", err)
		return orders, err
	}
	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
			&order.Description, &order.Username, &order.FileOriginal, &order.FileCopy, &order.Current)
		if err != nil {
			log.Printf("error GetDateOrders: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, err
}

// возвращаем количество приказов в промежутки дат
func (p *pgDb) GetCountDateOrdersByUsername(startDate, endDate time.Time, username string) (int, error) {
	rows, err := p.dbConn.Query("SELECT COUNT(*) FROM orders WHERE reg_date BETWEEN $1 AND $2 AND user_id = (SELECT id FROM users WHERE users.username = $3)", util.FormatDate(startDate, "2006-01-02"), util.FormatDate(endDate, "2006-01-02"), username)
	if err != nil {
		log.Printf("error GetCountDateOrders: %v", err)
		return checkCount(rows), err
	}
	return checkCount(rows), err
}

// возвращаем количество приказов в промежутки дат
// Идея: Формирование запроса из кусков в зависимости от того что приходит в функицю
func (p *pgDb) GetSearchOrders(order model.Order, startDate, endDate time.Time) ([]model.Order, error) {
	var sqlQiery string

	sqlQierySelect := `SELECT id, (SELECT name FROM hbtype WHERE hbtype.id = orders.doc_type_id), (SELECT name FROM hbkind WHERE hbkind.id = orders.kind_of_doc_id), (SELECT name FROM hblabel WHERE hblabel.id = orders.doc_label_id), 
	reg_date, reg_number, description, (SELECT username FROM users WHERE orders.user_id = users.id), file_original, file_copy, current FROM orders`

	sqlQieryRegDate := fmt.Sprintf("%s %s", sqlQierySelect, "WHERE reg_date >= $1 AND reg_date <= $2")

	elements := make(map[string]string)
	elements[fmt.Sprintf("%s %s %s", " INTERSECT", sqlQierySelect, "WHERE doc_type_id = (SELECT id FROM hbtype WHERE hbtype.name =")] = order.DocType
	elements[fmt.Sprintf("%s %s %s", " INTERSECT", sqlQierySelect, "WHERE kind_of_doc_id = (SELECT id FROM hbkind WHERE hbkind.name =")] = order.KindOfDoc
	elements[fmt.Sprintf("%s %s %s", " INTERSECT", sqlQierySelect, "WHERE doc_label_id = (SELECT id FROM hblabel WHERE hblabel.name =")] = order.DocLabel
	elements[fmt.Sprintf("%s %s %s", " INTERSECT", sqlQierySelect, "WHERE (reg_number =")] = order.RegNumber
	elements[fmt.Sprintf("%s %s %s", " INTERSECT", sqlQierySelect, "WHERE (description =")] = order.Description
	elements[fmt.Sprintf("%s %s %s", " INTERSECT", sqlQierySelect, "WHERE user_id = (SELECT id FROM users WHERE users.username =")] = order.Username

	var orderParams []string
	var orderValues []interface{}

	orderValues = append(orderValues, util.FormatDate(startDate, "2006-01-02"))
	orderValues = append(orderValues, util.FormatDate(endDate, "2006-01-02"))
	i := 3
	for orderParam, value := range elements {
		if value != "" {
			orderParams = append(orderParams, fmt.Sprintf("%s $%v)", orderParam, i))
			orderValues = append(orderValues, value)
			i++
		}
	}
	orders := []model.Order{}
	sqlQiery = fmt.Sprintf("%s %s %s", sqlQieryRegDate, strings.Join(orderParams[:], " "), "ORDER BY reg_date DESC")
	rows, err := p.dbConn.Query(fmt.Sprintf("%s", sqlQiery), orderValues...)
	if err != nil {
		log.Printf("error GetSearchOrders: %v", err)
		return orders, err
	}
	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(&order.ID, &order.DocType, &order.KindOfDoc, &order.DocLabel, &order.RegDate, &order.RegNumber,
			&order.Description, &order.Username, &order.FileOriginal, &order.FileCopy, &order.Current)
		if err != nil {
			log.Printf("error GetSearchOrders: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, err
}

// возвращаем количество приказов в промежутки дат
func (p *pgDb) GetCountDateOrders(startDate, endDate time.Time) (int, error) {
	rows, err := p.dbConn.Query("SELECT COUNT(*) FROM orders WHERE reg_date BETWEEN $1 AND $2", util.FormatDate(startDate, "2006-01-02"), util.FormatDate(endDate, "2006-01-02"))
	if err != nil {
		log.Printf("error GetCountDateOrders: %v", err)
		return checkCount(rows), err
	}
	return checkCount(rows), err
}

func (p *pgDb) CreateDepartament(departament model.Departament) error {
	_, err := p.dbConn.Exec("INSERT INTO departaments (title) VALUES ($1)", &departament.Title)
	if err != nil {
		log.Printf("error CreateDepartament: %v", err)
		return err
	}
	return err
}

func (p *pgDb) GetDepartaments() ([]model.Departament, error) {
	departaments := []model.Departament{}
	rows, err := p.dbConn.Query(`SELECT title FROM departaments`)
	if err != nil {
		log.Printf("error GetDepartament: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		departament := model.Departament{}
		err := rows.Scan(&departament.Title)
		if err != nil {
			log.Printf("error GetDepartament: %v", err)
			continue
		}
		departaments = append(departaments, departament)
	}
	return departaments, nil
}

func (p *pgDb) GetDepartament(departamentID int64) (model.Departament, error) {
	row := p.dbConn.QueryRow(`SELECT title FROM departaments WHERE id = $1`, departamentID)

	departament := model.Departament{}
	err := row.Scan(&departament.Title)
	if err != nil {
		log.Printf("error GetDepartament: %v", err)
		return departament, err
	}
	return departament, err
}

func (p *pgDb) UpdateDepartament(departament model.Departament) error {
	_, err := p.dbConn.Exec(`UPDATE departaments set title = $1 WHERE id = $2`, &departament.Title, &departament.ID)
	if err != nil {
		log.Printf("error UpdateDepartament: %v", err)
		return err
	}
	return err
}

func (p *pgDb) DeleteDepartament(id int64) error {
	_, err := p.dbConn.Exec("DELETE from departaments where id = $1", id)
	if err != nil {
		log.Printf("error DeleteDepartament: %v", err)
		return err
	}
	return err
}

func (p *pgDb) CreateHBKindOfDoc(hbkind model.HBKindOfDoc) error {
	_, err := p.dbConn.Exec("INSERT INTO hbkind (name) VALUES ($1)", &hbkind.Name)
	if err != nil {
		log.Printf("error CreateBKindOfDoc: %v", err)
		return err
	}
	return err
}

func (p *pgDb) GetHBKindOfDoc() ([]model.HBKindOfDoc, error) {
	hbkinds := []model.HBKindOfDoc{}
	rows, err := p.dbConn.Query(`SELECT id, name FROM hbkind`)
	if err != nil {
		log.Printf("error GetBKindOfDoc: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		hbkind := model.HBKindOfDoc{}
		err := rows.Scan(&hbkind.ID, &hbkind.Name)
		if err != nil {
			log.Printf("error GetBKindOfDoc: %v", err)
			continue
		}
		hbkinds = append(hbkinds, hbkind)
	}
	return hbkinds, nil
}

func (p *pgDb) CreateHBDocLabel(hblabel model.HBDocLabel) error {
	_, err := p.dbConn.Exec("INSERT INTO hblabel (name) VALUES ($1)", &hblabel.Name)
	if err != nil {
		log.Printf("error CreateHBDocLabel: %v", err)
		return err
	}
	return err
}

func (p *pgDb) GetHBDocLabel() ([]model.HBDocLabel, error) {
	hblabels := []model.HBDocLabel{}
	rows, err := p.dbConn.Query(`SELECT id, name FROM hblabel`)
	if err != nil {
		log.Printf("error GetHBDocLabel: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		hblabel := model.HBDocLabel{}
		err := rows.Scan(&hblabel.ID, &hblabel.Name)
		if err != nil {
			log.Printf("error GetHBDocLabel: %v", err)
			continue
		}
		hblabels = append(hblabels, hblabel)
	}
	return hblabels, nil
}

func (p *pgDb) CreateHBDocType(hbtype model.HBDocType) error {
	_, err := p.dbConn.Exec("INSERT INTO hbtype (name) VALUES ($1)", &hbtype.Name)
	if err != nil {
		log.Printf("error CreateHBDocType: %v", err)
		return err
	}
	return err
}

func (p *pgDb) GetHBDocType() ([]model.HBDocType, error) {
	hbtypes := []model.HBDocType{}
	rows, err := p.dbConn.Query(`SELECT id, name FROM hbtype`)
	if err != nil {
		log.Printf("error GetHBDocType: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		hbtype := model.HBDocType{}
		err := rows.Scan(&hbtype.ID, &hbtype.Name)
		if err != nil {
			log.Printf("error GetHBDocType: %v", err)
			continue
		}
		hbtypes = append(hbtypes, hbtype)
	}
	return hbtypes, nil
}

func (p *pgDb) Get2HBDocType(codeFragment string) ([]model.HBDocType, error) {
	hbtypes := []model.HBDocType{}
	rows, err := p.dbConn.Query(`SELECT id, name FROM hbtype WHERE name LIKE '%' || $1 || '%'`, codeFragment)
	if err != nil {
		log.Printf("error GetHBDocType: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		hbtype := model.HBDocType{}
		err := rows.Scan(&hbtype.ID, &hbtype.Name)
		if err != nil {
			log.Printf("error Get2HBDocType: %v", err)
			continue
		}
		hbtypes = append(hbtypes, hbtype)
	}
	return hbtypes, nil
}
