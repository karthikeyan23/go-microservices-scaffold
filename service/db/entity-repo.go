package db

import (
	"context"
	"database/sql"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go_scafold/service/model"
)

type repository struct {
	db     *sql.DB
	logger log.Logger
}

// New returns a concrete repository backed by Postgres.
func New(db *sql.DB, logger log.Logger) (model.Repository, error) {
	return &repository{
		db:     db,
		logger: log.With(logger, "repository", "postgres"),
	}, nil
}

func (r repository) Create(ctx context.Context, entity *model.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Get(ctx context.Context, id string) (*model.Entity, error) {
	var aEntity model.Entity
	err := r.db.QueryRowContext(ctx, "SELECT id, name, created_at FROM entity WHERE id = $1", id).Scan(
		&aEntity.ID, &aEntity.Name, &aEntity.CreatedAt)
	if err != nil {
		_ = level.Error(r.logger).Log("err", err)
		return nil, err
	}
	return &aEntity, nil
}

func (r repository) GetAll(ctx context.Context) ([]*model.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(ctx context.Context, entity *model.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
