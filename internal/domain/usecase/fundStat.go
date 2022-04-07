package usecase

import (
	"context"
	"time"

	"gitlab.llo.su/fond/radara/internal/data/repository"
	"gitlab.llo.su/fond/radara/internal/domain/builder"
	"gitlab.llo.su/fond/radara/internal/domain/dto"
)

type FundStatUseCase struct {
	WorkerTimeout time.Duration
	WriteTimeout  time.Duration
	ReadTimeout   time.Duration
	context       context.Context
	fundStatRepo  repository.FundStat
}

func NewFundStat(ctx context.Context, repo repository.FundStat) FundStatUseCase {
	return FundStatUseCase{
		WorkerTimeout: Timeout,
		WriteTimeout:  Timeout,
		ReadTimeout:   Timeout,
		context:       ctx,
		fundStatRepo:  repo,
	}
}

func (uc FundStatUseCase) AddMany(fundStatsDto []dto.FundStat, fundObjectId interface{}) error {
	if err := ValidateFundStats(fundStatsDto); err != nil {
		return err
	}

	fundStatAggs := builder.BuildFundStatsFromDtoSlice(fundStatsDto, fundObjectId)

	if err := uc.fundStatRepo.InsertMany(uc.context, fundStatAggs); err != nil {
		return err
	}

	return nil
}
