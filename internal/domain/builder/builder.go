package builder

import (
	"gitlab.llo.su/fond/radara/internal/domain/agg"
	"gitlab.llo.su/fond/radara/internal/domain/dto"
)

type Fund interface {
	BuildFundsFromDtoSlice(fundDto dto.Fund) agg.Fund
}

type FundStat interface {
	BuildFundStatsFromDtoSlice(fundStatsDto []dto.FundStat, fundObjectId interface{}) []agg.FundStat
}
