package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/piotrprz/nufito/shared"
)

func NewService() shared.NufitoService {
	db, err := sql.Open("postgres", "user=nufito dbname=nufito password=mysecretpassword host=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return &nufitoService{db: db}
}

type nufitoService struct {
	db *sql.DB
}

func (svc nufitoService) GetTrainers() ([]shared.Trainer, error) {
	rows, err := svc.db.Query("SELECT * FROM trainers")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var (
		trainers = []shared.Trainer{}
		id       string
		name     string
	)

	for rows.Next() {
		err1 := rows.Scan(&id, &name)
		if err1 != nil {
			log.Fatal(err1)
		}
		trainer := shared.Trainer{Name: name, Id: id}
		trainers = append(trainers, trainer)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return trainers, nil
}

func (svc *nufitoService) AddTrainer(trainer string) error {

	stmt, err := svc.db.Prepare("INSERT INTO trainers(name) VALUES($1) RETURNING id as LastInsertId;")
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
