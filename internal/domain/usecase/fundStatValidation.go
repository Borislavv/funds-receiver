package usecase

import (
	"fmt"

	"gitlab.llo.su/fond/radara/internal/domain/dto"
)

const ErrorPostFix string = "field required but got empty"

func ValidateFundStats(fundStatsDto []dto.FundStat) error {
	for _, fundStatDto := range fundStatsDto {
		if fundStatDto.Holding == "" {
			return fmt.Errorf("holding %s", ErrorPostFix)
		}

		if fundStatDto.Symbol == "" {
			return fmt.Errorf("symbol %s", ErrorPostFix)
		}
	}

	return nil
}
