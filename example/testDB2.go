package main

import (
	"database/sql"
	"fmt"
	// "os"
	"time"

	_ "github.com/coroot/coroot-node-agent/example/mysql"
)

func main() {
	i := 0
	for {
		i++
		j := i % 2
		switch j {
		case 0:
			testDB3()
		case 1:
			testDB4()
		}
		fmt.Println("Test DB oke !")
		time.Sleep(10 * time.Second)
	}
	// testDB()
	// fmt.Println("Test DB oke !")
}

func testDB3() {
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

func testDB4() {
	db, err := sql.Open("mysql", "annd2:password@tcp(localhost:3306)/hihi")
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
