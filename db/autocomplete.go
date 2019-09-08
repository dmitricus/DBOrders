package db

import (
	"log"

	"../model"
	_ "github.com/lib/pq"
)

func (p *pgDb) Get2HBDocType(codeFragment string) ([]model.HBDocType, error) {
	hbtypes := []model.HBDocType{}
	rows, err := p.dbConn.Query(`SELECT id, name FROM hbtype WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'`, codeFragment)
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
