package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending string = "pending"
	TransactionCompleted string = "completed"
	TransactionConfirmed string = "confirmed"
	TransactionError string = "error"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base `valid:"required"`
	PixKeyTo *PixKey `valid:"-"`
	AccountFrom *Account `valid:"-"`
	Status string `json:"status" valid:"notnull"`
	Amount float64 `json:"amount" valid:"notnull"`
	Description string `json:"description" valid:"notnull"`
	CancelDescription string `json:"cancel_description" valid:"-"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)
	
	if t.Amount <= 0 {
		return errors.New("Amount must be greater than 0.")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionConfirmed && t.Status != TransactionError {
		return errors.New("Unknown transaction status.")
	}

	if t.PixKeyTo.ID == t.AccountFrom.ID {
		return errors.New("Source and destination accounts must be different.")
	}

	if err != nil {
		return err 
	}
	
	return nil
}

func CreateTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction {
		Amount: amount,
		PixKeyTo: pixKeyTo,
		AccountFrom: accountFrom,
		Description: description,
		Status: TransactionPending,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()
	if (err != nil) {
		return nil, err
	}

	return &transaction, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Cancel() error {
	t.Status = TransactionError
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Confirm(description string) error {
	t.Status = TransactionConfirmed
	t.CancelDescription = description
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}