package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-module/carbon"

	"gitlab.llo.su/fond/radara/internal/data/repository"
	"gitlab.llo.su/fond/radara/internal/domain/dto"
)

const (
	// ratio of Fund fields to row numbers
	fundTitleLineNumber               = 0
	fundInceptionDateLineNumber       = 1
	fundStatDateLineNumber            = 2
	fundTotalAssetsLineNumber         = 3
	fundSharesOutStandingLineNumber   = 4
	fundExpenseRatioLineNumber        = 5
	fundTracksThisIndexLineNumber     = 6
	fundETFDatabaseCategoryLineNumber = 7
	fundIssuerLineNumber              = 8
	fundStructureLineNumber           = 9
	fundETFHomePageLineNumber         = 10

	// fundStats header row number
	fundStatsHeaderLineNumber = 11

	// ratio of FundStat fields to cell numbers
	fundStatHoldingCellNumber   = 0
	fundStatSymbolCellNumber    = 1
	fundStatWeightingCellNumber = 2
)

var (
	// the same data for each FundStat per file
	fundStatDate       = ""
	fundStatFileId int = 0

	// regexes
	intRgx   = regexp.MustCompile(`\d+`)
	floatRgx = regexp.MustCompile(`\d+\.\d+`)
)

// CSV parser
type ParserUseCase struct {
	etfDbFilesDir string

	etfDbIdsFilesRangeFrom int
	etfDbIdsFilesRangeTo   int

	fundUseCase     FundUseCase
	fundStatUseCase FundStatUseCase

	fetcherUseCase FetcherUseCase

	fundRepository     repository.Fund
	fundStatRepository repository.FundStat

	ctx context.Context
}

func NewParser(
	etfDbFilesDir string,
	etfDbIdsFilesRangeFrom int,
	etfDbIdsFilesRangeTo int,
	fundUseCase FundUseCase,
	fundStatUseCase FundStatUseCase,
	fetcherUseCase FetcherUseCase,
	fundRepository repository.Fund,
	fundStatRepository repository.FundStat,
	ctx context.Context,
) ParserUseCase {
	return ParserUseCase{
		etfDbFilesDir:          etfDbFilesDir,
		etfDbIdsFilesRangeFrom: etfDbIdsFilesRangeFrom,
		etfDbIdsFilesRangeTo:   etfDbIdsFilesRangeTo,
		fundUseCase:            fundUseCase,
		fundStatUseCase:        fundStatUseCase,
		fetcherUseCase:         fetcherUseCase,
		fundRepository:         fundRepository,
		fundStatRepository:     fundStatRepository,
		ctx:                    ctx,
	}
}

func (parser ParserUseCase) ParseFiles() error {
	// !!! ATTENTION: UNCOMMENT this code when auth problem will be solved !!!
	//
	// downloading and parsing these files
	// for fileId := parser.etfDbIdsFilesRangeFrom; fileId <= parser.etfDbIdsFilesRangeTo; fileId++ {
	// 	filepath, err := parser.fetcherUseCase.Download(fileId)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	fund, fundStats, err := parser.Parse(filepath)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	fmt.Printf("file %v successfully parsed\n", filepath)
	//
	// 	foundFund, err := parser.fundRepository.FindByTitle(parser.ctx, fund.Title)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	var objectId interface{}
	// 	if foundFund.ID.IsZero() {
	// 		objectId, err = parser.fundUseCase.AddOne(fund)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		objectId = foundFund.ID.Hex()
	// 	}
	//
	// 	var newFoundStat []dto.FundStat
	// 	for _, fundStat := range fundStats {
	// 		isExists, err := parser.fundStatRepository.IsExists(
	// 			parser.ctx,
	// 			carbon.Parse(fundStat.Date).Time,
	// 			fundStat.Holding,
	// 			fundStat.FileId,
	// 		)
	// 		if err != nil {
	// 			return err
	// 		}
	//
	// 		if !isExists {
	// 			newFoundStat = append(newFoundStat, fundStat)
	// 		}
	// 	}
	//
	// 	if len(newFoundStat) > 0 {
	// 		if err := parser.fundStatUseCase.AddMany(newFoundStat, insertedObjectId); err != nil {
	// 			return err
	// 		}
	// 	}
	//
	// 	// reset shared vars.
	// 	fundStatDate = ""
	// }

	// !!! ATTENTION: REMOVE this code when auth problem will be solved !!!
	//
	for fileId := parser.etfDbIdsFilesRangeFrom; fileId <= parser.etfDbIdsFilesRangeTo; fileId++ {
		fundStatFileId = fileId

		filepath := parser.etfDbFilesDir + "/" + fmt.Sprint(fileId) + FileExtension

		fund, fundStats, err := parser.Parse(filepath)
		if err != nil {
			return err
		}

		fmt.Printf("file %v successfully parsed\n", filepath)

		foundFund, err := parser.fundRepository.FindByTitle(parser.ctx, fund.Title)
		if err != nil {
			return err
		}

		var insertedObjectId interface{}
		if foundFund.ID.IsZero() {
			insertedObjectId, err = parser.fundUseCase.AddOne(fund)
			if err != nil {
				return err
			}
		} else {
			insertedObjectId = foundFund.ID.Hex()
		}

		var newFoundStat []dto.FundStat
		for _, fundStat := range fundStats {
			isExists, err := parser.fundStatRepository.IsExists(
				parser.ctx,
				carbon.Parse(fundStat.Date).Time,
				fundStat.Holding,
				fundStat.FileId,
			)
			if err != nil {
				return err
			}

			if !isExists {
				newFoundStat = append(newFoundStat, fundStat)
			}
		}

		if len(newFoundStat) > 0 {
			if err := parser.fundStatUseCase.AddMany(newFoundStat, insertedObjectId); err != nil {
				return err
			}
		}

		// reset var.
		fundStatDate = ""
	}

	return nil
}

func (parser ParserUseCase) Parse(filepath string) (dto.Fund, []dto.FundStat, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return dto.Fund{}, []dto.FundStat{}, err
	}
	defer file.Close()

	return parser.parse(csv.NewReader(file))
}

func (parser ParserUseCase) parse(reader *csv.Reader) (dto.Fund, []dto.FundStat, error) {
	// response
	var fund dto.Fund
	var fundStats []dto.FundStat

	lineIterator := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		// data parsing
		if lineIterator >= fundTitleLineNumber && lineIterator <= fundETFHomePageLineNumber {
			// Fund data parsing
			if err := parser.parseLineToFund(lineIterator, parser.parseLineSliceToString(lineIterator, line), &fund); err != nil {
				return fund, fundStats, err
			}
		} else if lineIterator == fundStatsHeaderLineNumber {
			// skipping header line
			lineIterator++
			continue
		} else if lineIterator > fundStatsHeaderLineNumber {
			// FundStat data parsing
			fundStat := dto.FundStat{
				FileId: fundStatFileId,
			}

			if err := parser.parseLineSliceToFundStat(line, &fundStat); err != nil {
				return fund, fundStats, err
			}

			fundStats = append(fundStats, fundStat)
		}

		lineIterator++
	}

	return fund, fundStats, nil
}

func (parser ParserUseCase) parseLineSliceToString(lineIterator int, line []string) string {
	splitedLine := strings.Split(line[0], ":")

	if lineIterator == fundTitleLineNumber {
		return strings.Trim(line[0], " ")
	} else if lineIterator == fundETFHomePageLineNumber {
		return strings.Trim(fmt.Sprintf("%v:%v", splitedLine[1], splitedLine[2]), " ")
	}

	return strings.Trim(splitedLine[1], " ")
}

func (parser ParserUseCase) parseLineToFund(lineIterator int, line string, fund *dto.Fund) error {
	if lineIterator == fundTitleLineNumber {
		fund.Title = line
	} else if lineIterator == fundInceptionDateLineNumber {
		fund.InceptionDate = line
	} else if lineIterator == fundStatDateLineNumber {
		if fundStatDate == "" {
			fundStatDate = line
		}
	} else if lineIterator == fundTotalAssetsLineNumber {
		totalAssets, err := parser.parseInt64(line)
		if err != nil {
			return err
		}

		fund.TotalAssets = totalAssets
	} else if lineIterator == fundSharesOutStandingLineNumber {
		sharesOutstanding, err := parser.parseInt64(line)
		if err != nil {
			return err
		}

		fund.SharesOutstanding = sharesOutstanding
	} else if lineIterator == fundExpenseRatioLineNumber {
		expenseRatio, err := parser.parseFloat64(line)
		if err != nil {
			return err
		}

		fund.ExpenseRatio = expenseRatio
	} else if lineIterator == fundTracksThisIndexLineNumber {
		fund.TracksThisIndex = line
	} else if lineIterator == fundETFDatabaseCategoryLineNumber {
		fund.ETFDatabaseCategory = line
	} else if lineIterator == fundIssuerLineNumber {
		fund.Issuer = line
	} else if lineIterator == fundStructureLineNumber {
		fund.Structure = line
	} else if lineIterator == fundETFHomePageLineNumber {
		fund.ETFHomePage = line
	}

	return nil
}

func (parser ParserUseCase) parseLineSliceToFundStat(line []string, fundStat *dto.FundStat) error {
	for cellIterator, cell := range line {
		if cellIterator == fundStatHoldingCellNumber {
			fundStat.Holding = cell
		} else if cellIterator == fundStatSymbolCellNumber {
			fundStat.Symbol = cell
		} else if cellIterator == fundStatWeightingCellNumber {
			weighting, err := parser.parseFloat64(cell)
			if err != nil {
				return err
			}

			fundStat.Weighting = weighting
		}

		fundStat.Date = fundStatDate
	}

	return nil
}

func (parser ParserUseCase) parseInt64(intAsString string) (int64, error) {
	return strconv.ParseInt(
		intRgx.FindString(intAsString),
		10,
		64,
	)
}

func (parser ParserUseCase) parseFloat64(floatAsString string) (float64, error) {
	return strconv.ParseFloat(
		floatRgx.FindString(floatAsString),
		64,
	)
}
