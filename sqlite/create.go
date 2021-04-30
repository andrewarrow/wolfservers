package sqlite

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"
)

func InsertRow(name, provider, privKey, pubKey string) {
	ts := time.Now()
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

	s := `insert into stakes (name, provider, producer_ip, producer_id, relay_ip, relay_id, ssh_key, ssh_key_pub, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	thing, _ := tx.Prepare(s)
	thing.Exec(name, provider, "producer", "", "relay", "", encodedStr, pubKey, ts)

	tx.Commit()
}
