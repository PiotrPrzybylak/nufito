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
	db, err := sql.Open("postgres", "user=nufito dbname=nufito password=mysecretpassword sslmode=disable")
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
	//	svc.Trainers = append(svc.Trainers, trainer)
	return nil
}
