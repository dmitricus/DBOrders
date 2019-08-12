package util

import (
	"math"
	"reflect"
)

type PaginationPage struct {
	PageNum, Ofset int
	Active         string
}

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

func Pagination(limit, all, linkLimit, start int) []PaginationPage {

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

	//fmt.Printf("%#v\n", allPages)
	// получаем от клиента текущую страницу и передаем в поиск текущего блока ссылок
	needChunk := searchPage(allPages, start)
	//fmt.Printf("needChunk %#v\n", needChunk)

	paginationPages := []PaginationPage{}
	// Собственно выводим ссылки из нужного чанка
	if len(allPages) > 0 {
		for pageNum, ofset := range allPages[needChunk] {
			paginationPage := PaginationPage{}

			// Делаем текущую страницу не активной (т.е. не ссылкой):
			if ofset == start {
				paginationPage.PageNum = ((pageNum + 1) + (linkLimit * needChunk))
				paginationPage.Ofset = ofset
				paginationPage.Active = "disabled"
				paginationPages = append(paginationPages, paginationPage)
				continue
			}
			paginationPage.PageNum = ((pageNum + 1) + (linkLimit * needChunk))
			paginationPage.Ofset = ofset
			paginationPage.Active = ""
			paginationPages = append(paginationPages, paginationPage)

		}
		return paginationPages
	}
	return paginationPages
}
