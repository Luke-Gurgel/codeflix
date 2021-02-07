package factory

import (
	"github.com/Luke-Gurgel/codeflix/application/usecase"
	"github.com/Luke-Gurgel/codeflix/infra/repo"
	"github.com/jinzhu/gorm"
)

func TransactionUseCase(database *gorm.DB) usecase.TransactionUseCase {
	pixRepo := repo.PixKeyRepoDB{DB: database}
	transactionRepo := repo.TransactionRepoDB{DB: database}
	transactionUseCase := usecase.TransactionUseCase{
		PixKeyRepository:      pixRepo,
		TransactionRepository: &transactionRepo,
	}

	return transactionUseCase
}
