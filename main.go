package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func schedule(db *sql.DB) {

	fmt.Println("bye.")
	var count int
	var id1 int
	start := time.Now()
	end := time.Now()
	interval := time.Now()
	update := time.Now()

	row := db.QueryRow("select * from logic;")
	err := row.Scan(&id1, &count, &start, &end, &interval)
	if err != nil {
		log.Fatal(err)
	}

	insertOutput := `insert into output values($1,$2,$3,$4,$5,$6,$7,$8)`
	inputQuery := `select * from inputTable where id=$1`
	update = start.Add(14*time.Minute + 59*time.Second) //(interval-1)*time.Minute

	fmt.Println(start)
	var auction int
	auction = 1
	t := 0
	for update != end {

		var id int
		var category string
		var prod_name string
		var description string
		var mrp string

		for i := 1; i <= count; i++ {

			row := db.QueryRow(inputQuery, t+i)
			err1 := row.Scan(&id, &category, &prod_name, &description, &mrp)
			if err1 != nil {
				log.Fatal(err1)
			}
			_, err2 := db.Exec(insertOutput, id, category, prod_name, description, mrp, start, update, auction)
			if err2 != nil {
				log.Fatal(err2)
			}

		}
		t = t + 5

		start = start.Add(15 * time.Minute) //interval*time.Minute --> reality
		update = update.Add(15 * time.Minute)
		auction++

	}

}
func main() {
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatal(err1)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	if _, err := db.Exec("TRUNCATE TABLE student;"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database student opened and ready.")
	schedule(db)
}
