package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/andrewarrow/wolfservers/files"
	_ "github.com/mattn/go-sqlite3"
)

func OpenTheDB() *sql.DB {
	db, err := sql.Open("sqlite3", files.UserHomeDir()+"/wolf.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func CreateSchema() {
	db := OpenTheDB()
	defer db.Close()

	sqlStmt := `
create table stakes (provider text, producer text, relay text, ssh_key text,
              ssh_key_pub text, created_at datetime, amount integer not null default 100);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		//fmt.Printf("%q\n", err)
		return
	}
}

func InsertStake() {
	ts := time.Now()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	s := `insert into stakes (provider, producer, relay, ssh_key, ssh_key_pub, created_at) values (?, ?, ?, ?, ?, ?)`
	iia, _ := tx.Prepare(s)
	iia.Exec("1", "2", "3", "4", "5", ts)

	tx.Commit()
	db.Close()
}
