package mongorepository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.llo.su/fond/radara/internal/domain/agg"
)

type FundStat struct {
	coll *mongo.Collection
	mu   *sync.Mutex
	buf  []interface{}
}

func NewFundStat(collection *mongo.Collection) FundStat {
	return FundStat{
		coll: collection,
		mu:   &sync.Mutex{},
		buf:  []interface{}{},
	}
}

func (repo FundStat) InsertMany(ctx context.Context, fundStats []agg.FundStat) error {
	repo.mu.Lock()

	repo.buf = repo.buf[:0]
	for _, fundStat := range fundStats {
		repo.buf = append(
			repo.buf,
			bson.D{
				{"fund", fundStat.Fund},
				{"holding", fundStat.FundStat.Holding},
				{"symbol", fundStat.FundStat.Symbol},
				{"weighting", fundStat.FundStat.Weighting},
				{"file_id", fundStat.FundStat.FileId},
				{"date", fundStat.FundStat.Date},
				{"createdAt", fundStat.Timestamp.CreatedAt},
			},
		)
	}
	result, err := repo.coll.InsertMany(ctx, repo.buf, options.InsertMany())
	if len(result.InsertedIDs) == 0 {
		return errors.New("no one `fundStat` document was created")
	}

	repo.mu.Unlock()

	return err
}

func (repo FundStat) IsExists(
	ctx context.Context,
	date time.Time,
	holding string,
	fileId int,
) (bool, error) {
	total, err := repo.coll.CountDocuments(
		ctx,
		bson.M{
			"date":    date,
			"holding": holding,
			"file_id": fileId,
		},
	)
	if err != nil {
		return false, err
	}

	if total > 1 {
		return true, errors.New(
			fmt.Sprintf(
				"%v `fundStat` duplicates found by date: %v, holding: %v, fileId: %v",
				total,
				date,
				holding,
				fileId,
			),
		)
	} else if total == 1 {
		return true, nil
	}

	return false, nil
}
