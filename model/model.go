package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Model is ...
type Model struct {
	db
}

// New is ...
func New(db db) *Model {
	return &Model{
		db: db,
	}
}


func decodeJson(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(target)
}

func toJsonString(obj interface{}) string {
	js, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("{\"error\":%q}", err.Error())
	}
	return string(js)
}