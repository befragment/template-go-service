package handler1

type usecase1 interface {
	SomeMethodUC1(ctx context.Context, id int) (int, error)
}