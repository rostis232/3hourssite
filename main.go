package main

import (
	"fmt"
	"github.com/upper/db/v4/adapter/postgresql"
	"html/template"
	"log"
	"net/http"
)

//Налаштування підключення до БД
var settings = postgresql.ConnectionURL{
	Database: `3hourssite`,
	Host:     `localhost:5432`,
	User:     `admin`,
	Password: `admin`,
}

//Створення структури статті
type Article struct {
	Id       uint16 `db:"id"`
	Title    string `db:"title"`
	Anons    string `db:"anons"`
	FullText string `db:"full_text"`
}

//Создание массива
var articles []Article

//Функція обробки індексної сторінки
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//Підключаємося до БД
	sess, err := postgresql.Open(settings)
	if err != nil {
		log.Fatal("Open: ", err)
	}
	//Відкладенне відключення до БД
	defer sess.Close()

	//Выборка данных

	res, err := sess.SQL().Query("SELECT * FROM articles")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//Обнуляем список, иначе при каждом обновлении страницы количество постов будет увеличиваться!!
	articles = []Article{}
	//В цикле все объекты структуры добавляем в слайс
	for res.Next() {
		var article Article
		err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)
		if err != nil {
			fmt.Println(err)
		}
		articles = append(articles, article)
	}
	// Здесь мы передаем на фронт шаблон и готовый слайс с объектами структуры
	t.ExecuteTemplate(w, "index", articles)
}

//Функція обробки сторінки створення статті
func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

//Функція запису данних про нову статтю в БД
func save_article(w http.ResponseWriter, r *http.Request) {
	//Отримуємо данні з форм і записуємо їх у змінні
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	//Підключаємося до БД
	sess, err := postgresql.Open(settings)
	if err != nil {
		log.Fatal("Open: ", err)
	}
	//Відкладенне відключення до БД
	defer sess.Close()
	//Добавление записи
	insert, err := sess.SQL().Query(fmt.Sprintf("INSERT INTO articles (title, anons, full_text) VALUES ('%s', '%s', '%s')", title, anons, full_text))
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()
	//Редірект на головну сторінку
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))) //НЕ ПРАЦЮЄ!!!
	http.HandleFunc("/", index)
	http.HandleFunc("/create/", create)
	http.HandleFunc("/save_article/", save_article)
	http.ListenAndServe(":8081", nil)
}

func main() {
	handleFunc()
}
