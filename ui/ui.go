package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path"

	//"time"
	"strconv"

	"../context"
	"../model"
	"../util"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	//"github.com/kabukky/httpscerts"
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
			session.Values["is_admin"] = u.IsAdmin
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
		output := []int{}
		sms := util.DateStatGenerate()
		for _, sm := range sms {
			count, err := m.GetCountDateOrders(sm.StartDate, sm.EndDate)
			if err != nil {
				log.Printf("{\"error\":%q}", err.Error())
				return
			}
			output = append(output, count)
			//log.Printf("Колличество записей: %v за период: %s %s", count, sm.StartDate.Format("2006-01-02"), sm.EndDate.Format("2006-01-02"))
			//log.Printf("Длинна массива: %v", len(output))
		}
		//log.Printf("Длинна массива: %v", len(output))

		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "index.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", output); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func ListOrdersHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Page struct {
			Orders                         []model.Order
			PaginationPages                []util.PaginationPage
			Next, Previous                 int
			NextIsActive, PreviousIsActive bool
		}
		// Получим первую и последнюю дату текущего года
		sm := util.DateYearGenerate()
		var (
			limit     = 7
			linkLimit = 5
			start     = 0
		)
		all, err := m.GetCountDateOrders(sm.StartDate, sm.EndDate)

		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if r.Method == "GET" {
			vars := mux.Vars(r)
			id := int(intVar(vars, "id"))
			if id != 0 {
				start = id
			}
		}
		paginationPages := util.Pagination(limit, all, linkLimit, start)
		log.Print("start: ", limit, all, linkLimit, start)

		orders, err := m.GetDateOrders(sm.StartDate, sm.EndDate, limit, start)
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}

		next := (start + limit)
		previous := (start - limit)
		previousIsActive := false
		nextIsActive := false
		if previous < 0 {
			previousIsActive = true
		}
		if next >= all {
			nextIsActive = true
		}

		page := Page{Orders: orders, PaginationPages: paginationPages, Next: next, Previous: previous, NextIsActive: nextIsActive, PreviousIsActive: previousIsActive}
		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "orders.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}

		if err := tmpl.ExecuteTemplate(w, "layout", page); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func ListArchiveOrdersHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Page struct {
			Orders                         []model.Order
			PaginationPages                []util.PaginationPage
			Next, Previous                 int
			NextIsActive, PreviousIsActive bool
		}
		// Получим первую и последнюю дату текущего года
		sm := util.DateYearGenerate()
		var (
			limit     = 7
			linkLimit = 5
			start     = 0
		)
		all, err := m.GetCountDateOrders(sm.StartDate, sm.EndDate)

		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if r.Method == "GET" {
			vars := mux.Vars(r)
			id := int(intVar(vars, "id"))
			if id != 0 {
				start = id
			}
		}
		paginationPages := util.Pagination(limit, all, linkLimit, start)
		log.Print("start: ", limit, all, linkLimit, start)

		orders, err := m.GetDateOrders(sm.StartDate, sm.EndDate, limit, start)
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}

		next := (start + limit)
		previous := (start - limit)
		previousIsActive := false
		nextIsActive := false
		if previous < 0 {
			previousIsActive = true
		}
		if next >= all {
			nextIsActive = true
		}

		page := Page{Orders: orders, PaginationPages: paginationPages, Next: next, Previous: previous, NextIsActive: nextIsActive, PreviousIsActive: previousIsActive}
		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "orders_archive.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}

		if err := tmpl.ExecuteTemplate(w, "layout", page); err != nil {
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
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
			}
			order.DocType = r.FormValue("DocType")
			order.KindOfDoc = r.FormValue("KindOfDoc")
			order.DocLabel = r.FormValue("DocLabel")
			//order.RegDate = r.FormValue("RegDate")
			order.RegNumber = r.FormValue("RegNumber")
			order.Description = r.FormValue("Description")
			order.Author = r.FormValue("Author")
			order.FileOriginal = r.FormValue("FileOriginal")
			order.FileCopy = r.FormValue("FileCopy")
			//order.Current = r.FormValue("Current")
			//order.OldOrderID = r.FormValue("OldOrderID")
			log.Println(order)
			log.Println(r.FormValue("Current"))
			/*
				err = m.UpdateOrder(order)
				if err != nil {
					fmt.Fprintf(w, "err: %s\n", err)
				}
			*/
			http.Redirect(w, r, "/orders", 301)
		}

		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "orders_edit.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", order); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func DeleteOrderHandler(config Config, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := int64(intVar(vars, "id"))

		if id != 0 {
			_, err := m.GetOrder(id)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
				return
			}
		}
		err = m.DeleteOrder(id)
		if err != nil {
			fmt.Fprintf(w, "err: %s", err)
			return
		}
		fmt.Fprintf(w, "Deleting order: %d", id)
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

		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "users_edit.html"))
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

func requireAdmin(h http.Handler, m *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "gowe")
		if err != nil {
			log.Printf("requireAdmin: err: %s\n", err)
			return
		}
		r = context.Set(r, "session", session)
		if isAdmin, ok := session.Values["is_admin"]; ok {
			if isAdmin == false {
				http.Error(w, "Admin required", http.StatusNotFound)
				return
			}
		} else {
			http.Redirect(w, r, "/login", 302)
		}
		h.ServeHTTP(w, r)
	}
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

	router.HandleFunc("/users", Use(ListUsersHandler(cfg, m), m, RequireLogin, requireAdmin))
	router.HandleFunc("/users/edit/{id:[0-9]+}", Use(EditUserHandler(cfg, m), m, RequireLogin, requireAdmin))
	router.HandleFunc("/users/edit/{id:[0-9]+}/delete", Use(DeleteUserHandler(cfg, m), m, RequireLogin, requireAdmin))

	//router.HandleFunc("/orders", Use(ListOrdersHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/orders/{id:[0-9]+}", Use(ListOrdersHandler(cfg, m), m, RequireLogin))
	//router.HandleFunc("/ordersarc", Use(ListArchiveOrdersHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/ordersarc/{id:[0-9]+}", Use(ListArchiveOrdersHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/orders/edit/{id:[0-9]+}", Use(EditOrderHandler(cfg, m), m, RequireLogin))
	router.HandleFunc("/orders/edit/{id:[0-9]+}/delete", Use(DeleteOrderHandler(cfg, m), m, RequireLogin, requireAdmin))

	router.PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir("assets/css"))))
	router.PathPrefix("/img/").Handler(
		http.StripPrefix("/img/", http.FileServer(http.Dir("assets/img"))))
	router.PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js"))))
	router.PathPrefix("/templates/").Handler(
		http.StripPrefix("/templates/", http.FileServer(http.Dir("assets/templates"))))

	h := Use(router.ServeHTTP, m, Logger, ContextManager)
	/*
		// Проверяем, доступен ли cert файл.
		err := httpscerts.Check("cert.pem", "key.pem")
		// Если он недоступен, то генерируем новый.
		if err != nil {
			err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:3000")
			if err != nil {
				log.Fatal("Ошибка: Не можем сгенерировать https сертификат.")
			}
		}
		go http.ListenAndServeTLS(listener, "cert.pem", "key.pem", h)
	*/
	go http.Serve(listener, h)
}
