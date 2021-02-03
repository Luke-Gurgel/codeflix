package model_test

import (
	"testing"

	"github.com/Luke-Gurgel/codeflix/domain/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestModel_CreatePixKey(t *testing.T) {
	code := "001"
	name := "Banco do Brasil"
	bank, _ := model.CreateBank(code, name)

	accountNumber := "abcnumber"
	ownerName := "Luke"
	account, _ := model.CreateAccount(bank, accountNumber, ownerName)

	kind := "email"
	key := "j@j.com"
	pixKey, err := model.CreatePixKey(kind, account, key)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(pixKey.ID))
	require.Equal(t, pixKey.Kind, kind)
	require.Equal(t, pixKey.Status, "active")

	kind = "cpf"
	_, err = model.CreatePixKey(kind, account, key)
	require.Nil(t, err)

	_, err = model.CreatePixKey("invalid_kind", account, key)
	require.NotNil(t, err)
}
