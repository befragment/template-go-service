package bcrypt

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrHashPassword    = errors.New("can't create password hash")
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(_ context.Context, password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("%w: empty password", ErrHashPassword)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrHashPassword, err)
	}

	return string(hash), nil
}

func (h *BcryptHasher) Compare(_ context.Context, hash, password string) error {
	if hash == "" || password == "" {
		return ErrInvalidPassword
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}

		return fmt.Errorf("compare password: %w", err)
	}

	return nil
}
