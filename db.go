package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//type conn *sql.DB
var conn *sql.DB

func DbConnect() (*sql.DB, error) {
	var err error

	if conn == nil {
		conn, err = sql.Open("mysql", DSN)
		if err != nil {
			fmt.Println(err)
		}
	}

	return conn, err
}

func DbClose() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

func DbExec(sql string, args ...interface{}) (sql.Result, error) {
	db, err := DbConnect()
	if err != nil {
		return nil, err
	}

	return db.Exec(sql, args...)
}

func DbSelect(sql string, args ...interface{}) (*sql.Rows, error) {
	db, err := DbConnect()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil

	/*for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
	rows.Close()*/
}

func DbSelectOne(sql string, args ...interface{}) (*sql.Row, error) {
	db, err := DbConnect()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(sql, args...)

	return row, nil
}
