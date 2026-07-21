package persistence

import (
	"context"
	"database/sql"

	"shafurui/internal/model"
	"shafurui/internal/repository"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) FindEnabledByLogin(ctx context.Context, login string) (*model.User, error) {
	const query = `
SELECT id, username, nickname, password, email
FROM ` + "`user`" + `
WHERE deletedAt IS NULL
  AND isEnabled = 1
  AND (username = ? OR email = ?)
LIMIT 1`

	return r.scanUser(r.db.QueryRowContext(ctx, query, login, login))
}

func (r *UserRepositoryImpl) FindEnabledByID(ctx context.Context, userID uint64) (*model.User, error) {
	const query = `
SELECT id, username, nickname, password, email
FROM ` + "`user`" + `
WHERE deletedAt IS NULL
  AND isEnabled = 1
  AND id = ?
LIMIT 1`

	return r.scanUser(r.db.QueryRowContext(ctx, query, userID))
}

func (r *UserRepositoryImpl) scanUser(row *sql.Row) (*model.User, error) {
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Nickname, &user.Password, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}
