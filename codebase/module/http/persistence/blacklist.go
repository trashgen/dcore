// +build ignore

package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type BlackListModule struct {
	db *sql.DB
}

func NewBlackListModule() *BlackListModule {
	dbConn, err := sql.Open("postgres", "user=postgres password=admin dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &BlackListModule{db: dbConn}
}

func (this *BlackListModule) Save(data string) {
	if !this.CheckExists(data) {
		if _, err := this.db.Exec(addToBlaclist(data)); err != nil {
			log.Fatalln(err.Error())
		}
	}
}

func (this *BlackListModule) CheckExists(id string) (exists bool) {
	res := this.db.QueryRow(checkInBlackList(id))
	if err := res.Scan(&exists); err != nil {
		log.Fatalln(err.Error())
	}
	return exists
}

func (this *BlackListModule) Close() {
	this.db.Close()
}

func addToBlaclist(ip string) string {
	return fmt.Sprintf("INSERT INTO blacklist (ip) VALUES ('%s')", ip)
}

func checkInBlackList(ip string) string {
	return fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM blacklist WHERE ip = '%s')", ip)
}
