package repository

import (
	"database/sql"
	"path/filepath"
	"runtime"
)

type SQLiteContext struct {
	connectionString string
	databaseKind     string
	dbConnection     *sql.DB
}

func (ins *SQLiteContext) Init() {
	_, b, _, _ := runtime.Caller(0)
	ins.connectionString = filepath.Join(filepath.Dir(b), "..", "database.db")
	ins.databaseKind = "sqlite3"
}

func (ins *SQLiteContext) Open() (*sql.DB, error) {
	var err error
	ins.dbConnection, err = sql.Open(ins.databaseKind, ins.connectionString)
	if err != nil {
		return nil, err
	}
	return ins.dbConnection, nil
}

var DbConn SQLiteContext
