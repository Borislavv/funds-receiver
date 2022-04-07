package agg

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.llo.su/fond/radara/internal/domain/entity"
)

type Fund struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Fund entity.Fund
}
