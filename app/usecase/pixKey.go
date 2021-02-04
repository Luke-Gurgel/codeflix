package app

import (
	"github.com/Luke-Gurgel/codeflix/domain/model"
)

type PixKeyUseCase struct {
	PixKeyRepository model.PixKeyRepositoryInterface
}

func (p *PixKeyUseCase) RegisterKey(key string, kind string, accountID string) (*model.PixKey, error) {
	account, err := p.PixKeyRepository.FindAccount(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.CreatePixKey(kind, account, key)
	if err != nil {
		return nil, err
	}

	pixKey, err = p.PixKeyRepository.RegisterKey(pixKey)
	if err != nil {
		return nil, err
	}

	return pixKey, nil
}

func (p *PixKeyUseCase) FindKey(key string, kind string) (*model.PixKey, error) {
	pixKey, err := p.PixKeyRepository.FindKeyByKind(key, kind)
	if err != nil {
		return nil, err
	}
	return pixKey, nil
}
