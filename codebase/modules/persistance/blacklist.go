package persistance

import (
    "fmt"
    "log"
    "database/sql"
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

func (this *BlackListModule) Save(ip string) {
    if ! this.CheckInBlackList(ip) {
        log.Printf("Save address [%s] to black list\n", ip)
        if _, err := this.db.Exec(addToBlaclist(ip)); err != nil {
            log.Fatalln(err.Error())
        }
    }
}

func (this *BlackListModule) CheckInBlackList(ip string) (exists bool) {
    res := this.db.QueryRow(checkInBlackList(ip))
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