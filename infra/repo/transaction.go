package repo

import (
	"fmt"

	"github.com/Luke-Gurgel/codeflix/domain/model"
	"gorm.io/gorm"
)

type TransactionRepoDB struct {
	DB *gorm.DB
}

func (r TransactionRepoDB) Register(transaction *model.Transaction) error {
	err := r.DB.Create(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

func (r TransactionRepoDB) Save(transaction *model.Transaction) error {
	err := r.DB.Save(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

func (r TransactionRepoDB) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	r.DB.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)

	if transaction.ID == "" {
		return nil, fmt.Errorf("Not transaction with that ID was found.")
	}

	return &transaction, nil
}
