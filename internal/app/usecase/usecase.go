package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type Repository interface {
	CreateOrGetFromStorage(ctx context.Context, url string, userID int) (string, error)
	GetOriginalURLFromStorage(ctx context.Context, url string) (string, error)
	CreateOrGetBatchFromStorage(ctx context.Context, batchURL *models.BatchURL, userID int) (*models.BatchURL, error)
	RegisterUser(ctx context.Context) (*models.User, error)
	GetBatchURLFromStorage(ctx context.Context, userID int) (*models.BatchURL, error)
	DeleteSliceURLFromStorage(ctx context.Context, shortURL []string, userID int) error
	GetStatsFromStorage(ctx context.Context) (*models.Stats, error)

	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

func ProcessBatchURLs(ctx context.Context,
	body []byte,
	repo Repository,
	token, globalURL string,
	l *zap.Logger) (models.BatchURL, error) {
	userID, _ := auth.GetUserID(token)
	l.Debug("start ProcessBatchURLs")

	var payload *models.BatchURL
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if len(*payload) < 1 {
		//nolint:goerr113
		return nil, errors.New("empty payload")
	}

	val, err := repo.CreateOrGetBatchFromStorage(ctx, payload, userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	for k := range *val {
		(*val)[k].ShortURL = utils.ReturnURL((*val)[k].ShortURL, globalURL)
	}

	return *val, nil
}
