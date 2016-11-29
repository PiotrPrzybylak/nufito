package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=nufito dbname=nufito password=mysecretpassword sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// age := 21
	//	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	rows, err := db.Query("SELECT * FROM trainers")

	fmt.Println(rows)
	fmt.Println(err)

	var (
		id   int
		name string
	)

	defer rows.Close()
	for rows.Next() {
		err1 := rows.Scan(&id, &name)
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Println(id, name)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}
