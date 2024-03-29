// Package db provides functionality for working with database storage.
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"

	//nolint:revive,nolintlint
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	//nolint:revive,nolintlint
	_ "github.com/jackc/pgx/v5/stdlib"
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

func InitDB(_ context.Context, connString string, myLogger *zap.Logger) (*PostgresDB, error) {
	//nolint:contextcheck
	err := RunMigrations(context.Background(), connString, myLogger)
	if err != nil {
		myLogger.Debug("Fail to run migrations", zap.String("msg", err.Error()))

		return nil, err
	}

	db, err := sql.Open("pgx", connString)

	if err != nil {
		myLogger.Debug("Fail open db connection", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}
	m := &PostgresDB{
		DB:       db,
		myLogger: myLogger,
	}
	myLogger.Debug("InitDB", zap.String("connString", connString))
	//nolint:contextcheck
	if err = db.PingContext(context.Background()); err != nil {
		myLogger.Debug("Fail to ping db", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	return m, nil
}

func RunMigrations(_ context.Context, connString string, _ *zap.Logger) error {
	const migrationsPath = "./db/migrations"
	// fmt.Println(migrationsPath)
	// fmt.Println(connString)
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
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (mdb PostgresDB) Close() {
	mdb.DB.Close()
}

func (mdb PostgresDB) CreateOrGetFromStorage(ctx context.Context, url string, userID int) (string, error) {
	shortURL, err := utils.CreateShortURL()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	event := models.NewEvent(shortURL, url, userID)
	query := `INSERT INTO public.urls(id, original_url, short_url, user_id)
	VALUES ($1, $2, $3, $4);`
	_, err = mdb.DB.ExecContext(ctx, query, event.UUID, event.OriginalURL, event.ShortURL, event.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			uniqError := models.NewUniqueError("original_url", err)
			mdb.myLogger.Debug("Failed exec CreateOrGetFromStorage UniqueViolation", zap.String("msg", err.Error()))
			shortURL, err = mdb.isOriginalURLExist(ctx, url)
			if err != nil {
				return "", err
			}

			// nolint:wrapcheck,nolintlint
			return shortURL, uniqError
		}
		mdb.myLogger.Debug("Failed exec CreateOrGetFromStorage", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}

	return shortURL, nil
}
func (mdb PostgresDB) GetOriginalURLFromStorage(ctx context.Context, shortURL string) (string, error) {
	mdb.myLogger.Debug("start GetOriginalURLFromStorage")
	query := `SELECT original_url, is_deleted
	FROM public.urls where short_url = $1;`
	var originalURL string
	var isDeleted sql.NullBool
	row := mdb.DB.QueryRowContext(ctx, query, shortURL)

	err := row.Scan(&originalURL, &isDeleted)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			mdb.myLogger.Debug("URL not exist", zap.String("msg", err.Error()))

			return "", fmt.Errorf("%w", err)
		}
		mdb.myLogger.Debug("Failed to check if url exist", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}

	if isDeleted.Valid && isDeleted.Bool {
		mdb.myLogger.Debug("Got Original URL", zap.String("msg", originalURL), zap.Bool("is deleted?", isDeleted.Bool))

		//nolint:goerr113,wrapcheck
		return "", models.NewIsDeletedError(shortURL, models.NewIsDeletedError(shortURL, errors.New("")))
	}

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
		if errors.Is(err, sql.ErrNoRows) {
			mdb.myLogger.Debug("URL not exist", zap.String("msg", err.Error()))

			return "", nil
		}
		mdb.myLogger.Debug("Failed to check if url exist", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}

	return shortURL, nil
}

func (mdb PostgresDB) CreateOrGetBatchFromStorage(ctx context.Context,
	batchURL *models.BatchURL,
	userID int) (*models.BatchURL, error) {
	mdb.myLogger.Debug("Start CreateOrGetBatchFromStorage", zap.Any("msg", *(batchURL)))
	tx, err := mdb.DB.Begin()
	if err != nil {
		mdb.myLogger.Debug("Failed to Begin Tx in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}
	mdb.myLogger.Debug("UserID", zap.Int("msg", userID))

	defer func() {
		if err = tx.Rollback(); err != nil {
			mdb.myLogger.Debug("error in rollback", zap.Error(err))
		}
	}()
	query := `INSERT INTO public.urls(id, original_url, short_url, user_id) VALUES ($1, $2, $3, $4);`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		mdb.myLogger.Debug("Failed to PrepareContext in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}
	defer stmt.Close()

	for k, v := range *batchURL {
		var shortURL string
		shortURL, err = mdb.isOriginalURLExist(ctx, v.OriginalURL)

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
			return nil, fmt.Errorf("%w", err)
		}

		event := models.NewEvent(shortURL, v.OriginalURL, userID)

		_, err = stmt.ExecContext(ctx, event.UUID, event.OriginalURL, event.ShortURL, event.UserID)
		if err != nil {
			mdb.myLogger.Debug("Failed exec ExecContext in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))

			return nil, fmt.Errorf("%w", err)
		}
		(*batchURL)[k].OriginalURL = ""
		(*batchURL)[k].ShortURL = shortURL
	}
	if err = tx.Commit(); err != nil {
		mdb.myLogger.Debug("Failed tx.Commit in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	return batchURL, nil
}

func (mdb PostgresDB) GetUserByID(_ context.Context, _ int) (*models.User, error) {
	//nolint:goerr113
	return nil, fmt.Errorf("not implemented")
}

func (mdb PostgresDB) getNextUserID(ctx context.Context) (int, error) {
	tx, err := mdb.DB.Begin()
	var userID int
	if err != nil {
		mdb.myLogger.Debug("Failed to Begin Tx in GetNextUserID", zap.String("msg", err.Error()))

		return -1, fmt.Errorf("%w", err)
	}

	defer func() {
		if err = tx.Rollback(); err != nil {
			mdb.myLogger.Debug("error in rollback", zap.Error(err))
		}
	}()
	query := `SELECT CASE
			WHEN count(id)<1 THEN 1
			ELSE max(id)+1
			END 
  		FROM users;`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		mdb.myLogger.Debug("Failed to PrepareContext in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))

		return -1, fmt.Errorf("%w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx)
	err = row.Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			mdb.myLogger.Debug("URL not exist", zap.String("msg", err.Error()))

			return -1, fmt.Errorf("%w", err)
		}
		mdb.myLogger.Debug("Failed to check if url exist", zap.String("msg", err.Error()))

		return -1, fmt.Errorf("%w", err)
	}

	if err = tx.Commit(); err != nil {
		mdb.myLogger.Debug("Failed tx.Commit in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))

		return -1, fmt.Errorf("%w", err)
	}

	return userID, nil
}

func (mdb PostgresDB) RegisterUser(ctx context.Context) (*models.User, error) {
	newUserID, err := mdb.getNextUserID(ctx)
	if err != nil {
		mdb.myLogger.Debug("Failed to getNextUserID in RegisterUser", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}
	mdb.myLogger.Debug("newUserID", zap.Int("msg", newUserID))

	tx, err := mdb.DB.Begin()

	if err != nil {
		mdb.myLogger.Debug("Failed to Begin Tx in RegisterUser", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	defer func() {
		if err = tx.Rollback(); err != nil {
			mdb.myLogger.Debug("error in rollback", zap.Error(err))
		}
	}()
	query := `INSERT INTO users (id) values ($1);`
	stmt, err := tx.PrepareContext(ctx, query)

	if err != nil {
		mdb.myLogger.Debug("Failed to PrepareContext in CreateOrGetBatchFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, newUserID)
	if err != nil {
		mdb.myLogger.Debug("Failed exec ExecContext in RegisterUser", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}
	if err = tx.Commit(); err != nil {
		mdb.myLogger.Debug("Failed tx.Commit in RegisterUser", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	return &models.User{
		ID: newUserID,
	}, nil
}

func (mdb PostgresDB) GetBatchURLFromStorage(ctx context.Context, userID int) (*models.BatchURL, error) {
	batchURL := &models.BatchURL{}
	tx, err := mdb.DB.Begin()
	if err != nil {
		mdb.myLogger.Debug("Failed to Begin Tx in GetBatchURLFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			mdb.myLogger.Debug("error in rollback", zap.Error(err))
		}
	}()

	query := `SELECT short_url, original_url from urls where user_id = $1`

	stmt, err := tx.PrepareContext(ctx, query)

	if err != nil {
		mdb.myLogger.Debug("Failer to PrepareContext in GetBatchURLFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		mdb.myLogger.Debug("Failer to QueryContext in GetBatchURLFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var shortURL, originalURL string
		err = rows.Scan(&shortURL, &originalURL)
		if err != nil {
			mdb.myLogger.Debug("Failer to Scan in GetBatchURLFromStorage", zap.String("msg", err.Error()))

			return nil, fmt.Errorf("%w", err)
		}

		event := &models.Event{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		}
		*batchURL = append(*batchURL, *event)
	}
	err = rows.Err()
	if err != nil {
		mdb.myLogger.Debug("Failer to Scan (rows.Err()) in GetBatchURLFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	return batchURL, nil
}

func (mdb PostgresDB) DeleteSliceURLFromStorage(ctx context.Context, shortURL []string, userID int) error {
	mdb.myLogger.Debug("started DeleteSliceURLFromStorage")
	mdb.myLogger.Debug("shortURL in DeleteSliceURLFromStorage", zap.Any("msg", shortURL))
	mdb.myLogger.Debug("userID in DeleteSliceURLFromStorage", zap.Int("msg", userID))

	tx, err := mdb.DB.Begin()
	if err != nil {
		mdb.myLogger.Debug("Failed to Begin Tx in DeleteSliceURLFromStorage", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}

	defer func() {
		if err = tx.Rollback(); err != nil {
			mdb.myLogger.Debug("error in rollback", zap.Error(err))
		}
	}()

	urls := squirrel.Update("urls").
		Set("is_deleted", true).
		Where(squirrel.And{
			squirrel.Eq{"short_url": shortURL},
			squirrel.Eq{"user_id": userID},
			squirrel.Eq{"is_deleted": nil}}).
		PlaceholderFormat(squirrel.Dollar)
	sql, args, err := urls.ToSql()

	if err != nil {
		mdb.myLogger.Debug("Failed to build sql usersID", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}
	// fmt.Println(sql)
	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		mdb.myLogger.Debug("Failed to exec sql", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}

	if err = tx.Commit(); err != nil {
		mdb.myLogger.Debug("Failed tx.Commit in DeleteSliceURLFromStorage", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}
	mdb.myLogger.Debug("Success commit DeleteSliceURLFromStorage")

	return nil
}

func (mdb PostgresDB) GetStatsFromStorage(ctx context.Context) (*models.Stats, error) {
	tx, err := mdb.DB.Begin()
	if err != nil {
		mdb.myLogger.Debug("Failed to Begin Tx in GetStatsFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	defer func() {
		if err = tx.Rollback(); err != nil {
			mdb.myLogger.Debug("error in rollback", zap.Error(err))
		}
	}()

	usersSubquery, _, _ := squirrel.Select("COUNT (DISTINCT id)").From("users").PlaceholderFormat(squirrel.Dollar).ToSql()

	urlsSubquery, _, _ := squirrel.Select("COUNT (DISTINCT id)").From("urls").PlaceholderFormat(squirrel.Dollar).ToSql()

	query, _, err := squirrel.Select(fmt.Sprintf("(%s) AS user_count",
		usersSubquery),
		fmt.Sprintf("(%s) AS url_count",
			urlsSubquery)).ToSql()

	if err != nil {
		mdb.myLogger.Debug("Failed to generate sql command", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	stmt, err := tx.PrepareContext(ctx, query)

	if err != nil {
		mdb.myLogger.Debug("Failed to PrepareContext in GetStatsFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx)
	var countUrls, countUsers int
	err = row.Scan(&countUsers, &countUrls)
	if err != nil {
		mdb.myLogger.Debug("failed to scan in GetStatsFromStorage", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	return &models.Stats{
		URLs:  countUrls,
		Users: countUsers,
	}, nil
}
