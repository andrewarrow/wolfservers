package sqlite

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"
)

func GetWolfPhrase() string {
	phrase := os.Getenv("WOLF_PHRASE")
	if len(phrase) < 36 {
		fmt.Println("wolves use longer phrases.")
		return ""
	}
	return phrase
}

func InsertPat(provider, pat string) {
	ts := time.Now()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	phrase := GetWolfPhrase()
	if phrase == "" {
		return
	}

	shhh := encrypt([]byte(pat), phrase)
	encodedStr := base64.StdEncoding.EncodeToString(shhh)

	s := `insert into pats (provider, pat, created_at) values (?, ?, ?)`
	thing, _ := tx.Prepare(s)
	thing.Exec(provider, encodedStr, ts)

	tx.Commit()
	DisplayCopyDropboxNotice()
}

func InsertRow(name, provider, privKey, pubKey string) {
	ts := time.Now()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	phrase := GetWolfPhrase()
	if phrase == "" {
		return
	}
	shhh := encrypt([]byte(privKey), phrase)
	encodedStr := base64.StdEncoding.EncodeToString(shhh)

	s := `insert into stakes (name, provider, producer_ip, producer_id, relay_ip, relay_id, ssh_key, ssh_key_pub, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	thing, _ := tx.Prepare(s)
	thing.Exec(name, provider, "producer", "", "relay", "", encodedStr, pubKey, ts)

	tx.Commit()
	DisplayCopyDropboxNotice()
}
func DisplayCopyDropboxNotice() {
	fmt.Println("")
	fmt.Println("cp ~/wolf.db ~/Dropbox/db-work/")
	fmt.Println("scp -i ~/.ssh/wolf-91F4  ~/wolf.db root@cyborg.st:")
	fmt.Println("")
}
func InsertPaymentRow(name, pv, ps, sv, ss, sa, pa string) {
	ts := time.Now()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	phrase := GetWolfPhrase()
	if phrase == "" {
		return
	}
	shhh := encrypt([]byte(pv), phrase)
	encodedPV := base64.StdEncoding.EncodeToString(shhh)
	shhh = encrypt([]byte(ps), phrase)
	encodedPS := base64.StdEncoding.EncodeToString(shhh)

	shhh = encrypt([]byte(sv), phrase)
	encodedSV := base64.StdEncoding.EncodeToString(shhh)
	shhh = encrypt([]byte(ss), phrase)
	encodedSS := base64.StdEncoding.EncodeToString(shhh)

	shhh = encrypt([]byte(sa), phrase)
	encodedSA := base64.StdEncoding.EncodeToString(shhh)
	shhh = encrypt([]byte(pa), phrase)
	encodedPA := base64.StdEncoding.EncodeToString(shhh)

	sql := `insert into payment (name, pv, ps, sv, ss, sa, pa, created_at) values (?, ?, ?, ?, ?, ?, ?, ?)`
	thing, _ := tx.Prepare(sql)
	thing.Exec(name, encodedPV, encodedPS,
		encodedSV, encodedSS,
		encodedSA, encodedPA,
		ts)

	tx.Commit()
	DisplayCopyDropboxNotice()
}
func InsertNodeRow(name, v, s, c string) {
	ts := time.Now()
	db := OpenTheDB()
	defer db.Close()
	tx, _ := db.Begin()

	phrase := GetWolfPhrase()
	if phrase == "" {
		return
	}
	shhh := encrypt([]byte(v), phrase)
	encodedV := base64.StdEncoding.EncodeToString(shhh)
	shhh = encrypt([]byte(s), phrase)
	encodedS := base64.StdEncoding.EncodeToString(shhh)
	shhh = encrypt([]byte(c), phrase)
	encodedC := base64.StdEncoding.EncodeToString(shhh)

	sql := `insert into nodes (name, counter, vkey, skey, created_at) values (?, ?, ?, ?, ?)`
	thing, _ := tx.Prepare(sql)
	thing.Exec(name, encodedC, encodedV, encodedS, ts)

	tx.Commit()
	DisplayCopyDropboxNotice()
}
