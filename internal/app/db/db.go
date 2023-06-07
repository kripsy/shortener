package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB interface {
	Ping() error
	Close()
}

type PostgresDB struct {
	DB *sql.DB
}

var _ DB = &PostgresDB{}

// func InitDB(host, port, user, password, dbName string) (*PostgresDB, error) {
// 	ps := fmt.Sprintf("host=%s user=%s password=%s sslmode=disable port=%s",
// 		host, user, password, port)

// 	db, err := sql.Open("pgx", ps)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m := &PostgresDB{
// 		DB: db,
// 	}

//		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//		defer cancel()
//		if err = db.PingContext(ctx); err != nil {
//			return nil, err
//		}
//		return m, nil
//	}
func InitDB(connString string) (*PostgresDB, error) {

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}
	m := &PostgresDB{
		DB: db,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return m, nil
}

func (mdb PostgresDB) Ping() error {
	err := mdb.DB.Ping()
	return err
}

func (mdb PostgresDB) Close() {
	mdb.DB.Close()
}
