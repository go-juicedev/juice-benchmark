package benchmark

import (
	"context"
	"database/sql"
)

//go:generate juicecli impl -t UserRepository -o repository_impl.go
type UserRepository interface {
	Create(ctx context.Context, user *JuiceUser) (sql.Result, error)
	BatchCreate(ctx context.Context, users []*JuiceUser) (sql.Result, error)
	QueryAll(ctx context.Context) ([]*JuiceUser, error)
	QueryWithLimit(ctx context.Context, limit int) ([]*JuiceUser, error)
}
