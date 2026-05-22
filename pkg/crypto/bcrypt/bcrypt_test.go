package bcrypt

import (
	"context"
	"errors"
	"strings"
	"testing"

	xbcrypt "golang.org/x/crypto/bcrypt"
)

func TestNewBcryptHasher(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tname        string
		inputCost    int
		expectedCost int
	}{
		{
			tname:        "zero -> bcrypt.DefaultCost",
			inputCost:    0,
			expectedCost: xbcrypt.DefaultCost,
		},
		{
			tname:        "non-zero -> passed through",
			inputCost:    4,
			expectedCost: 4,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.tname, func(t *testing.T) {
			t.Parallel()

			h := NewBcryptHasher(tc.inputCost)
			if h == nil {
				t.Fatalf("expected hasher, got nil")
			}
			if h.cost != tc.expectedCost {
				t.Fatalf("expected cost %d, got %d", tc.expectedCost, h.cost)
			}
		})
	}
}

func TestBcryptHasher_Hash(t *testing.T) {
	t.Parallel()

	h := NewBcryptHasher(4)

	ctx := context.Background()

	tests := []struct {
		tname         string
		password      string
		wantHashEmpty bool
		wantErr       bool
		wantErrIs     error
		wantCompareOK bool
	}{
		{
			tname:         "empty password",
			password:      "",
			wantHashEmpty: true,
			wantErr:       true,
			wantErrIs:     ErrHashPassword,
		},
		{
			tname:         "success",
			password:      "secret-pass",
			wantHashEmpty: false,
			wantErr:       false,
			wantErrIs:     nil,
			wantCompareOK: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.tname, func(t *testing.T) {
			t.Parallel()

			hash, err := h.Hash(ctx, tc.password)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tc.wantErrIs != nil && !errors.Is(err, tc.wantErrIs) {
					t.Fatalf("expected error to match %v, got: %v", tc.wantErrIs, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			}

			if tc.wantHashEmpty {
				if hash != "" {
					t.Fatalf("expected empty hash, got: %q", hash)
				}
				return
			}

			if hash == "" {
				t.Fatalf("expected non-empty hash")
			}

			if tc.wantCompareOK {
				if err := h.Compare(ctx, hash, tc.password); err != nil {
					t.Fatalf("expected Compare to succeed, got: %v", err)
				}
			}
		})
	}
}

func TestBcryptHasher_Compare(t *testing.T) {
	t.Parallel()

	h := NewBcryptHasher(4)
	ctx := context.Background()

	hashOK, err := h.Hash(ctx, "secret-pass")
	if err != nil {
		t.Fatalf("expected Hash to succeed, got: %v", err)
	}

	tests := []struct {
		tname           string
		hash            string
		password        string
		wantErrIs       error
		wantNil         bool
		wantErrContains string
		wantErrNotIs    error
	}{
		{
			tname:     "empty hash",
			hash:      "",
			password:  "secret-pass",
			wantErrIs: ErrInvalidPassword,
			wantNil:   false,
		},
		{
			tname:     "empty password",
			hash:      hashOK,
			password:  "",
			wantErrIs: ErrInvalidPassword,
			wantNil:   false,
		},
		{
			tname:     "mismatch",
			hash:      hashOK,
			password:  "wrong-pass",
			wantErrIs: ErrInvalidPassword,
			wantNil:   false,
		},
		{
			tname:    "match",
			hash:     hashOK,
			password: "secret-pass",
			wantNil:  true,
		},
		{
			tname:           "invalid hash format",
			hash:            "not-a-bcrypt-hash",
			password:        "secret-pass",
			wantNil:         false,
			wantErrNotIs:    ErrInvalidPassword,
			wantErrContains: "compare password",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.tname, func(t *testing.T) {
			t.Parallel()

			err := h.Compare(ctx, tc.hash, tc.password)

			if tc.wantNil {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if tc.wantErrIs != nil && !errors.Is(err, tc.wantErrIs) {
				t.Fatalf("expected error to match %v, got: %v", tc.wantErrIs, err)
			}
			if tc.wantErrNotIs != nil && errors.Is(err, tc.wantErrNotIs) {
				t.Fatalf("expected error to not match %v, got: %v", tc.wantErrNotIs, err)
			}
			if tc.wantErrContains != "" && !strings.Contains(err.Error(), tc.wantErrContains) {
				t.Fatalf("expected error to contain %q, got: %v", tc.wantErrContains, err)
			}
		})
	}
}
