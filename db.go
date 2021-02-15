package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

var schema string = 
`BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "providers" (
	"id"	INTEGER,
	"name"	TEXT,
	"url"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "entries" (
	"id"	INTEGER,
	"provider_id"	INTEGER,
	"url"	TEXT,
	"title"	TEXT,
	"published_at"	INTEGER,
	"author"	TEXT,
	"fetched_at"	INTEGER,
	"json"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
COMMIT;
`

func DBInit(){
	var err error
	DB, err = sqlx.Connect("sqlite3", "./db.sqlite3")
	DB.MustExec(schema)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database loaded...")
}
