package model_test

import (
	"testing"

	"github.com/Luke-Gurgel/codeflix/domain/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestModel_CreateBank(t *testing.T) {
	code := "001"
	name := "Banco do Brasil"

	bank, err := model.CreateBank(code, name)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(bank.ID))
	require.Equal(t, bank.Code, code)
	require.Equal(t, bank.Name, name)

	_, err = model.CreateBank("", "")
	require.NotNil(t, err)
}
