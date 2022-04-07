package usecase

import (
	"fmt"

	"gitlab.llo.su/fond/radara/internal/domain/dto"
)

func ValidateFunds(fundDto dto.Fund) error {
	if fundDto.Title == "" {
		return fmt.Errorf("title %s", ErrorPostFix)
	}

	if fundDto.InceptionDate == "" {
		return fmt.Errorf("inceptionDate %s", ErrorPostFix)
	}

	if fundDto.TracksThisIndex == "" {
		return fmt.Errorf("tracksThisIndex %s", ErrorPostFix)
	}

	if fundDto.ETFDatabaseCategory == "" {
		return fmt.Errorf("etfDatabaseCategory %s", ErrorPostFix)
	}

	if fundDto.Issuer == "" {
		return fmt.Errorf("issuer %s", ErrorPostFix)
	}

	if fundDto.Structure == "" {
		return fmt.Errorf("structure %s", ErrorPostFix)
	}

	if fundDto.ETFHomePage == "" {
		return fmt.Errorf("etfHomePage %s", ErrorPostFix)
	}

	return nil
}
