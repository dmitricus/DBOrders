package main

import (
	"os"

	//"time"

	"github.com/kpmy/odf/generators"
	"github.com/kpmy/odf/mappers"
	"github.com/kpmy/odf/mappers/attr"
	"github.com/kpmy/odf/model"
	_ "github.com/kpmy/odf/model/stub" //не забываем загрузить код скрытой реализации
	"github.com/kpmy/odf/xmlns"
)

func main() {
	if output, err := os.Create("demo2.odf"); err == nil {
		//закроем файл потом
		defer output.Close()
		//создадим пустую модель документа
		m := model.ModelFactory()
		//создадим форматтер
		fm := &mappers.Formatter{}
		//присоединим форматтер к документу
		fm.ConnectTo(m)
		//установим тип документа, в данном случае это текстовый документ
		fm.MimeType = xmlns.MimeText
		//инициализируем внутренние структуры
		fm.Init()
		//запишем простую строку
		//fm.WriteString("Привет, Дима!")
		{ //first page
			//fm.WriteString("\n\n\n\n\n\n\n\n\n\n")
			fm.RegisterFont("Times New Roman", "Times New Roman")
			fm.SetAttr(new(attr.TextAttributes).FontFace("Times New Roman").Size(14).Bold()).SetAttr(new(attr.ParagraphAttributes).AlignCenter())
			fm.WritePara("Администрация Владимирской области\nДепартамент социальной защиты населения\n")
			fm.SetAttr(new(attr.TextAttributes).FontFace("Times New Roman").Size(18).Bold())
			fm.WriteString("П Р И К А З")
			fm.WriteLn()
			fm.SetAttr(new(attr.TextAttributes).FontFace("Times New Roman").Size(14).Bold())
			fm.WritePara("___.___.______ г.						       №______\n")
			fm.WriteLn()
			fm.SetAttr(nil)
			fm.WritePara("\n")

			/*fm.SetAttr(new(attr.TextAttributes).Bold())
			fm.WriteString(time.Now().String())
			fm.SetAttr(new(attr.ParagraphAttributes).PageBreak()) // новая страница
			fm.WritePara("")
			fm.SetAttr(nil)*/
		}

		{
			fm.SetAttr(new(attr.TextAttributes).FontFace("Times New Roman").Italic().Size(14))
			fm.WriteString("Об утверждении положения о\nработе с внешними носителями\nинформации\n")
		}
		{ //fm.SetAttr(new(attr.TextAttributes).Bold().Size(18))
			//fm.WritePara(strings.ToUpper("Заголовок"))

			fm.WriteLn()

			paragraph := []string{
				"В целях защиты от утечек конфиденциальной информации через периферийные устройства департамента социальной защиты населения (далее - департамент) п р и к а з ы в а ю:",
				"1. Утвердить положение о работе со съемными носителями конфиденциальной информации департамента согласно приложению №1 к настоящему приказу.",
				"2. Государственным служащим департамента до 20.02.2017 г. сдать служебные съемные носители информации в информационно – компьютерный отдел.",
				"3. Информационно – компьютерному отделу департамента до 01.03.2017:",
				"- привести в соответствие рабочие станции пользователей согласно приложению №1 к настоящему приказу;",
				"- установить средства защиты от утечек конфиденциальной информации через периферийные устройства согласно приложению №2 к настоящему приказу.",
				"4. Утвердить форму журнала учёта съемных носителей конфиденциальной информации (персональных данных) согласно приложению №3 к настоящему приказу.",
				"5. Государственным служащим департамента указанным в приложении №2 к настоящему приказу до 01.03.2017 г. получить зарегистрированные носители информации и расписаться в журнале учёта съемных носителей конфиденциальной информации (персональных данных).",
				"6. Признать утратившим силу приказ департамента от __.__.2016 №___.",
			}
			for _, para := range paragraph {
				//fm.WriteLn()
				fm.SetAttr(new(attr.TextAttributes).FontFace("Times New Roman").Size(14)).SetAttr(new(attr.ParagraphAttributes).AlignJustify())
				fm.WritePara("\t" + para)
			}
		}
		{
			fm.WriteLn()
			fm.WritePara("\n\n")
			fm.SetAttr(new(attr.TextAttributes).FontFace("Times New Roman").Size(14)).SetAttr(new(attr.ParagraphAttributes).AlignJustify())
			fm.WriteString("Директор департамента				            Л.Е. Кукушкина\n")
		}
		//сохраним файл
		generators.GeneratePackage(m, nil, output, fm.MimeType)
	}
}
