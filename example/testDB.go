package main

import (
	"database/sql"
	"fmt"
	"os"
	// "time"

	_ "github.com/coroot/coroot-node-agent/example/mysql"
)

func main() {
	// i := 0
	// for {
	// 	i++
	// 	j := i % 2
	// 	switch j {
	// 	case 0:
	// 		testDB()
	// 	// case 1:
	// 	// 	testDB2()
	// 	}
	// fmt.Println("Test DB oke !")
	// 	time.Sleep(10 * time.Second)
	// }
	testDB()
	fmt.Println("Test DB oke !")
}

func testDB() {
	db, err := sql.Open("mysql", "annd2:password@tcp(localhost:3306)/mysql")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insert, err := db.Query("SELECT * FROM sbtest2")
	// insert, err := db.Query("SELECT User                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               FROM user")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func testDB2() {
	db, err := sql.Open("mysql", os.Getenv("URL"))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM sbtest1 WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	err = stmt.QueryRow("localhost").Err()
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
}

/*
export QUERY='SELECT * FROM sbtest1'
export URL='annd2:password@tcp(localhost:3306)/hihi'

*/
