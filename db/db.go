package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/piotrprz/nufito/shared"
)

func NewService() shared.NufitoService {
	return &nufitoService{}
}

type nufitoService struct {
}

func (svc nufitoService) GetTrainers() ([]string, error) {
	db, err := sql.Open("postgres", "user=nufito dbname=nufito password=mysecretpassword host=db sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	// age := 21
	//	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	rows, err := db.Query("SELECT * FROM trainers")

	trainers := []string{}

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
		trainers = append(trainers, name)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return trainers, nil
}

func (svc *nufitoService) AddTrainer(trainer string) error {
	db, err := sql.Open("postgres", "user=nufito dbname=nufito password=mysecretpassword host=db sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO trainers(name) VALUES($1) RETURNING id as LastInsertId;")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(trainer)
	if err != nil {
		log.Fatal(err)
	}
	// Does not work in pg driver :(
	// lastId, err := res.LastInsertId()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	log.Printf("affected = %d\n", rowCnt)

	return nil
}
