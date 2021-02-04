package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionConfirmed string = "confirmed"
	TransactionError     string = "error"
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
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	AccountFromID     string   `gorm:"column:account_from_id;type:uuid;" valid:"notnull"`
	Amount            float64  `json:"amount" gorm:"type:float" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	PixKeyToID        string   `gorm:"column:pix_key_to_id;type:uuid;" valid:"notnull"`
	Status            string   `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description       string   `json:"description" gorm:"type:varchar(255)" valid:"-"`
	CancelDescription string   `json:"cancel_description" gorm:"type:varchar(255)" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("Amount must be greater than 0.")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionConfirmed && t.Status != TransactionError {
		return errors.New("Unknown transaction status.")
	}

	if t.PixKeyTo.AccountID == t.AccountFromID {
		return errors.New("Source and destination accounts must be different.")
	}

	if err != nil {
		return err
	}

	return nil
}

func CreateTransaction(accountFrom *Account, pixKeyTo *PixKey, amount float64, id string, description string) (*Transaction, error) {
	transaction := Transaction{
		Amount:        amount,
		PixKeyTo:      pixKeyTo,
		PixKeyToID:    pixKeyTo.ID,
		AccountFrom:   accountFrom,
		AccountFromID: accountFrom.ID,
		Status:        TransactionPending,
		Description:   description,
	}

	transaction.ID = id
	transaction.CreatedAt = time.Now()

	if id == "" {
		transaction.ID = uuid.NewV4().String()
	}

	err := transaction.isValid()
	if err != nil {
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

func (t *Transaction) Cancel(description string) error {
	t.CancelDescription = description
	t.Status = TransactionError
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}
