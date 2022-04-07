package usecase

import (
	"context"
	"time"

	"gitlab.llo.su/fond/radara/internal/data/repository"
	"gitlab.llo.su/fond/radara/internal/domain/builder"
	"gitlab.llo.su/fond/radara/internal/domain/dto"
)

const (
	Timeout time.Duration = time.Duration(10) * time.Second
)

type FundUseCase struct {
	WorkerTimeout time.Duration
	WriteTimeout  time.Duration
	ReadTimeout   time.Duration
	context       context.Context
	fundRepo      repository.Fund
}

func NewFund(ctx context.Context, repo repository.Fund) FundUseCase {
	return FundUseCase{
		WorkerTimeout: Timeout,
		WriteTimeout:  Timeout,
		ReadTimeout:   Timeout,
		context:       ctx,
		fundRepo:      repo,
	}
}

func (uc FundUseCase) AddOne(fundDto dto.Fund) (interface{}, error) {
	if err := ValidateFunds(fundDto); err != nil {
		return nil, err
	}

	fundsAgg := builder.BuildFundsFromDtoSlice(fundDto)

	insertedObjectId, err := uc.fundRepo.InsertOne(uc.context, fundsAgg)
	if err != nil {
		return nil, err
	}

	return insertedObjectId, nil
}
