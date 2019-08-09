package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"../../db"
	"../../model"
	"../../util"
	"golang.org/x/crypto/bcrypt"
)

//argsWithProg := os.Args

type Config struct {
	ListenSpec string

	Db db.Config
}

func processFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", "localhost:3000", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "host=localhost port=5432 user=postgres password=31yu*#km dbname=gowebapp sslmode=disable", "DB Connect String")

	flag.Parse()
	return cfg
}

func Run(cfg *Config) (*model.Model, error) {
	log.Printf("Starting, HTTP on: %s\n", cfg.ListenSpec)
	// Инициализация соединение с БД
	db, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return nil, err
	}
	// Создание модели БД
	m := model.New(db)

	return m, err
}

func adduser(m *model.Model, username, password, email, title string, is_admin bool) {
	t := time.Now()
	flag.Parse()

	if username == "" || password == "" {
		fmt.Println("Need username and password")
		return
	}

	u := model.User{}
	u, err := m.GetUserByUsername(username)
	if err == nil {
		u.IsAdmin = is_admin
		u.Email = email
		u.Title = title
		fmt.Println("Updating user and promoting to: %s", username)
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}
		u.Password = string(hash)
		fmt.Printf("%+v\n", u)
		err = m.UpdateUser(u)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	} else {
		fmt.Println("Creating user and promoting to: %s", username)
		u.Username = username
		u.Created = t
		u.Email = email
		u.IsAdmin = is_admin
		u.Title = title
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}
		u.Password = string(hash)
		fmt.Printf("%+v\n", u)
		err = m.CreateUser(u)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	}
}

func generateOrders(m *model.Model, username string) {
	docType := []string{0: "приказ", 1: "распоряжение", 2: "постановление"}
	kindOfDoc := []string{0: "По основной (профильной) деятельности", 1: "По личному составу", 2: "По административно-хозяйственным вопросам"}
	docLabel := []string{0: "ПД", 1: "ДСП", 2: "Свободный доступ (для общего пользования)"}

	for i := 1; i < 200; i++ {
		o := model.Order{}

		o.DocType = util.RandString(docType)            // Тип документа (приказ, распоряжение)
		o.KindOfDoc = util.RandString(kindOfDoc)        // Вид документа (личный состав, основная деятельность)
		o.DocLabel = util.RandString(docLabel)          // Пометка секретности (персональные данные, ДСП)
		o.RegDate = util.RanDate()                      // Дата регистрации
		o.RegNumber = fmt.Sprintf("%d", rand.Intn(900)) // Регистрационный номер
		o.Description = "О работе в ГИС ОГ"             // Описание
		o.FileOriginal = "ссылка"                       // Оригинальный файл
		o.FileCopy = "ссылка"                           // Копия файла
		o.Current = true                                // Флаг действия документа
		o.Username = username                           // Автор

		err := m.CreateOrder(o)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		} else {
			fmt.Printf("generateOrders: %s\n", o)
		}
	}

}

func generateDepartaments(m *model.Model) {
	fmt.Printf("Start generateDepartaments")
	titles := []string{
		0:  "Информационно-компьютерный отдел",
		1:  "Отдел бухгалтерского учёта и отчётности",
		2:  "Отдел государственных закупок для государственных нужд",
		3:  "Отдел кадров и делопроизводства",
		4:  "Отдел компенсационных выплат и социальных гарантий",
		5:  "Отдел контроля и надзора в сфере социального обслуживания",
		6:  "Отдел организации назначения детских пособий и социальных выплат",
		7:  "Отдел организации социального обслуживания населения в стационарных учреждениях",
		8:  "Отдел по делам пожилых людей и инвалидов",
		9:  "Экономико-финансовый отдел",
		10: "Сектор правового обеспечения",
		11: "Сектор социального обслуживания семьи и детей, находящихся в трудной жизненной ситуации"}

	for _, title := range titles {
		d := model.Departament{}
		d.Title = title
		fmt.Printf("title: %s\n", title)
		err := m.CreateDepartament(d)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	}
}

func generateHBDocType(m *model.Model) {
	fmt.Printf("Start generateHBDocType")
	docType := []string{0: "приказ", 1: "распоряжение", 2: "постановление"}

	for _, name := range docType {
		hb := model.HBDocType{}
		hb.Name = name
		fmt.Printf("title: %s\n", name)
		err := m.CreateHBDocType(hb)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	}
}

func generateHBKindOfDoc(m *model.Model) {
	fmt.Printf("Start generateHBKindOfDoc")
	kindOfDoc := []string{0: "По основной (профильной) деятельности", 1: "По личному составу", 2: "По административно-хозяйственным вопросам"}

	for _, name := range kindOfDoc {
		hb := model.HBKindOfDoc{}
		hb.Name = name
		fmt.Printf("title: %s\n", name)
		err := m.CreateHBKindOfDoc(hb)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	}
}

func generateHBDocLabel(m *model.Model) {
	fmt.Printf("Start generateHBDocLabel")
	docLabel := []string{0: "ПД", 1: "ДСП", 2: "Свободный доступ (для общего пользования)"}

	for _, name := range docLabel {
		hb := model.HBDocLabel{}
		hb.Name = name
		fmt.Printf("title: %s\n", name)
		err := m.CreateHBDocLabel(hb)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	}
}


func main() {
	cfg := processFlags()
	m, err := Run(cfg)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	//generateDepartaments(m)
	//generateHBDocType(m)
	//generateHBKindOfDoc(m)
	//generateHBDocLabel(m)
	//adduser(m, "admin", "12345", "admin@uszn.avo.ru", "Информационно-компьютерный отдел", true)
	//adduser(m, "dmitrieva_av", "12345", "dmitrieva@uszn.avo.ru", "Отдел организации назначения детских пособий и социальных выплат", false)
	generateOrders(m, "dmitrieva_av")
}
