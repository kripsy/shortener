package proto

import (
	context "context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/middleware"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/usecase"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	status "google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedShortenerServer
	MyLogger *zap.Logger
	URLRepo  string
	Repo     handlers.Repository
}

func InitServer(urlRepo string, repo handlers.Repository,
	myLogger *zap.Logger,
	ts *net.IPNet,
	isSecure bool) (*grpc.Server, error) {
	s := &Server{
		MyLogger: myLogger,
		URLRepo:  urlRepo,
		Repo:     repo,
	}

	m := middleware.InitMyMiddleware(myLogger, repo, ts)
	interceptors := []grpc.UnaryServerInterceptor{
		m.GrpcRequestLogger,
		m.GrpcTrustedSubnetMiddleware,
		m.GrpcJWTInterceptor,
	}

	var srv *grpc.Server
	if isSecure {
		creds, err := credentials.NewServerTLSFromFile(utils.ServerCertPath, utils.PrivateKeyPath)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		srv = grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	} else {
		srv = grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	}

	RegisterShortenerServer(srv, s)

	return srv, nil
}

func (s *Server) SaveURL(ctx context.Context, req *SaveURLRequest) (*SaveURLResponse, error) {
	token, err := utils.GetTokenFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	userID, _ := auth.GetUserID(token)
	s.MyLogger.Debug("start SaveURL")
	isUniqueError := false
	val, err := s.Repo.CreateOrGetFromStorage(ctx, req.Url, userID)
	if err != nil {
		var ue *models.UniqueError
		if errors.As(err, &ue) {
			isUniqueError = true
		} else {
			s.MyLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	result := utils.ReturnURL(val, s.URLRepo)
	return &SaveURLResponse{
		Result:        result,
		IsUniqueError: isUniqueError,
	}, nil
}

func (s *Server) GetURL(ctx context.Context, req *GetURLRequest) (*GetURLResponse, error) {
	s.MyLogger.Debug("start SaveURL")
	val, err := s.Repo.GetOriginalURLFromStorage(ctx, req.Url)
	if err != nil {
		s.MyLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result := utils.ReturnURL(val, s.URLRepo)

	return &GetURLResponse{
		Url: result,
	}, nil
}

func (s *Server) GetStats(ctx context.Context, req *GetStatsRequest) (*GetStatsResponse, error) {
	s.MyLogger.Debug("start SaveURL")
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	stats, err := s.Repo.GetStatsFromStorage(ctx)
	if err != nil {

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetStatsResponse{
		Urls:  int32(stats.URLs),
		Users: int32(stats.Users),
	}, nil
}

func (s *Server) SaveBatchURL(ctx context.Context, req *SaveBatchURLRequest) (*SaveBatchURLResponse, error) {
	token, err := utils.GetTokenFromMetadata(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, _ = auth.GetUserID(token)

	body, err := json.Marshal(req.UrlBatch)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	result, err := usecase.ProcessBatchURLs(ctx, body, s.Repo, token, s.URLRepo, s.MyLogger)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	r := сonvertToSaveBatchURLResponse(result)

	s.MyLogger.Debug("start SaveURL")

	return r, nil
}

func сonvertToSaveBatchURLResponse(batchURL models.BatchURL) *SaveBatchURLResponse {
	var responseObjects []*SaveBatchURLResponse_URLObject

	for _, event := range batchURL {
		responseObject := &SaveBatchURLResponse_URLObject{
			CorrelationId: event.CorrelationID,
			ShortUrl:      event.ShortURL,
		}
		responseObjects = append(responseObjects, responseObject)
	}

	return &SaveBatchURLResponse{
		UrlBatch: responseObjects,
	}
}
