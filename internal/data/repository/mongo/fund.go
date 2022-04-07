package mongorepository

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.llo.su/fond/radara/internal/domain/agg"
)

const (
	ErrNoDocumentsWasFound = "mongo: no documents in result"
)

type Fund struct {
	coll *mongo.Collection
	mu   *sync.Mutex
	buf  []interface{}
}

func NewFund(collection *mongo.Collection) Fund {
	return Fund{
		coll: collection,
		mu:   &sync.Mutex{},
		buf:  []interface{}{},
	}
}

// InsertOne - will return an objectId with errors as nil, or objectId as nil with error
func (repo Fund) InsertOne(ctx context.Context, fund agg.Fund) (interface{}, error) {
	repo.mu.Lock()

	result, err := repo.coll.InsertOne(ctx, fund, options.InsertOne())
	if err != nil {
		return nil, err
	}

	repo.mu.Unlock()

	return result.InsertedID, nil
}

func (repo Fund) FindByTitle(ctx context.Context, title string) (agg.Fund, error) {
	var resAggFund agg.Fund

	snglRes := repo.coll.FindOne(ctx, bson.M{"fund.title": title}, options.FindOne())

	if err := snglRes.Decode(&resAggFund); err != nil {
		if err.Error() == ErrNoDocumentsWasFound {
			return resAggFund, nil
		}

		return resAggFund, err
	}

	return resAggFund, nil
}
