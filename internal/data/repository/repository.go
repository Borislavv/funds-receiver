package repository

import (
	"context"
	"time"

	"gitlab.llo.su/fond/radara/internal/domain/agg"
)

type Store struct {
	Fund     Fund
	FundStat FundStat
}

type Fund interface {
	InsertOne(ctx context.Context, fund agg.Fund) (interface{}, error)
	FindByTitle(ctx context.Context, title string) (agg.Fund, error)
}

type FundStat interface {
	InsertMany(ctx context.Context, fundStats []agg.FundStat) error
	IsExists(ctx context.Context, date time.Time, holding string, fileId int) (bool, error)
}
