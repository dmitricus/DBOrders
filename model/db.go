package model

import (
	"time"
)

type db interface {
	GetUsers() ([]User, error)
	GetUser(id int64) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	DeleteUser(id int64) error
	GetOrders(limit, offset int) ([]Order, error)
	GetUserByUsername(username string) (User, error)
	GetOrder(id int64) (Order, error)
	DeleteOrder(id int64) error
	UpdateOrder(order Order) error
	CreateOrder(order Order) error
	GetDateOrders(startDate, endDate time.Time, limit, offset int) ([]Order, error)
	GetDateUserByUsername(startDate, endDate time.Time, username string, limit, offset int) ([]Order, error)
	GetCountDateOrdersByUsername(startDate, endDate time.Time, username string) (int, error)
	GetCountDateOrders(startDate, endDate time.Time) (int, error)
	GetSearchOrders(order Order, startDate, endDate time.Time) ([]Order, error)
	GetDepartaments() ([]Departament, error)
	GetDepartament(departamentID int64) (Departament, error)
	CreateDepartament(departament Departament) error
	CreateHBKindOfDoc(hbkind HBKindOfDoc) error
	GetHBKindOfDoc() ([]HBKindOfDoc, error)
	CreateHBDocLabel(hblabel HBDocLabel) error
	GetHBDocLabel() ([]HBDocLabel, error)
	CreateHBDocType(hbtype HBDocType) error
	GetHBDocType() ([]HBDocType, error)
}
