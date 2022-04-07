package radara

import (
	"context"
	"time"

	"github.com/serge64/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"gitlab.llo.su/fond/radara/internal/data/repository"
	mongoRepository "gitlab.llo.su/fond/radara/internal/data/repository/mongo"
	"gitlab.llo.su/fond/radara/internal/domain/usecase"
)

const (
	MongoFundCollection     = "fund"
	MongoFundStatCollection = "fundStat"
)

type Config struct {
	MongoURI      string `env:"MONGO_URI,default=mongodb://localhost:27017/"`
	MongoDatabase string `env:"MONGO_DATABASE,default=radara"`

	EtfDbUrl               string `env:"ETF_DB_URL,default=https://etfdb.com/holdings-export"`
	EtfDbFilesDir          string `env:"ETF_DB_FILES_DIR,default=/tmp/etfdb/files"`
	EtfDbIdsFilesRangeFrom int    `env:"ETF_DB_IDS_FILES_RANGE_FROM,default=1"`
	EtfDbIdsFilesRangeTo   int    `env:"ETF_DB_IDS_FILES_RANGE_TO,default=4694"`

	WorkerTimeout     int `env:"WORKER_TIMEOUT,default=2"`
	WriteMongoTimeout int `env:"WRITE_MONGO_TIMEOUT,default=1"`
	ReadMongoTimeout  int `env:"READ_MONGO_TIMEOUT,default=1"`
}

func Run() error {
	var config Config
	if err := env.Unmarshal(&config); err != nil {
		return err
	}

	// create context
	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := mongo.Connect(cancelCtx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(cancelCtx)

	// checking the db is available
	if err := client.Ping(cancelCtx, readpref.Primary()); err != nil {
		return err
	}

	// configure interface for communicate with db
	store := repository.Store{
		Fund: mongoRepository.NewFund(
			client.Database(config.MongoDatabase).Collection(MongoFundCollection),
		),
		FundStat: mongoRepository.NewFundStat(
			client.Database(config.MongoDatabase).Collection(MongoFundStatCollection),
		),
	}

	// init. Fund repository
	fundUseCase := usecase.NewFund(cancelCtx, store.Fund)
	fundUseCase.WorkerTimeout = time.Duration(config.WorkerTimeout) * time.Second
	fundUseCase.WriteTimeout = time.Duration(config.WriteMongoTimeout) * time.Second
	fundUseCase.ReadTimeout = time.Duration(config.ReadMongoTimeout) * time.Second

	// init. FundStat repository
	fundStatUseCase := usecase.NewFundStat(cancelCtx, store.FundStat)
	fundStatUseCase.WorkerTimeout = time.Duration(config.WorkerTimeout) * time.Second
	fundStatUseCase.WriteTimeout = time.Duration(config.WriteMongoTimeout) * time.Second
	fundStatUseCase.ReadTimeout = time.Duration(config.ReadMongoTimeout) * time.Second

	// init. csv files fetcher
	fetcherUseCase := usecase.NewFetcher(config.EtfDbUrl, config.EtfDbFilesDir)

	// init. parser
	parserUseCase := usecase.NewParser(
		config.EtfDbFilesDir,
		config.EtfDbIdsFilesRangeFrom,
		config.EtfDbIdsFilesRangeTo,
		fundUseCase,
		fundStatUseCase,
		fetcherUseCase,
		store.Fund,
		store.FundStat,
		ctx,
	)

	if err := parserUseCase.ParseFiles(); err != nil {
		return err
	}

	return nil
}
