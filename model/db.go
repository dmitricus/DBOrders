package model

type db interface {
	GetUsers() ([]User, error)
	GetUser(id int64) (User, error)
	GetUserByUsername(username string) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	DeleteUser(id int64) error
	GetOrders() ([]Order, error)
	GetOrder() (Order, error)
	DeleteOrder(id int64) error
	UpdateOrder(order Order) error
	CreateOrder(order Order) error
}
