package grpc

import (
	"context"

	"github.com/Luke-Gurgel/codeflix/application/grpc/pb"
	"github.com/Luke-Gurgel/codeflix/application/usecase"
)

type PixKeyGrpcService struct {
	PixUseCase usecase.PixKeyUseCase
	pb.UnimplementedPixServiceServer
}

func (s *PixKeyGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreationResult, error) {
	key, err := s.PixUseCase.RegisterKey(in.Key, in.Kind, in.AccountID)
	if err != nil {
		return &pb.PixKeyCreationResult{
			Status: "Pix key not registered",
			Error:  err.Error(),
		}, err
	}

	return &pb.PixKeyCreationResult{
		Id:     key.ID,
		Status: "Pix key created",
	}, nil
}

func (s *PixKeyGrpcService) Find(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	pixKey, err := s.PixUseCase.FindKey(in.Key, in.Kind)

	if err != nil {
		return &pb.PixKeyInfo{}, err
	}

	return &pb.PixKeyInfo{
		Id:   pixKey.ID,
		Kind: pixKey.Kind,
		Key:  pixKey.Key,
		Account: &pb.Account{
			AccountID:     pixKey.AccountID,
			AccountNumber: pixKey.Account.Number,
			BankID:        pixKey.Account.BankID,
			BankName:      pixKey.Account.Bank.Name,
			OwnerName:     pixKey.Account.OwnerName,
			CreatedAt:     pixKey.Account.CreatedAt.String(),
		},
		CreatedAt: pixKey.CreatedAt.String(),
	}, nil
}

func CreatePixKeyGrpcService(usecase usecase.PixKeyUseCase) *PixKeyGrpcService {
	return &PixKeyGrpcService{PixUseCase: usecase}
}
