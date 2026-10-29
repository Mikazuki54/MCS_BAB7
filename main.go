package testbab7

import (
	"database/sql"
	"fmt"
	"log"
	"mcs_bab_7/database"
	"mcs_bab_7/routers"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbName   = ""
)

var (
	DB  *sql.DB
	err error
)

func main() {
	var PORT = ":50001"
	psqlInfo := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, password, dbName,
	)

	DB, err = sql.Open("posgres", psqlInfo)

	if err != nil {
		log.Fatal("gagal membuka DB", psqlInfo)
	}

	database.DBMigrations(DB)

	defer DB.Close()

	routers.StartServer().Run(PORT)
	fmt.Println("DB berhasil tersambung")
}
