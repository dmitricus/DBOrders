package model

import (
	"time"
)

// Order is ...
type Order struct {
	ID           int64     // Идентификатор
	DocType      string    // Тип документа (приказ, распоряжение)
	KindOfDoc    string    // Вид документа (личный состав, основная деятельность)
	DocLabel     string    // Пометка секретности (персональные данные, ДСП)
	RegDate      time.Time // Дата регистрации
	RegNumber    string    // Регистрационный номер
	Description  string    // Описание
	Username     string    // Идентификатор пользователя
	FileOriginal string    // Оригинальный файл
	FileCopy     string    // Копия файла
	Current      bool      // Флаг действия документа
}
