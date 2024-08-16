package pgs

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	Name      string
	PollCount int32
}

type DbClient struct {
	Ctx  context.Context
	Pool *pgxpool.Pool
}

func (cli *DbClient) Connect(ctx context.Context, cfg DbConfig) error {
	connectionStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
	pgxConfig, err := pgxpool.ParseConfig(connectionStr)
	if err != nil {
		return err
	}
	pgxConfig.MaxConns = cfg.PollCount

	db, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return err
	}

	err = db.Ping(ctx)
	if err != nil {
		return err
	}

	cli.Pool = db
	cli.Ctx = ctx

	return nil
}
