package usecase1

import (
	"context"
	"time"

	// note that domain entities can be imported and used in any layer
	"github.com/befragment/template-go/internal/domain/entity1" 
)

// Here we define interface for the object `Clock` from /pkg
// Then we use this interface in usecase1 like `uc.clock.Now()`
type clock interface {
	Now() time.Time
}

// We also use interfaces for the repository
// Call it like `uc.repo.R1method()`
type repository1 interface {
	R1method() (domain.NewDomainEntity1, error)
}

// Add logger if needed
type logger interface {
	Info(args ...interface{})
}

// UUID string generation if needed
type uuidgen interface {
	New() string
}

// use for hashing passwords
type hasher interface {
	Hash(_ context.Context, password string) (string, error)
	Compare(_ context.Context, hash, password string) error
}
