package agg

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.llo.su/fond/radara/internal/domain/entity"
	"gitlab.llo.su/fond/radara/internal/domain/vo"
)

type FundStat struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	FundStat  entity.FundStat
	Timestamp vo.Timestamp

	// objectId of `Fund` entity
	Fund interface{}
}
