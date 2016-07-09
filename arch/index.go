package archaeology

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type IndexDB struct {
	Db *sql.DB
}

func (index *IndexDB) mustQuery(stmt string) *sql.Rows {
	rows, err := index.Db.Query(stmt)
	if nil != err {
		log.Fatal(err)
	}
	return rows
}

func (index *IndexDB) Open(cfg ArchCfg) error {
	db, err := sql.Open("sqlite3", cfg.IndexDbPath)
	if err != nil {
		log.Fatal(err)
	}
	index.Db = db
	log.Print("Opened ", cfg.IndexDbPath)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (index *IndexDB) Close() error {
	return index.Db.Close()
}

func (index *IndexDB) HasTable(name string) bool {
	sqlStmt := `show tables in `
	// 	sqlStmt := `SELECT TABLE_NAME
	// FROM INFORMATION_SCHEMA.TABLES
	// WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_CATALOG='dbName'`

	rows := index.mustQuery(sqlStmt)
	defer rows.Close()

	for rows.Next() {
		var tableName *string
		err := rows.Scan(tableName)
		if nil != err {
			log.Fatal(err)
		}
		if *tableName == name {
			return true
		}
	}
	return false
}

func (index *IndexDB) mustExec(stmt string) {
	_, err := index.Db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func (index *IndexDB) CreateTables() {
	// Check 'hashes' table
	if !index.HasTable("hashes") {
		sqlStmt := `
		create table hashes (VersionId INT primary key, Hash TEXT);
		`
		index.mustExec(sqlStmt)
	}

	if !index.HasTable("versions") {
		sqlStmt := `
		create table versions (Path TEXT primary key, VersionId INT);
		`
		index.mustExec(sqlStmt)
	}

}
