package main

import (
	"log"
	"net/http"
)

func main() {
	// Указываем, что статические файлы находятся в папке "static"
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Запускаем сервер на порту 8080
	log.Println("Сервер запущен на http://localhost:8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal(err)
	}
}
