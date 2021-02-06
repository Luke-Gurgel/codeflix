package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/Luke-Gurgel/codeflix/application/grpc/pb"
	"github.com/Luke-Gurgel/codeflix/application/usecase"
	"github.com/Luke-Gurgel/codeflix/infra/repo"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(db *gorm.DB, port int) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pixRepo := repo.PixKeyRepoDB{DB: db}
	pixUseCase := usecase.PixKeyUseCase{PixKeyRepository: pixRepo}
	pixGrpcService := CreatePixKeyGrpcService(pixUseCase)
	pb.RegisterPixServiceServer(grpcServer, pixGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Could not listen", err)
	}

	log.Printf("gRPC server listening on port %d", port)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Could not start grpc server", err)
	}
}
