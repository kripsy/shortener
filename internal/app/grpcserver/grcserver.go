package grpcserver

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/bufbuild/protovalidate-go"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/middleware"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/usecase"
	"github.com/kripsy/shortener/internal/app/utils"
	pb "github.com/kripsy/shortener/pkg/api/shortener/v1/gen"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedShortenerServer
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

	pb.RegisterShortenerServer(srv, s)

	return srv, nil
}

func (s *Server) SaveURL(ctx context.Context, req *pb.SaveURLRequest) (*pb.SaveURLResponse, error) {
	token, err := utils.GetTokenFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
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

			return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
		}
	}

	result := utils.ReturnURL(val, s.URLRepo)

	return &pb.SaveURLResponse{
		Result:        result,
		IsUniqueError: isUniqueError,
	}, nil
}

func (s *Server) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	s.MyLogger.Debug("start SaveURL")

	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
	}

	val, err := s.Repo.GetOriginalURLFromStorage(ctx, req.Url)
	if err != nil {
		s.MyLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))

		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
	}

	result := utils.ReturnURL(val, s.URLRepo)

	return &pb.GetURLResponse{
		Url: result,
	}, nil
}

func (s *Server) GetStats(ctx context.Context, _ *emptypb.Empty) (*pb.GetStatsResponse, error) {
	s.MyLogger.Debug("start SaveURL")
	stats, err := s.Repo.GetStatsFromStorage(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	return &pb.GetStatsResponse{
		Urls:  int32(stats.URLs),
		Users: int32(stats.Users),
	}, nil
}

func (s *Server) SaveBatchURL(ctx context.Context, req *pb.SaveBatchURLRequest) (*pb.SaveBatchURLResponse, error) {
	token, err := utils.GetTokenFromMetadata(ctx)

	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
	}
	_, _ = auth.GetUserID(token)

	body, err := json.Marshal(req.UrlBatch)
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, "invalid request"))
	}

	batch := models.BatchURL{}
	err = json.Unmarshal(body, &batch)
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, "invalid request"))
	}
	// for _, v := range req.UrlBatch {
	// 	batch = append(batch, models.Event{
	// 		CorrelationID: v.CorrelationId,
	// 		OriginalURL:   v.OriginalUrl,
	// 	})
	// }

	result, err := usecase.ProcessBatchURLs(ctx, &batch, s.Repo, token, s.URLRepo, s.MyLogger)
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}
	r := utils.Convert2BatchURLResponse(result)

	s.MyLogger.Debug("start SaveURL")

	return r, nil
}
