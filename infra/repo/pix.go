package repo

import (
	"fmt"

	"github.com/Luke-Gurgel/codeflix/domain/model"
	"gorm.io/gorm"
)

type PixKeyRepoDB struct {
	DB *gorm.DB
}

func (r PixKeyRepoDB) AddBank(bank *model.Bank) error {
	err := r.DB.Create(bank).Error
	if err != nil {
		return err
	}
	return nil
}

func (r PixKeyRepoDB) AddAccount(account *model.Account) error {
	err := r.DB.Create(account).Error
	if err != nil {
		return err
	}
	return nil
}

func (r PixKeyRepoDB) RegisterKey(pixKey *model.PixKey) (*model.PixKey, error) {
	err := r.DB.Create(pixKey).Error
	if err != nil {
		return nil, err
	}
	return pixKey, nil
}

func (r PixKeyRepoDB) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	var pixKey model.PixKey
	r.DB.Preload("Account.Bank").First(&pixKey, "kind = ? and key = ?", kind, key)

	if pixKey.ID == "" {
		return nil, fmt.Errorf("No key was found.")
	}

	return &pixKey, nil
}

func (r PixKeyRepoDB) FindAccount(id string) (*model.Account, error) {
	var account model.Account
	r.DB.Preload("Bank").First(&account, "id = ?", id)

	if account.ID == "" {
		return nil, fmt.Errorf("No account with that ID was found.")
	}

	return &account, nil
}

func (r PixKeyRepoDB) FindBank(id string) (*model.Bank, error) {
	var bank model.Bank
	r.DB.First(&bank, "id = ?", id)

	if bank.ID == "" {
		return nil, fmt.Errorf("No bank with that ID was found.")
	}

	return &bank, nil
}
