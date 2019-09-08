package db

import (
	"log"

	"../model"
	_ "github.com/lib/pq"
)

func (p *pgDb) Select() ([]model.User, error) {
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
