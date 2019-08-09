package util

import (
	"math/rand"
	"time"
)

/*
t := time.Now()
formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
*/

type StatMonth struct {
	StartDate, EndDate time.Time
}

// Форматирование времени
func FormatDate(t time.Time, format string) string {
	return t.Format(format) // Аналогично: YYYY-MM-DD
}

// Генерация начала и конец дат всех месяцев в текущем году
func DateStatGenerate() []StatMonth {
	sms := []StatMonth{}
	currentLocation := time.Now()
	for month := 1; month < 13; month++ {

		firstOfMonth := time.Date(currentLocation.Year(), time.Month(month), 1, 0, 0, 0, 0, currentLocation.Location())
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		sm := StatMonth{firstOfMonth, lastOfMonth}
		sms = append(sms, sm)
		//fmt.Println(firstOfMonth.Format("2006-01-02"))
		//fmt.Println(lastOfMonth.Format("2006-01-02"))
	}
	return sms
}

// Генерация начала и конец дат всех месяцев в текущем году
func DateYearGenerate() StatMonth {
	sm := StatMonth{}
	currentLocation := time.Now()
	firstOfYear := time.Date(currentLocation.Year(), 1, 1, 0, 0, 0, 0, currentLocation.Location())
	lastOfYear := time.Date(currentLocation.Year(), 12, 1, 0, 0, 0, 0, currentLocation.Location())
	sm = StatMonth{firstOfYear, lastOfYear}
	return sm
}

// Случайное значение даты
func RanDate() time.Time {
	min := time.Date(2019, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 12, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// Случайное значение массива строк
func RandString(array []string) string {
	return array[rand.Intn(len(array))]
}
