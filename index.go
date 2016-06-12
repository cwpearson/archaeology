package archaeology

import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type IndexDB struct {
	Db *sql.DB
}

func (index *IndexDB) Open(cfg ArchCfg) error {
	db, err := sql.Open("sqlite3", cfg.IndexDbPath)
	if err != nil {
		log.Fatal(err)
	}
	index.Db = db
	log.Print("Opened ", cfg.IndexDbPath)

	return err
}

func (index *IndexDB) Close() error {
	return index.Db.Close()
}

func (index *IndexDB) CreateTable() error {
	sqlStmt := `
	create table files (path TEXT primary key, hash TEXT);
	delete from files;
	`

	_, err := index.Db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	return err;
}

