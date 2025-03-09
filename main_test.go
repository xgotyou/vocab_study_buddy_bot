package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("YDB_CONNECTION_STRING", "grpc://localhost:2136/local")
	InitDB(context.Background())

	err := clearDB()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func clearDB() error {
	db, err := sql.Open("ydb", os.Getenv("YDB_CONNECTION_STRING"))
	defer db.Close()

	if err != nil {
		return err
	}

	row := db.QueryRow("DELETE FROM dictionary_entries;")
	err = row.Err()

	if err != nil {
		return err
	}

	log.Println("Database successefully cleaned.")
	return nil
}
