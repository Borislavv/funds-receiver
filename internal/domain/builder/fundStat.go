package builder

import (
	"time"

	"github.com/golang-module/carbon"

	"gitlab.llo.su/fond/radara/internal/domain/agg"
	"gitlab.llo.su/fond/radara/internal/domain/dto"
	"gitlab.llo.su/fond/radara/internal/domain/entity"
	"gitlab.llo.su/fond/radara/internal/domain/vo"
)

func BuildFundStatsFromDtoSlice(fundStatsDto []dto.FundStat, fundObjectId interface{}) []agg.FundStat {
	var fundStatsAgg []agg.FundStat

	for _, fundStatDto := range fundStatsDto {
		fundStatAgg := agg.FundStat{
			FundStat: entity.FundStat{
				Holding:   fundStatDto.Holding,
				Symbol:    fundStatDto.Symbol,
				Weighting: fundStatDto.Weighting,
				FileId:    fundStatDto.FileId,
				Date:      carbon.Parse(fundStatDto.Date).Time,
			},
			Fund: fundObjectId,
			Timestamp: vo.Timestamp{
				CreatedAt: time.Now(),
			},
		}

		fundStatsAgg = append(fundStatsAgg, fundStatAgg)
	}

	return fundStatsAgg
}
