package sqlite

import (
	"database/sql"
	"fmt"
)

func List() {
	CreateSchema()
	InsertStake()
	db := OpenTheDB()
	defer db.Close()
	ListRows(db)
}

func ListRows(db *sql.DB) {
	rows, err := db.Query(fmt.Sprintf("select provider,ssh_key from stakes"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s1 string
		var s2 string
		rows.Scan(&s1, &s2)
		fmt.Println(s1, s2)
	}
}
