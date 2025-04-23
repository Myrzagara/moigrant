package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {

	db, err := sql.Open("postgres", "user=gen_user password=Passcode@123 dbname=postgre host=93.183.82.176 port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 🔥 API endpoint — отдельно и ДО других обработчиков
	http.HandleFunc("/cities", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM cities")
		if err != nil {
			http.Error(w, "Ошибка запроса", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var cities []City
		for rows.Next() {
			var city City
			if err := rows.Scan(&city.ID, &city.Name); err != nil {
				http.Error(w, "Ошибка данных", http.StatusInternalServerError)
				return
			}
			cities = append(cities, city)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cities)
	})

	// Если у тебя есть HTML/JS/CSS-файлы — подгружаем их отсюда
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs) // ← это обработчик для HTML

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
