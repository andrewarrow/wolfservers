package sqlite

import "fmt"

func List() {
	CreateSchema()
	InsertStake()
	fmt.Println("vim-go")
}
