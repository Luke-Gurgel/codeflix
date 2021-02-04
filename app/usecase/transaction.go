package app

import (
	"github.com/Luke-Gurgel/codeflix/domain/model"
)

type TransactionUseCase struct {
	PixKeyRepository      model.PixKeyRepositoryInterface
	TransactionRepository model.TransactionRepositoryInterface
}

func (t *TransactionUseCase) Register(accountID string, amount float64, pixKeyTo string, pixKeyToKind string, description string) (*model.Transaction, error) {
	account, err := t.PixKeyRepository.FindAccount(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixKeyRepository.FindKeyByKind(pixKeyTo, pixKeyToKind)
	if err != nil {
		return nil, err
	}

	transaction, err := model.CreateTransaction(account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionUseCase) Confirm(transactionID string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionConfirmed

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionUseCase) Complete(transactionID string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionCompleted

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionUseCase) Error(transactionID string, reason string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionConfirmed
	transaction.CancelDescription = reason

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
