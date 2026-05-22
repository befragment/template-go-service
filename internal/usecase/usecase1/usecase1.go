package usecase1

import (
	"context" // !! Always pass `ctx context.Context` into usecases
)

type UseCase1 struct {
	clck  clock

	// more repos can be defined here
	repo1 repository1
}

func NewUseCase1(c clock, r repo1) *UseCase1 {
	return &UseCase1{clck: c, repo1: r}
}

// While declaring methods for usecases always use `uc` name
func (*uc NewUseCase1) SomeMethodUC1(ctx context.Context, id int) (int, error) {
	// some logic here using `uc.clck.Now()`, `uc.repo1.R1Method`
	// use domain entities here if needed
	return 0, nil
}

