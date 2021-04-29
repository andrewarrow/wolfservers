package sqlite

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
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
create table stakes (provider text, producer text, relay text, ssh_key text, ssh_key_pub text, created_at datetime, amount integer not null default 100);

CREATE VIEW view_stakes as select provider, ssh_key from stakes order by created_at desc;


`

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

	phrase := os.Getenv("WOLF_PHRASE")
	if len(phrase) < 36 {
		fmt.Println("wolves use longer phrases.")
		return
	}
	data := "ssh priv key"
	shhh := encrypt([]byte(data), phrase)
	encodedStr := base64.StdEncoding.EncodeToString(shhh)

	s := `insert into stakes (provider, producer, relay, ssh_key, ssh_key_pub, created_at) values (?, ?, ?, ?, ?, ?)`
	iia, _ := tx.Prepare(s)
	iia.Exec("1", "2", "3", encodedStr, "5", ts)

	tx.Commit()
}
