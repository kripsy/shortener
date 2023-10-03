package proto

import (
	context "context"
	"log"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/middleware"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	UnimplementedGreeterServer
	MyLogger *zap.Logger
	URLRepo  string
}

func InitServer(urlRepo string, repo handlers.Repository,
	myLogger *zap.Logger,
	ts *net.IPNet,
	isSecure bool) (*grpc.Server, error) {
	s := &Server{
		MyLogger: myLogger,
		URLRepo:  urlRepo,
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

	RegisterGreeterServer(srv, s)

	return srv, nil
}

func (s *Server) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	return &HelloResponse{Message: "Hello, " + req.Name}, nil
}
