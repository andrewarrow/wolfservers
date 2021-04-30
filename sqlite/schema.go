package sqlite

import (
	"database/sql"
	"fmt"

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
create table stakes (name text, provider text, producer_ip text, producer_id text, relay_ip text, relay_id text, ssh_key text, ssh_key_pub text, created_at datetime, amount integer not null default 100);
create table nodes (name text, counter text, vkey text, skey text, created_at datetime);

CREATE VIEW view_stakes as select provider_ip, ssh_key from stakes order by created_at desc;


`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		//fmt.Printf("%q\n", err)
		return
	}
}
