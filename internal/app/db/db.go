package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type DB interface {
	Ping() error
	Close()
}

type PostgresDB struct {
	DB       *sql.DB
	myLogger *zap.Logger
}

var _ DB = &PostgresDB{}

func InitDB(connString string, myLogger *zap.Logger) (*PostgresDB, error) {

	err := RunMigrations(context.Background(), connString, myLogger)
	if err != nil {
		myLogger.Debug("Fail to run migrations", zap.String("msg", err.Error()))
		return nil, err
	}

	db, err := sql.Open("pgx", connString)

	if err != nil {
		myLogger.Debug("Fail open db connection", zap.String("msg", err.Error()))
		return nil, err
	}
	m := &PostgresDB{
		DB:       db,
		myLogger: myLogger,
	}
	myLogger.Debug("InitDB", zap.String("connString", connString))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		myLogger.Debug("Fail to ping db", zap.String("msg", err.Error()))
		return nil, err
	}
	return m, nil
}

func RunMigrations(ctx context.Context, connString string, myLogger *zap.Logger) error {
	const migrationsPath = "./db/migrations"
	fmt.Println(migrationsPath)
	fmt.Println(connString)
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), connString)
	if err != nil {
		return fmt.Errorf("failed to get new migrate instance: %w", err)
	}
	fmt.Println("success")

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to DB: %w", err)
		}
	}
	return nil
}

func (mdb PostgresDB) Ping() error {
	err := mdb.DB.Ping()
	return err
}

func (mdb PostgresDB) Close() {
	mdb.DB.Close()
}

func (mdb PostgresDB) CreateOrGetFromStorage(ctx context.Context, url string) (string, error) {

	shortURL, err := utils.CreateShortURL()
	if err != nil {
		return "", err
	}
	event := models.NewEvent(shortURL, url)
	query := `INSERT INTO public.urls(id, original_url, short_url)
	VALUES ($1, $2, $3);`
	_, err = mdb.DB.ExecContext(ctx, query, event.UUID, event.OriginalURL, event.ShortURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			uniqError := models.NewUniqueError("original_url", err)
			mdb.myLogger.Debug("Failed exec CreateOrGetFromStorage UniqueViolation", zap.String("msg", err.Error()))
			shortURL, err = mdb.isOriginalURLExist(ctx, url)
			if err != nil {
				return "", err
			}
			return shortURL, uniqError
		}
		mdb.myLogger.Debug("Failed exec CreateOrGetFromStorage", zap.String("msg", err.Error()))
		return "", err
	}
	return shortURL, nil
}
func (mdb PostgresDB) GetOriginalURLFromStorage(ctx context.Context, shortURL string) (string, error) {
	mdb.myLogger.Debug("start GetOriginalURLFromStorage")
	query := `SELECT original_url
	FROM public.urls where short_url = $1;`
	var originalURL string
	row := mdb.DB.QueryRowContext(ctx, query, shortURL)

	err := row.Scan(&originalURL)

	if err != nil {
		if err == sql.ErrNoRows {
			mdb.myLogger.Debug("URL not exist", zap.String("msg", err.Error()))
			return "", err
		}
		mdb.myLogger.Debug("Failed to check if url exist", zap.String("msg", err.Error()))
		return "", err
	}

	mdb.myLogger.Debug("Got Original URL", zap.String("msg", originalURL))
	return originalURL, nil
}

func (mdb PostgresDB) isOriginalURLExist(ctx context.Context, url string) (string, error) {
	mdb.myLogger.Debug("start isOriginalURLExist")
	query := `SELECT short_url
	FROM public.urls where original_url = $1;`
	var shortURL string
	row := mdb.DB.QueryRowContext(ctx, query, url)

	err := row.Scan(&shortURL)

	if err != nil {
		if err == sql.ErrNoRows {
			mdb.myLogger.Debug("URL not exist", zap.String("msg", err.Error()))
			return "", nil
		}
		mdb.myLogger.Debug("Failed to check if url exist", zap.String("msg", err.Error()))
		return "", err
	}

	return shortURL, nil
}

func (mdb PostgresDB) CreateOrGetBatchFromStorage(ctx context.Context, batchURL *models.BatchURL) (*models.BatchURL, error) {
	mdb.myLogger.Debug("Start CreateOrGetBatchFromStorage", zap.Any("msg", *(batchURL)))
	tx, err := mdb.DB.Begin()
	if err != nil {
		mdb.myLogger.Debug("Failed to Begin Tx in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))
		return nil, err
	}

	defer tx.Rollback()
	query := `INSERT INTO public.urls(id, original_url, short_url) VALUES ($1, $2, $3);`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		mdb.myLogger.Debug("Failed to PrepareContext in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))
		return nil, err
	}
	defer stmt.Close()

	for k, v := range *batchURL {
		fmt.Println(v.OriginalURL)
		shortURL, err := mdb.isOriginalURLExist(ctx, v.OriginalURL)

		if err != nil {
			mdb.myLogger.Debug("Failed to check if url exist", zap.String("msg", err.Error()))
		}
		if shortURL != "" {
			mdb.myLogger.Debug("URL is exist")
			// return shortURL, nil
			(*batchURL)[k].ShortURL = shortURL
			(*batchURL)[k].OriginalURL = ""
			continue
		}
		mdb.myLogger.Debug("URL not exists")

		shortURL, err = utils.CreateShortURL()
		if err != nil {
			return nil, err
		}

		event := models.NewEvent(shortURL, v.OriginalURL)

		_, err = stmt.ExecContext(ctx, event.UUID, event.OriginalURL, event.ShortURL)
		if err != nil {
			mdb.myLogger.Debug("Failed exec ExecContext in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))
			return nil, err
		}
		(*batchURL)[k].OriginalURL = ""
		(*batchURL)[k].ShortURL = shortURL
	}
	tx.Commit()
	return batchURL, nil
}

func (mdb PostgresDB) GetUserByID(ctx context.Context, ID uint64) (*auth.User, error) {

	return nil, fmt.Errorf("not implemented")
}
