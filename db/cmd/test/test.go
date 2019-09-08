package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

/*
type Stat_month struct {
	StartDate, EndDate time.Time
}

func date_stat_generate() []Stat_month {
	sms := []Stat_month{}
	currentLocation := time.Now()
	for month := 1; month < 13; month++ {

		firstOfMonth := time.Date(currentLocation.Year(), time.Month(month), 1, 0, 0, 0, 0, currentLocation.Location())
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		sm := Stat_month{firstOfMonth, lastOfMonth}
		sms = append(sms, sm)
		//fmt.Println(firstOfMonth.Format("2006-01-02"))
		//fmt.Println(lastOfMonth.Format("2006-01-02"))
	}
	return sms
}

// Случайное значение даты
func randate() time.Time {
	min := time.Date(2019, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 12, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// Случайное значение массива строк
func randstring(array []string) string {
	return array[rand.Intn(len(array))]
}
*/
var (
	limit     int //- Количество элементов на странице - 7
	all       int //- Общее количество элементов - 110
	linkLimit int //- Количество ссылок в состоянии - 5
	start     int //- Текущее смещение ( для первой страницы будет отсутствовать поэтому эту ситуацию нужно учесть)
)

func ArrayChunk(s []int, size int) [][]int {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]int
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// pagesList - массив с чанками
// needPage - Здесь это наш GET - параметр (START)
// Вернёт int - индекс нужного чанка:
func searchPage(pagesList [][]int, needPage int) int {
	for chunk, pages := range pagesList {
		if ok, _ := in_array(needPage, pages); ok {
			return chunk
		}
	}
	return 0
}

func pagination(start int) {
	limit = 7
	all = 110
	linkLimit = 5

	pages := math.Ceil(float64(all) / float64(limit)) // кол-во страниц
	var pagesArr []int

	// Заполняем массив: ключ - это номер страницы, значение - это смещение для БД.
	// Нумерация здесь нужна с единицы, а смещение с шагом = кол-ву материалов на странице.
	for i := 0; i < int(pages); i++ {
		pagesArr = append(pagesArr, i*int(limit))
	}
	// Теперь что бы на странице отображать нужное кол-во ссылок
	// дробим массив на чанки:
	allPages := ArrayChunk(pagesArr, linkLimit)

	fmt.Printf("%#v\n", allPages)
	// получаем от клиента текущую страницу и передаем в поиск текущего блока ссылок
	needChunk := searchPage(allPages, start)
	//fmt.Printf("needChunk %#v\n", needChunk)

	// Собственно выводим ссылки из нужного чанка
	for pageNum, ofset := range allPages[needChunk] {
		// Делаем текущую страницу не активной (т.е. не ссылкой):
		if ofset == start {
			fmt.Printf("%#v %v %s\n", ((pageNum + 1) + (linkLimit * needChunk)), ofset, "Текущая")
			continue
		}
		fmt.Printf("%#v %#v\n", ((pageNum + 1) + (linkLimit * needChunk)), ofset)
		//$htmlOut .= '<li><a href="'.link.'&'.varName.'='. ofset .'">'. pageNum . '</a></li>';
	}
}

// Случайное значение даты
func RanDate() time.Time {
	min := time.Date(2019, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 12, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// Форматирование времени
func FormatDate(t time.Time, format string) string {
	return t.Format(format) // Аналогично: YYYY-MM-DD
}

func ReplicatorParenthesis(number int) string {
	var buffer bytes.Buffer
	if number > 0 && number <= 5 {
		for n := 1; n <= number; n++ {
			buffer.WriteString("(")
		}
	}
	return buffer.String()
}

func main() {
	/*
		for i := 1; i < 50; i++ {
			fmt.Println(randate())
		}
	*/
	/*
		colors := []string{2: "blue", 0: "red", 1: "green"}
		for i := 1; i < 50; i++ {
			fmt.Println(randstring(colors))
		}
	*/
	// посылаем текущее смещение бд
	//pagination(0)
	//fmt.Println(RanDate())
	//fmt.Println(FormatDate(RanDate(), "03-01-2006"))

	fmt.Printf("%s", ReplicatorParenthesis(5))
}
