package model_test

import (
	"testing"

	"github.com/Luke-Gurgel/codeflix/domain/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestModel_CreateTransaction(t *testing.T) {
	code := "001"
	name := "Banco do Brasil"
	bank, _ := model.CreateBank(code, name)

	srcAccountNumber := "abcnumber"
	srcAccountOwnerName := "Luke"
	srcAccount, _ := model.CreateAccount(bank, srcAccountNumber, srcAccountOwnerName)

	destAccountNumber := "xyznumber"
	destAccountOwnerName := "Thais"
	destAccount, _ := model.CreateAccount(bank, destAccountNumber, destAccountOwnerName)

	require.NotEqual(t, srcAccount.ID, destAccount.ID)

	kind := "email"
	key := "j@j.com"
	pixKey, _ := model.CreatePixKey(kind, destAccount, key)

	amount := 3.10
	description := "Paying my bills yall"
	transaction, err := model.CreateTransaction(srcAccount, pixKey, amount, "", description)

	require.Nil(t, err)
	require.NotNil(t, uuid.FromStringOrNil(transaction.ID))
	require.Equal(t, transaction.Status, model.TransactionPending)
	require.Equal(t, transaction.Description, description)
	require.Empty(t, transaction.CancelDescription)
	require.Equal(t, transaction.Amount, amount)

	pixKeySameAccount, _ := model.CreatePixKey(kind, srcAccount, key)

	_, err = model.CreateTransaction(srcAccount, pixKeySameAccount, amount, "", description)
	require.NotNil(t, err)

	_, err = model.CreateTransaction(srcAccount, pixKey, 0, "", description)
	require.NotNil(t, err)
}

func TestModel_ChangeTransactionStatus(t *testing.T) {
	code := "001"
	name := "Banco do Brasil"
	bank, _ := model.CreateBank(code, name)

	srcAccountNumber := "abcnumber"
	srcAccountOwnerName := "Luke"
	srcAccount, _ := model.CreateAccount(bank, srcAccountNumber, srcAccountOwnerName)

	destAccountNumber := "xyznumber"
	destAccountOwnerName := "Thais"
	destAccount, _ := model.CreateAccount(bank, destAccountNumber, destAccountOwnerName)

	kind := "email"
	key := "j@j.com"
	pixKey, _ := model.CreatePixKey(kind, destAccount, key)

	amount := 3.10
	description := "Paying my bills yall"
	transaction, _ := model.CreateTransaction(srcAccount, pixKey, amount, "", description)

	_ = transaction.Complete()
	require.Equal(t, transaction.Status, model.TransactionCompleted)

	cancelDescription := "I'm actually broke"
	_ = transaction.Cancel(cancelDescription)
	require.Equal(t, transaction.Status, model.TransactionError)
	require.Equal(t, transaction.CancelDescription, cancelDescription)
}
