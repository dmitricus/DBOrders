package model

type db interface {
	SelectUsers() ([]*User, error)
	SelectUser(aut map[string]interface{}) ([]*User, bool)
	SelectOrders() ([]Order, error)
	DeleteOrder(id string) error
	EditOrders(id string) (Order, error)
	UpdateOrders(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID string) error
	CreateOrders(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID string) error
}
