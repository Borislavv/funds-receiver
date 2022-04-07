package builder

import (
	"github.com/golang-module/carbon"

	"gitlab.llo.su/fond/radara/internal/domain/agg"
	"gitlab.llo.su/fond/radara/internal/domain/dto"
	"gitlab.llo.su/fond/radara/internal/domain/entity"
)

func BuildFundsFromDtoSlice(fundDto dto.Fund) agg.Fund {
	return agg.Fund{
		Fund: entity.Fund{
			Title:               fundDto.Title,
			InceptionDate:       carbon.Parse(fundDto.InceptionDate).Time,
			TotalAssets:         fundDto.TotalAssets,
			SharesOutstanding:   fundDto.SharesOutstanding,
			ExpenseRatio:        fundDto.ExpenseRatio,
			TracksThisIndex:     fundDto.TracksThisIndex,
			ETFDatabaseCategory: fundDto.ETFDatabaseCategory,
			Issuer:              fundDto.Issuer,
			Structure:           fundDto.Structure,
			ETFHomePage:         fundDto.ETFHomePage,
		},
	}
}
