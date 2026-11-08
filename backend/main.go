package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 30; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		log.Println("waiting for db:", err)
		time.Sleep(2 * time.Second)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS products(
		id SERIAL PRIMARY KEY, name TEXT NOT NULL, price NUMERIC NOT NULL)`); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			rows, err := db.Query("SELECT id,name,price FROM products ORDER BY id")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer rows.Close()
			out := []Product{}
			for rows.Next() {
				var p Product
				rows.Scan(&p.ID, &p.Name, &p.Price)
				out = append(out, p)
			}
			json.NewEncoder(w).Encode(out)
		case http.MethodPost:
			var p Product
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.Name == "" {
				http.Error(w, "bad request", 400)
				return
			}
			err := db.QueryRow("INSERT INTO products(name,price) VALUES($1,$2) RETURNING id",
				p.Name, p.Price).Scan(&p.ID)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(p)
		default:
			http.Error(w, "method not allowed", 405)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
