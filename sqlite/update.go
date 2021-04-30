package sqlite

import (
	"encoding/base64"
	"fmt"
	"os"
)

func UpdateIps(name, producer, relay string) {
	CreateSchema()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	s := `update stakes set producer_ip=?, relay_ip=? where name=?`
	thing, _ := tx.Prepare(s)
	thing.Exec(producer, relay, name)
	tx.Commit()
}
func UpdateIds(name, producer, relay string) {
	CreateSchema()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	s := `update stakes set producer_id=?, relay_id=? where name=?`
	thing, _ := tx.Prepare(s)
	thing.Exec(producer, relay, name)
	tx.Commit()
}
func UpdateRow(name, privKey, pubKey string) {
	CreateSchema()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	phrase := os.Getenv("WOLF_PHRASE")
	if len(phrase) < 36 {
		fmt.Println("wolves use longer phrases.")
		return
	}
	shhh := encrypt([]byte(privKey), phrase)
	encodedStr := base64.StdEncoding.EncodeToString(shhh)

	s := `update stakes set ssh_key=?, ssh_key_pub=? where name=?`
	thing, _ := tx.Prepare(s)
	thing.Exec(encodedStr, pubKey, name)

	tx.Commit()
}
