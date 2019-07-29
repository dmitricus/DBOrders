package ui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path"
	"time"

	"../model"
	"github.com/gorilla/mux"
)

// Config is ...
type Config struct {
	Assets http.FileSystem
}

// Start is ...
func Start(cfg Config, m *model.Model, listener net.Listener) {

	server := &http.Server{
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16}
	
	router := mux.NewRouter()

	http.Handle("/", requireLogin(indexHandler(m), m))
	http.Handle("/users", requireLogin(usersHandler(m), m))
	http.Handle("/login", loginHandler(m))
	http.Handle("/logout", logoutHandler(m))
	http.Handle("/delete", requireLogin(DeleteHandler(m), m))
	router.Handle("/edit/{id:[0-9]+}", requireLogin(EditOrders(m), m)).Methods("GET")
	router.Handle("/edit/{id:[0-9]+}", requireLogin(EditOrderHandler(m), m)).Methods("POST")
	http.Handle("/create", requireLogin(CreateOrderHandler(m), m))
	http.Handle("/js/", http.FileServer(cfg.Assets))
	http.Handle("/img/", http.FileServer(cfg.Assets))
	http.Handle("/css/", http.FileServer(cfg.Assets))
	http.Handle("/templates/", http.FileServer(cfg.Assets))

	go server.Serve(listener)
}

const (
	cdnReact           = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"
	cdnReactDom        = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"
	cdnBabelStandalone = "https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.24.0/babel.min.js"
	cdnAxios           = "https://cdnjs.cloudflare.com/ajax/libs/axios/0.16.1/axios.min.js"
)

var loginFormTmpl = []byte(`
<html>
	<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<!-- Bootstrap CSS -->
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
		<link rel="stylesheet" href="/css/template.css" >
		<title>БД Приказов</title>
	</head>
	<body>
	<form name="authForm" class="form-signin" >
	<img class="mb-4" src="/img/logo.png" alt="" width="72" height="72">
	<h1 class="h3 mb-3 font-weight-normal">Пожалуйста, войдите в систему</h1>
	<label for="authLogin" class="sr-only">Логин</label>
	<input type="text" name="authLogin" class="form-control" placeholder="Логин" required autofocus>
	<label for="authPass" class="sr-only">Пароль</label>
	<input type="password" name="authPass" class="form-control" placeholder="Пароль" required>
	<div class="checkbox mb-3">
		<label>
		<input type="checkbox" value="remember-me"> Запомнить меня
		</label>
	</div>
	<button class="btn btn-lg btn-primary btn-block" id="submit" type="submit">Вход</button>
	</form>
	<script>
	document.getElementById("submit").addEventListener("click", function (e) {
		e.preventDefault();
	   // получаем данные формы
	   let authForm = document.forms["authForm"];
	   let username = authForm.elements["authLogin"].value;
	   let password = authForm.elements["authPass"].value;
	   // сериализуем данные в json
	   let user = JSON.stringify({username: username, password: password});
	   let request = new XMLHttpRequest();
	   // посылаем запрос на адрес "/login"
		request.open("POST", "/login", true);   
		request.setRequestHeader("Content-Type", "application/json");
		request.addEventListener("load", function () {
		   // получаем и парсим ответ сервера
			let receivedUser = JSON.parse(request.response);
			console.log(receivedUser.username, "-", receivedUser.password);   // смотрим ответ сервера
		});
		request.send(user);
		location.replace("/"); 
	});
	</script>
	<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
	</body>
</html>
`)

func indexHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//fmt.Fprintf(w, indexHTML)
		orders, err := m.Orders(w, r)
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		//log.Printf("Orders: %v\n", orders)
		/*wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}*/
		tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "index.html"))
		if err != nil {
			log.Printf("{\"error\":%q}", err.Error())
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout", orders); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	})
}

func usersHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := m.Users()
		if err != nil {
			http.Error(w, "This is an error", http.StatusBadRequest)
			return
		}
		// Запаковка в JSON
		js, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "This is an error", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(js))
	})
}

func requireLogin(h http.Handler, m *model.Model) http.Handler {
	loginURL := "/login"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.IsAuthenticated(w, r) {
			http.Redirect(w, r, loginURL, http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)

	})
}

func loginHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Write(loginFormTmpl)
			return
		}
		if r.Method != http.MethodGet {
			m.Login(w, r)
			return
		}
	})
}

func logoutHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			log.Printf("Метод: %v\n", http.MethodGet)
			m.Logout(w, r)
			return
		}
		if r.Method == http.MethodPost {
			log.Printf("Метод: %v\n", http.MethodPost)
			m.Logout(w, r)
			return
		}
	})
}

func DeleteHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		err := m.DeleteModelOrders(id)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		} else {
			http.Redirect(w, r, "/", 301)
		}
	})
}

// возвращаем пользователю страницу для редактирования объекта
func EditOrders(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		order, err := m.EditModelOrder(id)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		} else {
			tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "edit.html"))
			if err != nil {
				log.Printf("{\"error\":%q}", err.Error())
				return
			}
			if err := tmpl.ExecuteTemplate(w, "layout", order); err != nil {
				log.Println(err.Error())
				http.Error(w, http.StatusText(500), 500)
			}
		}
	})
}

// получаем измененные данные и сохраняем их в БД
func EditOrderHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		ID := r.FormValue("ID")
		DocType := r.FormValue("DocType")
		KindOfDoc := r.FormValue("KindOfDoc")
		DocLabel := r.FormValue("DocLabel")
		RegDate := r.FormValue("RegDate")
		RegNumber := r.FormValue("RegNumber")
		Description := r.FormValue("Description")
		Author := r.FormValue("Author")
		FileOriginal := r.FormValue("FileOriginal")
		FileCopy := r.FormValue("FileCopy")
		Current := r.FormValue("Current")
		OldOrderID := r.FormValue("OldOrderID")
		
		err = m.UpdateModelOrder(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	})
}

func CreateOrderHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Page struct {
			Title string
		}
		p := Page{Title: "Добавить приказ"}
		if r.Method == "POST" {

			err := r.ParseForm()
			if err != nil {
				log.Println(err)
			}
			ID := r.FormValue("ID")
			DocType := r.FormValue("DocType")
			KindOfDoc := r.FormValue("KindOfDoc")
			DocLabel := r.FormValue("DocLabel")
			RegDate := r.FormValue("RegDate")
			RegNumber := r.FormValue("RegNumber")
			Description := r.FormValue("Description")
			Author := r.FormValue("Author")
			FileOriginal := r.FormValue("FileOriginal")
			FileCopy := r.FormValue("FileCopy")
			Current := r.FormValue("Current")
			OldOrderID := r.FormValue("OldOrderID")

			err = m.CreateModelOrder(ID, DocType, KindOfDoc, DocLabel, RegDate, RegNumber, Description, Author, FileOriginal, FileCopy, Current, OldOrderID)

			//_, err = p.dbConn.Exec("insert into productdb.Products (model, company, price) values (?, ?, ?)", model, company, price)

			if err != nil {
				log.Println(err)
			}
			http.Redirect(w, r, "/", 301)
		} else {
			tmpl, err := template.ParseFiles(path.Join("assets/templates", "layout.html"), path.Join("assets/templates", "create.html"))
			if err != nil {
				log.Printf("{\"error\":%q}", err.Error())
				return
			}
			if err := tmpl.ExecuteTemplate(w, "layout", p); err != nil {
				log.Println(err.Error())
				http.Error(w, http.StatusText(500), 500)
			}
		}
	})
}
