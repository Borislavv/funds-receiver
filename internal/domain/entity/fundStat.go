package entity

import (
	"time"
)

type FundStat struct {
	// Holding title
	//
	// required: true
	// example: Apple Inc.
	Holding string `json:"holding" bson:"holding"`

	// Symbol of holding
	//
	// required: true
	// examle: AAPL
	Symbol string `json:"symbol" bson:"symbol"`

	// Weighting of holding
	//
	// required: true
	// example: 6.85
	//			(percents)
	Weighting float64 `json:"weighting" bson:"weighting"`

	// FileId
	//
	// required: true
	// example: 29
	FileId int `json:"file_id" bson:"file_id"`

	// Date of stat. from file
	//
	// pattern: \d{4}-\d{2}-\d{2}
	// required: true
	// example: 2021-11-21
	Date time.Time `json:"date" bson:"date"`
}
