package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"

	"../context"
	"../model"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Config is ...
type Config struct {
	Assets http.FileSystem
}

var store = sessions.NewCookieStore(
	[]byte("You probably want to change this"),
	[]byte("Seriously, I mean it. Change it."))

///// Authentication
func LoginHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := context.Get(r, "session").(*sessions.Session)

		loginTmpl := "assets/templates/login.html"

		params := struct {
			Flashes []interface{}
		}{}

		if r.Method == "GET" {
			params.Flashes = session.Flashes()
			s, err := loadTmpl(loginTmpl, params)
			if err != nil {
				log.Printf("error loading template: %s\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
			session.Save(r, w)
			fmt.Fprint(w, s)

		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
			}
			username := r.Form["username"][0]
			password := r.Form["password"][0]

			u, err := m.GetUserByUsername(username)
			if err != nil {
				log.Printf("error: %s\n", err)
				session.AddFlash("err: " + err.Error())
				err = session.Save(r, w)
				if err != nil {
					log.Printf("error saving session: %s\n", err)
				}
				http.Redirect(w, r, "/login", 301)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			if err != nil {
				session.AddFlash("err: " + err.Error())
				err = session.Save(r, w)
				if err != nil {
					log.Printf("error saving session: %s\n", err)
				}

				http.Redirect(w, r, "/login", 301)
				return
			}

			session.Values["id"] = u.ID
			err = session.Save(r, w)
			if err != nil {
				log.Printf("error saving session: %s\n", err)
			}
			http.Redirect(w, r, "/", 301)
		}
	}
}

func LogoutHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := context.Get(r, "session").(*sessions.Session)
		delete(session.Values, "id")
		session.Save(r, w)
		http.Redirect(w, r, "/login", 301)
	}
}

func indexHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := m.GetOrders()
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "index.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", orders); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func ListOrdersHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := m.GetOrders()
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "orders.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", orders); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func EditOrderHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := int64(intVar(vars, "id"))
		order := model.Order{}

		if id != 0 {
			order, err = m.GetOrder(id)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
				return
			}
		}

		if r.Method == "POST" {
			r.ParseForm()
			order.Username = r.Form["username"][0]
			order.Email = r.Form["email"][0]

			err = m.UpdateOrder(user)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
			}
			http.Redirect(w, r, "/users", 301)

		}

		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "edit_users.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", user); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
func ListUsersHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := m.GetUsers()
		if err != nil {
			log.Printf("err: %+v\n", err.Error())
			return
		}

		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "users.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", users); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
func intVar(vars map[string]string, k string) int64 {
	var vv int64
	if v, ok := vars[k]; ok {
		vv, _ = strconv.ParseInt(v, 0, 32)
	}
	return vv
}

func EditUserHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := int64(intVar(vars, "id"))
		user := model.User{}

		if id != 0 {
			user, err = m.GetUser(id)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
				return
			}
		}

		if r.Method == "POST" {
			r.ParseForm()
			user.Username = r.Form["username"][0]
			user.Email = r.Form["email"][0]

			err = m.UpdateUser(user)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
			}
			http.Redirect(w, r, "/users", 301)

		}

		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "edit_users.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", user); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func DeleteUserHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := int64(intVar(vars, "id"))

		if id != 0 {
			_, err := m.GetUser(id)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
				return
			}
		}
		err = m.DeleteUser(id)
		if err != nil {
			fmt.Fprintf(w, "err: %s", err)
			return
		}

		fmt.Fprintf(w, "Deleting user: %d", id)
	}
}

///// MIDDLEWARE
func Use(handler http.HandlerFunc, m *model.Model, mids ...func(http.Handler, *model.Model) http.HandlerFunc) http.HandlerFunc {
	for _, mid := range mids {
		handler = mid(handler, m)
	}
	return handler
}

func ContextManager(h http.Handler, m *model.Model) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := store.Get(r, "gowe")
		if err != nil {
			log.Printf("ContextManager: err: %s\n", err)
			return
		}

		r = context.Set(r, "session", session)

		if id, ok := session.Values["id"]; ok {
			u, err := m.GetUser(id.(int64))
			if err != nil {
				r = context.Set(r, "user", nil)
			} else {
				r = context.Set(r, "user", u)
			}
		} else {
			r = context.Set(r, "user", nil)
		}

		h.ServeHTTP(w, r)

		context.Clear(r)
	})
}

func RequireLogin(h http.Handler, m *model.Model) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u := context.Get(r, "user"); u != nil {
			h.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	})
}

///// HELPERS
func loadTmpl(path string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Error parsing template: %s", path)
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Println(err)
	}

	return buf.String(), err
}

func Logger(h http.Handler, m *model.Model) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s", r.URL)
		h.ServeHTTP(w, r)
	})
}

// Start is ...
func Start(cfg Config, m *model.Model, listener net.Listener) {
	router := mux.NewRouter()

	router.HandleFunc("/", Use(indexHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/login", LoginHandler(cfg, m))
	router.HandleFunc("/logout", Use(LogoutHandler(cfg), m, RequireLogin))
	router.HandleFunc("/users", Use(ListUsersHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/orders", Use(ListOrdersHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/edit/{id:[0-9]+}", Use(EditUserHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/edit/{id:[0-9]+}/delete", Use(DeleteUserHandler(cfg, m), m, RequireLogin))

	router.PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir("assets/css"))))
	router.PathPrefix("/img/").Handler(
		http.StripPrefix("/img/", http.FileServer(http.Dir("assets/img"))))
	router.PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js"))))
	router.PathPrefix("/templates/").Handler(
		http.StripPrefix("/templates/", http.FileServer(http.Dir("assets/templates"))))

	h := Use(router.ServeHTTP, m, Logger, ContextManager)

	//log.Printf("Listening on %s\n", listener)
	go http.Serve(listener, h)
}
