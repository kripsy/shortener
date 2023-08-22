package inmemorystorage

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type TestParams struct {
	testLogger     *zap.Logger
	testPrefixAddr string
	TestStorage    map[string]models.Event
}

func getParamsForTest() *TestParams {
	tl, _ := logger.InitLog("Debug")

	tp := &TestParams{
		testLogger:     tl,
		testPrefixAddr: "http://localhost:8080",
		TestStorage: map[string]models.Event{
			"1": {
				UUID:          1,
				ShortURL:      "ShortURL1",
				OriginalURL:   "OriginalURL1",
				CorrelationID: "1",
				UserID:        1,
			},
			"2": {
				UUID:          2,
				ShortURL:      "ShortURL2",
				OriginalURL:   "OriginalURL2",
				CorrelationID: "2",
				UserID:        1,
			},
			"3": {
				UUID:          3,
				ShortURL:      "ShortURL3",
				OriginalURL:   "OriginalURL3",
				CorrelationID: "3",
				UserID:        2,
			},
		},
	}
	return tp
}

func TestGetUserByID(t *testing.T) {
	paramTest := getParamsForTest()

	type fields struct {
		storage  map[string]models.Event
		myLogger *zap.Logger
	}
	type args struct {
		ctx context.Context
		ID  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "first success getting user",
			fields: fields{
				storage:  paramTest.TestStorage,
				myLogger: paramTest.testLogger,
			},
			args: args{
				ctx: context.Background(),
				ID:  2,
			},
			want: &models.User{
				ID: 2,
			},
			wantErr: false,
		},
		{
			name: "first failed getting user",
			fields: fields{
				storage:  paramTest.TestStorage,
				myLogger: paramTest.testLogger,
			},
			args: args{
				ctx: context.Background(),
				ID:  100500,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := InMemoryStorage{
				storage:  tt.fields.storage,
				myLogger: tt.fields.myLogger,
			}
			got, err := m.GetUserByID(tt.args.ctx, tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryStorage.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryStorage.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	paramTest := getParamsForTest()

	type fields struct {
		storage  map[string]models.Event
		myLogger *zap.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "first success register user",
			fields: fields{
				storage:  paramTest.TestStorage,
				myLogger: paramTest.testLogger,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := InMemoryStorage{
				storage:  tt.fields.storage,
				myLogger: tt.fields.myLogger,
			}
			got, err := m.RegisterUser(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryStorage.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got)
		})
	}
}

func BenchmarkCreateOrGetFromStorageWithoutPointer(b *testing.B) {
	paramTest := getParamsForTest()
	m := InMemoryStorage{
		storage:  paramTest.TestStorage,
		myLogger: paramTest.testLogger,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.CreateOrGetFromStorageWithoutPointer(context.Background(), fmt.Sprintf("http://example.com/%d", i+1), 1)
	}
}
func BenchmarkCreateOrGetFromStorage(b *testing.B) {
	paramTest := getParamsForTest()
	m := InMemoryStorage{
		storage:  paramTest.TestStorage,
		myLogger: paramTest.testLogger,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.CreateOrGetFromStorage(context.Background(), fmt.Sprintf("http://example.com/%d", i+1), 1)
	}
}

// func BenchmarkMemoryStorageCreateOrGetFromStorage(b *testing.B) {
// 	paramTest := getParamsForTest()
// 	m := InMemoryStorage{
// 		storage:  paramTest.TestStorage,
// 		myLogger: paramTest.testLogger,
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		events := GenerateEvents(100) // Генерирует 100 событий
// 		m.CreateOrGetBatchFromStorage(context.Background(), &events, 1)
// 	}
// }

// GenerateEvents создает множество событий Event
func GenerateEvents(count int) models.BatchURL {
	events := make(models.BatchURL, count)

	for i := 0; i < count; i++ {
		events[i] = models.Event{
			UUID:          i + 1, // Уникальный идентификатор для каждого события
			ShortURL:      "",
			OriginalURL:   fmt.Sprintf("http://example.com/%d", i+1),
			CorrelationID: fmt.Sprintf("correlation_id_%d", i+1),
			UserID:        rand.Intn(100) + 1, // Произвольный UserID в диапазоне от 1 до 100
			IsDeleted:     false,
		}
	}

	return events
}

// GenerateEvents создает множество событий Event
func GenerateURL(count int) string {
	return fmt.Sprintf("http://example.com/%d", count+1)
}
