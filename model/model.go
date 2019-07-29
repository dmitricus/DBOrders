package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	sessionName   = "default"
	sessionSecret = base64.StdEncoding.EncodeToString([]byte(sessionName))
	sessionStore  = sessions.NewCookieStore([]byte(sessionSecret), nil)
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

// UserJson is ...
type UserJson struct {
	Username string `json:"username", db:"username"`
	Password string `json:"password", db:"password"`
}

// Users is ...
func (m *Model) Users() ([]*User, error) {
	return m.SelectUsers()
}

// IsAuthenticated returns true if the user has a signed session cookie.
func (m *Model) IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	fmt.Println("authenticated: ", session.Values["authenticated"])
	// Проверьте, аутентифицирован ли пользователь
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
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

func (m *Model) Orders(w http.ResponseWriter, r *http.Request) ([]Order, error) {
	orders, err := m.SelectOrders()
	if err != nil {
		fmt.Sprintf("{\"error\":%q}", err.Error())
		return orders, err
	}
	return orders, err
}

func (m *Model) Login(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := &UserJson{}
	/*
		body, err := ioutil.ReadAll(r.Body)

		log.Printf("ответ: %q", body)
		err = json.Unmarshal(body, user)

		if err != nil {
			log.Printf("ошибка декодирования ответа: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("синтаксическая ошибка при смещении байта %d", e.Offset)
			}
			log.Printf("ответ: %q", body)
			return
		}
		//
		defer r.Body.Close()
	*/
	err = decodeJson(r, user)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Ошибка: ", err)
		return
	}
	// Аутентификация идет здесь
	//fmt.Println("Исходящие данные: ", user.Username, user.Password)
	aut := map[string]interface{}{
		"autLogin": user.Username,
		"autPass":  user.Password,
	}
	_, result := m.db.SelectUser(aut)
	//fmt.Println("Ответ от базы данных: ", result)
	if result {
		sessionStore.Options = &sessions.Options{
			MaxAge: 10 * 3600,
		}

		if flashes := session.Flashes(); len(flashes) > 0 {
			// Use the flash values.
		} else {
			// Set a new flash.
			session.AddFlash("Привет, успешный вход!")
		}
		// Установить пользователя в качестве аутентифицированного
		session.Values["authenticated"] = true
		session.Save(r, w)
		return
	}
}

func (m *Model) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отменить аутентификацию пользователея

	if flashes := session.Flashes(); len(flashes) > 0 {
		// Use the flash values.
	} else {
		// Set a new flash.
		session.AddFlash("Пока, сессия закрыта!")
	}
	session.Values["authenticated"] = false
	session.Options.MaxAge = 0
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}

func (m *Model) DeleteModelOrders(id string) error {
	err := m.DeleteOrder(id)
	if err != nil {
		fmt.Sprintf("{\"error\":%q}", err.Error())
		return err
	}
	return err
}

func (m *Model) EditModelOrder(id string) (Order, error) {
	o, err := m.EditOrders(id)
	if err != nil {
		fmt.Sprintf("{\"error\":%q}", err.Error())
		return o, err
	}
	return o, err
}

func (m *Model) UpdateModelOrder(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID string) error {
	err := m.UpdateOrders(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID)
	if err != nil {
		fmt.Sprintf("{\"error\":%q}", err.Error())
		return err
	}
	return err
}

func (m *Model) CreateModelOrder(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID string) error {
	err := m.CreateOrders(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID)
	if err != nil {
		fmt.Sprintf("{\"error\":%q}", err.Error())
		return err
	}
	return err
}