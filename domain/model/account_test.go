package model_test

import (
	"testing"

	"github.com/Luke-Gurgel/codeflix/domain/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestModel_CreateAccount(t *testing.T) {
	code := "001"
	name := "Banco do Brasil"
	bank, _ := model.CreateBank(code, name)

	accountNumber := "abcnumber"
	ownerName := "Luke"
	account, err := model.CreateAccount(bank, accountNumber, ownerName)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(account.ID))
	require.Equal(t, account.Number, accountNumber)
	require.Equal(t, account.Bank.ID, bank.ID)

	_, err = model.CreateAccount(bank, "", ownerName)
	require.NotNil(t, err)
}
