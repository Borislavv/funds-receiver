package usecase

import (
	"gitlab.llo.su/fond/radara/internal/domain/dto"
)

// ParserUseCase
type Parser interface {
	ParseFiles() error
	Parse(filepath string) (dto.Fund, []dto.FundStat, error)
}

// FetcherUseCase
type Fetcher interface {
	Download(fileId int) (string, error)
}

// FundUseCase
type Fund interface {
	AddOne(fundDto dto.Fund) (interface{}, error)
}

// FundStatUseCase
type FundStat interface {
	AddMany(fundStatsDto []dto.FundStat, fundObjectId interface{}) error
}
