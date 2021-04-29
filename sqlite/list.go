package sqlite

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
)

func List() {
	db := OpenTheDB()
	defer db.Close()
	ListRows(db)
}

func ListRows(db *sql.DB) {
	rows, err := db.Query(fmt.Sprintf("select name,provider,ssh_key from stakes"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	phrase := os.Getenv("WOLF_PHRASE")

	for rows.Next() {
		var s1 string
		var s2 string
		var s3 string
		rows.Scan(&s1, &s2, &s3)
		decodedBytes, _ := base64.StdEncoding.DecodeString(s3)
		shhh := decrypt(decodedBytes, phrase)
		fmt.Println(s1, s2, len(shhh))
	}
}
